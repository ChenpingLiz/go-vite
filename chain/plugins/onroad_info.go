package chain_plugins

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/vitelabs/go-vite/chain/db"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/vitepb"
	"math/big"
	"sync"
)

type OnRoadInfo struct {
	chain Chain

	store *chain_db.Store
	mu    sync.RWMutex
}

func newOnRoadInfo(store *chain_db.Store, chain Chain) Plugin {
	return &OnRoadInfo{
		store: store,
		chain: chain,
	}
}

func (or *OnRoadInfo) InsertAccountBlocks(blocks []*ledger.AccountBlock) error {
	or.mu.Lock()
	defer or.mu.Unlock()

	batch := or.store.NewBatch()
	for _, v := range blocks {
		if err := or.writeOnRoadInfo(batch, v); err != nil {
			return err
		}
	}
	or.store.Write(batch)
	return nil
}

func (or *OnRoadInfo) InsertSnapshotBlocks(blocks []*ledger.SnapshotBlock) error {
	return nil
}

func (or *OnRoadInfo) DeleteChunks(chunks []*ledger.SnapshotChunk) error {
	or.mu.Lock()
	defer or.mu.Unlock()

	batch := or.store.NewBatch()
	for _, chunk := range chunks {
		for _, block := range chunk.AccountBlocks {
			if err := or.deleteOnRoadInfo(batch, block); err != nil {
				return err
			}
		}
	}
	or.store.Write(batch)
	return nil
}

func (or *OnRoadInfo) GetOnRoadAccountInfo(addr *types.Address) (*ledger.OnRoadAccountInfo, error) {
	or.mu.RLock()
	defer or.mu.RUnlock()

	omMap, err := or.readOnRoadInfo(addr)
	if err != nil {
		return nil, err
	}
	onroadInfo := &ledger.OnRoadAccountInfo{
		AccountAddress:      *addr,
		TotalNumber:         0,
		TokenBalanceInfoMap: make(map[types.TokenTypeId]*ledger.TokenBalanceInfo),
	}
	balanceMap := onroadInfo.TokenBalanceInfoMap
	for k, v := range omMap {
		balanceMap[k] = &ledger.TokenBalanceInfo{
			TotalAmount: v.TotalAmount,
			Number:      v.Number,
		}
		onroadInfo.TotalNumber += v.Number
	}
	return onroadInfo, nil
}

func (or *OnRoadInfo) writeOnRoadInfo(batch *leveldb.Batch, block *ledger.AccountBlock) error {
	if block.IsSendBlock() {
		key := CreateOnRoadInfoKey(&block.ToAddress, &block.TokenId)
		om, err := or.getMeta(key)
		if err != nil {
			return err
		}
		if om != nil {
			om.TotalAmount.Add(&om.TotalAmount, block.Amount)
		} else {
			totalAmount := big.NewInt(0)
			if block.Amount != nil {
				totalAmount.Add(totalAmount, block.Amount)
			}
			om.TotalAmount = *totalAmount
		}
		om.Number++
		return or.writeMeta(batch, key, om)
	} else {
		fromBlock, err := or.chain.GetAccountBlockByHash(&block.FromBlockHash)
		if err != nil {
			return err
		}
		if fromBlock == nil {
			return errors.New("failed to find onroad by recv")
		}
		key := CreateOnRoadInfoKey(&fromBlock.ToAddress, &fromBlock.TokenId)
		om, err := or.getMeta(key)
		if err != nil {
			return err
		}
		if om == nil {
			return errors.New("conflict, failed to remove onroad cause info meta is nil")
		}
		if om.TotalAmount.Cmp(fromBlock.Amount) == -1 {
			return errors.New("conflict with amount of onroad info")
		} else if om.TotalAmount.Cmp(fromBlock.Amount) == 0 {
			or.deleteMeta(batch, key)
			return nil
		} else {
			om.TotalAmount.Sub(&om.TotalAmount, fromBlock.Amount)
			om.Number--
			return or.writeMeta(batch, key, om)
		}
	}
}

