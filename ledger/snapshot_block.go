package ledger

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/vitepb"
	"time"
)

var snapshotBlockLog = log15.New("module", "ledger/snapshot_block")

type SnapshotContent map[types.Address]*SnapshotContentItem

func (sc SnapshotContent) DeProto(pb *vitepb.SnapshotContent) {
	for addrString, snapshotItem := range pb.Content {
		addr, _ := types.HexToAddress(addrString)
		accountBlockHash, _ := types.BytesToHash(snapshotItem.AccountBlockHash)

		sc[addr] = &SnapshotContentItem{
			AccountBlockHeight: snapshotItem.AccountBlockHeight,
			AccountBlockHash:   accountBlockHash,
		}
	}
}

func (sc SnapshotContent) Proto() *vitepb.SnapshotContent {
	pb := &vitepb.SnapshotContent{
		Content: make(map[string]*vitepb.SnapshotItem),
	}

	for addr, snapshotItem := range sc {
		pb.Content[addr.String()] = &vitepb.SnapshotItem{
			AccountBlockHash:   snapshotItem.AccountBlockHash.Bytes(),
			AccountBlockHeight: snapshotItem.AccountBlockHeight,
		}
	}
	return pb
}

func (sc *SnapshotContent) Serialize() ([]byte, error) {
	pb := sc.Proto()
	buf, err := proto.Marshal(pb)
	if err != nil {
		snapshotBlockLog.Error("proto.Marshal failed, error is "+err.Error(), "method", "SnapshotContent.Serialize")
	}
	return buf, nil
}
func (sc *SnapshotContent) Deserialize(buf []byte) error {
	pb := &vitepb.SnapshotContent{}
	unmarshalErr := proto.Unmarshal(buf, pb)
	if unmarshalErr != nil {
		snapshotBlockLog.Error("proto.Unmarshal failed, error is "+unmarshalErr.Error(), "method", "SnapshotContent.Deserialize")
	}

	sc.DeProto(pb)
	return nil
}

func (*SnapshotContent) Hash() types.Hash {
	return types.Hash{}
}

type SnapshotContentItem struct {
	AccountBlockHeight uint64
	AccountBlockHash   types.Hash
}

type SnapshotBlock struct {
	Hash types.Hash

	PrevHash types.Hash
	Height   uint64
	producer *types.Address

	PublicKey ed25519.PublicKey
	Signature []byte

	Timestamp *time.Time

	SnapshotHash    *types.Hash
	SnapshotContent SnapshotContent
}

func (sb *SnapshotBlock) ComputeHash() types.Hash {
	var source []byte
	// PrevHash
	source = append(source, sb.PrevHash.Bytes()...)

	// Height
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, sb.Height)
	source = append(source, heightBytes...)

	// PublicKey
	source = append(source, sb.PublicKey...)

	// Timestamp
	unixTimeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(unixTimeBytes, uint64(sb.Timestamp.Unix()))
	source = append(source, unixTimeBytes...)

	// SnapshotHash
	if sb.SnapshotHash != nil {
		source = append(source, sb.SnapshotHash.Bytes()...)
	}

	hash, _ := types.BytesToHash(crypto.Hash256(source))
	return hash
}

func (sb *SnapshotBlock) Producer() types.Address {
	if sb.producer == nil {
		producer := types.PubkeyToAddress(sb.PublicKey)
		sb.producer = &producer
	}
	return *sb.producer
}

func (sb *SnapshotBlock) VerifySignature() bool {
	isVerified, verifyErr := crypto.VerifySig(sb.PublicKey, sb.Hash.Bytes(), sb.Signature)
	if verifyErr != nil {
		accountBlockLog.Error("crypto.VerifySig failed, error is "+verifyErr.Error(), "method", "VerifySignature")
	}
	return isVerified
}

func (sb *SnapshotBlock) proto() *vitepb.SnapshotBlock {
	pb := &vitepb.SnapshotBlock{}
	pb.Hash = sb.Hash.Bytes()
	pb.PrevHash = sb.PrevHash.Bytes()
	pb.Height = sb.Height
	pb.PublicKey = sb.PublicKey
	pb.Signature = sb.Signature
	pb.Timestamp = sb.Timestamp.UnixNano()
	if sb.SnapshotHash != nil {
		pb.SnapshotHash = sb.SnapshotHash.Bytes()
	}
	return pb
}

func (sb *SnapshotBlock) Proto() *vitepb.SnapshotBlock {
	pb := sb.proto()
	if sb.SnapshotContent != nil {
		pb.SnapshotContent = sb.SnapshotContent.Proto()
	}

	return pb
}

func (sb *SnapshotBlock) DeProto(pb *vitepb.SnapshotBlock) {
	sb.Hash, _ = types.BytesToHash(pb.Hash)
	sb.PrevHash, _ = types.BytesToHash(pb.PrevHash)
	sb.Height = pb.Height
	sb.PublicKey = pb.PublicKey
	sb.Signature = pb.Signature

	timestamp := time.Unix(0, pb.Timestamp)
	sb.Timestamp = &timestamp

	if len(pb.SnapshotHash) >= 0 {
		snapshotHash, _ := types.BytesToHash(pb.SnapshotHash)
		sb.SnapshotHash = &snapshotHash
	}

	if pb.SnapshotContent != nil {
		sb.SnapshotContent = SnapshotContent{}
		sb.SnapshotContent.DeProto(pb.SnapshotContent)
	}
}

func (sb *SnapshotBlock) DbSerialize() ([]byte, error) {
	pb := sb.proto()
	buf, err := proto.Marshal(pb)
	if err != nil {
		snapshotBlockLog.Error("proto.Marshal failed, error is "+err.Error(), "method", "SnapshotBlock.DbSerialize")
	}
	return buf, nil
}

func (sb *SnapshotBlock) Serialize() ([]byte, error) {
	pb := sb.Proto()
	buf, err := proto.Marshal(pb)
	if err != nil {
		snapshotBlockLog.Error("proto.Marshal failed, error is "+err.Error(), "method", "SnapshotBlock.Serialize")
	}
	return buf, nil
}

func (sb *SnapshotBlock) Deserialize(buf []byte) error {
	pb := &vitepb.SnapshotBlock{}
	unmarshalErr := proto.Unmarshal(buf, pb)
	if unmarshalErr != nil {
		snapshotBlockLog.Error("proto.Unmarshal failed, error is "+unmarshalErr.Error(), "method", "SnapshotBlock.Deserialize")
	}

	sb.DeProto(pb)
	return nil
}

func GetGenesisSnapshotBlock() *SnapshotBlock {
	timestamp := time.Unix(1537361101, 0)
	genesisSnapshotBlock := &SnapshotBlock{
		Height:    1,
		Timestamp: &timestamp,
		PublicKey: GenesisPublicKey,
	}
	genesisSnapshotBlock.Hash = genesisSnapshotBlock.ComputeHash()
	genesisSnapshotBlock.Signature, _ = hex.DecodeString("2147fb12ea96ab8561c02c9333ad4e0afc8420f036107582c269bb7e2ebf16443536996bacebef17455703de8a9a6c95998ed3fb3a7a4f44adb0c196572fb20b")

	return genesisSnapshotBlock
}