func (or *OnRoadInfo) deleteOnRoadInfo(batch *leveldb.Batch, block *ledger.AccountBlock) error {
	if block.IsSendBlock() {
		key := CreateOnRoadInfoKey(&block.ToAddress, &block.TokenId)
		om, err := or.getMeta(key)
		if err != nil {
			return err
		}
		if om == nil {
			return errors.New("conflict, failed to remove onroad cause info meta is nil")
		}
		if om.TotalAmount.Cmp(block.Amount) == -1 {
			return errors.New("conflict with amount of onroad info")
		} else if om.TotalAmount.Cmp(block.Amount) == 0 {
			or.deleteMeta(batch, key)
			return nil
		} else {
			om.TotalAmount.Sub(&om.TotalAmount, block.Amount)
			om.Number--
			return or.writeMeta(batch, key, om)
		}
	} else {
		fromBlock, err := or.chain.GetAccountBlockByHash(&block.FromBlockHash)
		if err != nil {
			return err
		}
		if fromBlock == nil {
			return errors.New("failed to find onroad by recv")
		}
		key := CreateOnRoadInfoKey(&fromBlock.ToAddress, &fromBlock.TokenId)
		om, err := or.getMeta(key)
		if err != nil {
			return err
		}
		if om != nil {
			om.TotalAmount.Add(&om.TotalAmount, fromBlock.Amount)
		} else {
			totalAmount := big.NewInt(0)
			if fromBlock.Amount != nil {
				totalAmount.Add(totalAmount, fromBlock.Amount)
			}
			om.TotalAmount = *totalAmount
		}
		om.Number++
		return or.writeMeta(batch, key, om)
	}
}

func (or *OnRoadInfo) readOnRoadInfo(addr *types.Address) (map[types.TokenTypeId]*onroadMeta, error) {
	omMap := make(map[types.TokenTypeId]*onroadMeta)
	iter := or.store.NewIterator(util.BytesPrefix(CreateOnRoadInfoPrefixKey(addr)))
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		tokenTypeIdBytes := key[1+types.AddressSize : 1+types.AddressSize+types.TokenTypeIdSize]
		tokenTypeId, err := types.BytesToTokenTypeId(tokenTypeIdBytes)
		if err != nil {
			return nil, err
		}
		om := &onroadMeta{}
		if err := om.deserialize(iter.Value()); err != nil {
			return nil, err
		}
		omMap[tokenTypeId] = om
	}
	if err := iter.Error(); err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	return omMap, nil
}

func (or *OnRoadInfo) getMeta(key []byte) (*onroadMeta, error) {
	value, err := or.store.Get(key)
	if err != nil {
		return nil, err
	}
	if len(value) <= 0 {
		return nil, nil
	}
	om := &onroadMeta{}
	if err := om.deserialize(value); err != nil {
		return nil, err
	}
	return om, nil
}

func (or *OnRoadInfo) writeMeta(batch *leveldb.Batch, key []byte, meta *onroadMeta) error {
	dataSlice, sErr := meta.serialize()
	if sErr != nil {
		return sErr
	}
	batch.Put(key, dataSlice)
	return nil
}

func (or *OnRoadInfo) deleteMeta(batch *leveldb.Batch, key []byte) {
	batch.Delete(key)
}

type onroadMeta struct {
	TotalAmount big.Int
	Number      uint64
}

func (om *onroadMeta) proto() *vitepb.OnroadMeta {
	pb := &vitepb.OnroadMeta{}
	pb.Num = om.Number
	pb.Amount = om.TotalAmount.Bytes()
	return pb
}

func (om *onroadMeta) deProto(pb *vitepb.OnroadMeta) {
	om.Number = pb.Num
	totalAmount := big.NewInt(0)
	if len(pb.Amount) > 0 {
		totalAmount.SetBytes(pb.Amount)
	}
	om.TotalAmount = *totalAmount
}

func (om *onroadMeta) serialize() ([]byte, error) {
	return proto.Marshal(om.proto())
}

func (om *onroadMeta) deserialize(buf []byte) error {
	pb := &vitepb.OnroadMeta{}
	if err := proto.Unmarshal(buf, pb); err != nil {
		return err
	}
	om.deProto(pb)
	return nil
}

func CreateOnRoadInfoKey(addr *types.Address, tId *types.TokenTypeId) []byte {
	key := make([]byte, 0, 1+types.AddressSize+types.TokenTypeIdSize)
	key = append(key, OnRoadInfoKeyPrefix)
	key = append(key, addr.Bytes()...)
	key = append(key, tId.Bytes()...)
	return key
}

func CreateOnRoadInfoPrefixKey(addr *types.Address) []byte {
	key := make([]byte, 0, 1+types.AddressSize)
	key = append(key, OnRoadInfoKeyPrefix)
	key = append(key, addr.Bytes()...)
	return key
}