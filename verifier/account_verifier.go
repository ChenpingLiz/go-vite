package verifier

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/vitelabs/go-vite/chain"
	"github.com/vitelabs/go-vite/common/math"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto"
	"github.com/vitelabs/go-vite/generator"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/pow"
	"github.com/vitelabs/go-vite/vm/util"
	"github.com/vitelabs/go-vite/vm_db"
	"math/big"
)

type AccountType int

const (
	AccountTypeNotSure AccountType = iota
	AccountTypeGeneral
	AccountTypeContract
)

func isAccTypeGeneral(sureAccType AccountType) bool {
	if sureAccType == AccountTypeContract {
		return false
	}
	return true
}

type AccountVerifier struct {
	chain     chain.Chain
	consensus consensus

	log log15.Logger
}

func NewAccountVerifier(chain chain.Chain, consensus consensus) *AccountVerifier {
	return &AccountVerifier{
		chain:     chain,
		consensus: consensus,

		log: log15.New("module", "AccountVerifier"),
	}
}

func (v *AccountVerifier) verifyReferred(block *ledger.AccountBlock) (VerifyResult, *AccBlockPendingTask, error) {
	pendingTask := &AccBlockPendingTask{}

	accType, err := v.verifyAccAddress(block)
	if err != nil {
		return FAIL, pendingTask, err
	}
	if accType == AccountTypeNotSure {
		pendingTask.AccountTask = append(pendingTask.AccountTask,
			&AccountPendingTask{Addr: &block.AccountAddress, Hash: &block.Hash})
		return PENDING, pendingTask, nil
	}
	isGeneralAddr := isAccTypeGeneral(accType)

	if err := v.verifySelf(block, isGeneralAddr); err != nil {
		return FAIL, pendingTask, err
	}

	verifyDependencyResult, err := v.verifyDependency(pendingTask, block, isGeneralAddr)
	if verifyDependencyResult != SUCCESS {
		return verifyDependencyResult, pendingTask, err
	}
	return SUCCESS, nil, nil
}

// check address's existence and validity（Height can't be lower than 1 and sendBlock can't stand at 1)
func (v *AccountVerifier) verifyAccAddress(block *ledger.AccountBlock) (AccountType, error) {
	mLog := v.log.New("method", "verifyAccAddress")
	if block.Height < 1 {
		return AccountTypeNotSure, errors.New("block.Height mustn't be lower than 1")
	}
	isContract1stCheck, err := v.chain.IsContractAccount(&block.AccountAddress)
	if err != nil {
		mLog.Error("1st check IsContractAccount, err" + err.Error())
		return AccountTypeNotSure, err
	}
	if isContract1stCheck {
		return AccountTypeContract, nil
	}
	if block.Height > 1 {
		if firstBlock, _ := v.chain.GetAccountBlockByHeight(&block.AccountAddress, 1); firstBlock != nil {
			isContract, err := v.chain.IsContractAccount(&block.AccountAddress)
			if err != nil {
				mLog.Error("check firstBlock IsContractAccount, err" + err.Error())
				return AccountTypeNotSure, err
			}
			if isContract {
				return AccountTypeContract, nil
			}
			return AccountTypeGeneral, nil
		}
	} else {
		if block.IsReceiveBlock() {
			sendBlock, sErr := v.chain.GetAccountBlockByHash(&block.Hash)
			if sErr != nil {
				mLog.Error(sErr.Error())
			}
			if sendBlock != nil {
				isContract2ndCheck, err := v.chain.IsContractAccount(&block.AccountAddress)
				if err != nil {
					mLog.Error("2nd check IsContractAccount, err" + err.Error())
					return AccountTypeNotSure, err
				}
				if isContract2ndCheck {
					return AccountTypeContract, nil
				}
				return AccountTypeGeneral, nil
			}
		} else {
			return AccountTypeNotSure, errors.New("fatal: sendblock.height can't be 1")
		}
	}
	return AccountTypeNotSure, nil
}

// Block itself coming into Verifier indicate that it referr to a snapshot, which maybe confimed >=0;
// Contarct.recv's must check its send whether is satisfied confirmed over custom times at least;
func (v *AccountVerifier) verifyComfirmedTimes(block *ledger.AccountBlock, isGeneralAddr bool) error {
	mLog := v.log.New("method", "verfiySnapshotLimit")
	var timesLimit uint64 = 0
	if block.IsReceiveBlock() {
		if !isGeneralAddr {
			meta, err := v.chain.GetContractMeta(&block.AccountAddress)
			if err != nil {
				mLog.Error(err.Error(), "call", "GetContractMeta")
				return err
			}
			timesLimit = uint64(meta.SendConfirmedTimes)
			if timesLimit <= 0 {
				return nil
			}
			sendConfirmedTimes, err := v.chain.GetConfirmedTimes(&block.FromBlockHash)
			if err != nil {
				mLog.Error(err.Error(), "call", "GetConfirmedTimes")
				return err
			}
			if sendConfirmedTimes < timesLimit {
				mLog.Error(fmt.Sprintf("err:%v, contract(addr:%v,ct:%v), from(hash:%v,ct:%v),",
					ErrVerifyConfirmedTimesNotEnough.Error(), block.AccountAddress, timesLimit, block.FromBlockHash, sendConfirmedTimes))
				return ErrVerifyConfirmedTimesNotEnough
			}
		}
	}
	return nil
}

func (v *AccountVerifier) verifySelf(block *ledger.AccountBlock, isGeneralAddr bool) error {
	if block.IsSendBlock() {
		if err := v.verifySendBlockIntergrity(block, isGeneralAddr); err != nil {
			return err
		}
	} else {
		if err := v.verifyRecvBlockIntergrity(block, isGeneralAddr); err != nil {
			return err
		}
	}
	if isGeneralAddr || block.IsReceiveBlock() {
		if err := v.verifySignature(block); err != nil {
			return err
		}
	}
	if err := v.verifyProducerLegality(block, isGeneralAddr); err != nil {
		return err
	}
	if err := v.verifyNonce(block, isGeneralAddr); err != nil {
		return err
	}
	if err := v.verifyHash(block); err != nil {
		return err
	}
	return nil
}

func (v *AccountVerifier) verifyDependency(pendingTask *AccBlockPendingTask, block *ledger.AccountBlock, isGeneralAddr bool) (VerifyResult, error) {
	if block.IsReceiveBlock() {
		// check the existence of recvBlock's send
		sendBlock, err := v.chain.GetAccountBlockByHash(&block.FromBlockHash)
		if err != nil {
			return FAIL, err
		}
		if sendBlock == nil {
			pendingTask.AccountTask = append(pendingTask.AccountTask,
				&AccountPendingTask{Addr: nil, Hash: &block.FromBlockHash},
				&AccountPendingTask{Addr: &block.AccountAddress, Hash: &block.Hash},
			)
			return PENDING, nil
		}
		// check whether is already received
		isReceived, err := v.chain.IsReceived(&block.FromBlockHash)
		if err != nil {
			return FAIL, err
		}
		if isReceived {
			return FAIL, errors.New("block is already received successfully")
		}
	}
	// check whether it's the available time the block can be snapshoted,
	// now specially for the block's confirmedTimes which it relay on.
	if err := v.verifyComfirmedTimes(block, isGeneralAddr); err != nil {
		return FAIL, err
	}
	// check whether the prev is snapshoted
	prevBlock, err := v.chain.GetAccountBlockByHash(&block.PrevHash)
	if err != nil {
		return FAIL, err
	}
	if prevBlock == nil {
		pendingTask.AccountTask = append(pendingTask.AccountTask,
			&AccountPendingTask{Addr: &block.AccountAddress, Hash: &block.PrevHash},
			&AccountPendingTask{Addr: &block.AccountAddress, Hash: &block.Hash},
		)
		return PENDING, nil
	}
	return SUCCESS, nil
}

func (v *AccountVerifier) verifySendBlockIntergrity(block *ledger.AccountBlock, isGeneralAddr bool) error {
	if block.Amount == nil {
		block.Amount = big.NewInt(0)
	} else {
		if block.Amount.Sign() < 0 || block.Amount.BitLen() > math.MaxBigIntLen {
			return errors.New("sendBlock.Amount out of bounds")
		}
		if block.TokenId == types.ZERO_TOKENID {
			return errors.New("sendBlock.TokenId can't be ZERO_TOKENID when amount has value")
		}
	}
	if block.Fee == nil {
		block.Fee = big.NewInt(0)
	} else {
		if block.Fee.Sign() < 0 || block.Fee.BitLen() > math.MaxBigIntLen {
			return errors.New("sendBlock.Fee out of bounds")
		}
	}
	if block.FromBlockHash != types.ZERO_HASH {
		return errors.New("sendBlock.FromBlockHash must be ZERO_HASH")
	}
	if isGeneralAddr {
		if block.Height <= 1 {
			return errors.New("general's sendBlock.Height must be larger than 1")
		}
	} else {
		if block.Height != 0 {
			return errors.New("contract's sendBlock.Height must be 0")
		}
		if len(block.Signature) != 0 || len(block.PublicKey) != 0 {
			return errors.New("signature and publicKey of the contract's send must be nil")
		}
	}
	return nil
}

func (v *AccountVerifier) verifyRecvBlockIntergrity(block *ledger.AccountBlock, isGeneralAddr bool) error {
	mLog := v.log.New("method", "verifyRecvBlockIntergrity")
	if block.TokenId != types.ZERO_TOKENID {
		return errors.New("recvBlock.TokenId must be ZERO_TOKENID")
	}
	if block.Amount != nil && block.Amount.Cmp(big.NewInt(0)) != 0 {
		return errors.New("recvBlock.Amount can't be anything other than nil or 0 ")
	}
	if block.Fee != nil && block.Fee.Cmp(big.NewInt(0)) != 0 {
		return errors.New("recvBlock.Fee can't be anything other than nil or 0")
	}
	if block.Height <= 0 {
		return errors.New("recvBlock.Height must be larger than 0")
	}
	if !isGeneralAddr && block.SendBlockList != nil {
		for _, sendBlock := range block.SendBlockList {
			if err := v.verifySendBlockIntergrity(sendBlock, true); err != nil {
				mLog.Error(fmt.Sprintf("err:%v, contractAddr:%v, recv-subSends(%v, %v)",
					err.Error(), block.AccountAddress, block.Hash, sendBlock.Hash))
				return err
			}
		}
	}
	return nil
}

func (v *AccountVerifier) verifySignature(block *ledger.AccountBlock) error {
	if len(block.Signature) <= 0 || len(block.PublicKey) <= 0 {
		return errors.New("signature and publicKey all must have value")
	}
	isVerified, verifyErr := crypto.VerifySig(block.PublicKey, block.Hash.Bytes(), block.Signature)
	if verifyErr != nil {
		v.log.Error(verifyErr.Error(), "call", "VerifySig")
		return ErrVerifySignatureFailed
	}
	if !isVerified {
		return ErrVerifySignatureFailed
	}
	return nil
}

func (v *AccountVerifier) verifyHash(block *ledger.AccountBlock) error {
	computedHash := block.ComputeHash()
	if block.Hash.IsZero() {
		return errors.New("hash can't be allzero")
	}
	if computedHash != block.Hash {
		//verifier.log.Error("checkHash failed", "originHash", block.Hash, "computedHash", computedHash)
		return ErrVerifyHashFailed
	}
	return nil
}

func (v *AccountVerifier) verifyNonce(block *ledger.AccountBlock, isGeneralAddr bool) error {
	if len(block.Nonce) != 0 {
		if !isGeneralAddr {
			return errors.New("nonce of contractAddr's block must be nil")
		}
		if block.Difficulty == nil {
			return errors.New("difficulty can't be nil")
		}
		if len(block.Nonce) != 8 {
			return errors.New("nonce length doesn't satisfy with 8")
		}
		hash256Data := crypto.Hash256(block.AccountAddress.Bytes(), block.PrevHash.Bytes())
		if !pow.CheckPowNonce(block.Difficulty, block.Nonce, hash256Data) {
			return ErrVerifyNonceFailed
		}
	} else {
		if block.Difficulty != nil {
			return errors.New("difficulty must be nil when nonce is nil")
		}
	}
	return nil
}

func (v *AccountVerifier) verifyProducerLegality(block *ledger.AccountBlock, isGeneralAddr bool) error {
	if isGeneralAddr {
		if types.PubkeyToAddress(block.PublicKey) != block.AccountAddress {
			return errors.New("general-account's publicKey doesn't match with the address")
		}
	} else {
		if block.IsReceiveBlock() {
			if result, err := v.consensus.VerifyAccountProducer(block); !result {
				if err != nil {
					v.log.Error(err.Error())
				}
				return errors.New("contract-block's producer is illegal")
			}
		}
	}
	return nil
}

func (v *AccountVerifier) vmVerify(block *ledger.AccountBlock, snapshotHash *types.Hash) (vmBlock *vm_db.VmAccountBlock, err error) {
	vLog := v.log.New("method", "VerifyforVM")

	var fromBlock *ledger.AccountBlock
	var states *util.GlobalStatus
	var recvErr error
	if block.IsReceiveBlock() {
		fromBlock, recvErr = v.chain.GetAccountBlockByHash(&block.FromBlockHash)
		if recvErr != recvErr {
			return nil, recvErr
		}
		if fromBlock == nil {
			return nil, errors.New("failed to find the recvBlock's fromBlock")
		}

		states, recvErr = v.chain.GetContractRandomGlobalStatus(&block.AccountAddress, &block.FromBlockHash)
		if recvErr != nil {
			return nil, recvErr
		}
	}
	gen, err := generator.NewGenerator2(v.chain, block.AccountAddress, snapshotHash, &block.PrevHash, states)
	if err != nil {
		vLog.Error("new generator error," + err.Error())
		return nil, ErrVerifyForVmGeneratorFailed
	}

	genResult, err := gen.GenerateWithBlock(*block, fromBlock)
	if err != nil {
		vLog.Error("generator block error," + err.Error())
		return nil, ErrVerifyForVmGeneratorFailed
	}
	if genResult.VmBlock == nil {
		if genResult.Err != nil {
			errInf := fmt.Sprintf("sbHash %v, addr %v,", states.SnapshotBlock.Hash, block.AccountAddress)
			if block.IsReceiveBlock() {
				errInf += fmt.Sprintf("fromHash %v", block.FromBlockHash)
			}
			vLog.Error(genResult.Err.Error(), "detail", errInf)
			return nil, genResult.Err
		}
		return nil, errors.New("vm failed, blockList is empty")
	}

	// verify vm result block's hash
	if block.Hash != vmBlock.AccountBlock.Hash {
		return nil, errors.New("Inconsistent execution results in vm.")
	}
	return vmBlock, nil
}

func (v *AccountVerifier) verifyIsReceivedSucceed(block *ledger.AccountBlock) (bool, error) {
	return v.chain.IsReceived(&block.FromBlockHash)
}

type AccBlockPendingTask struct {
	AccountTask []*AccountPendingTask
}

func (task *AccBlockPendingTask) pendingHashListToStr() string {
	var pendHashListStr string
	for k, v := range task.AccountTask {
		pendHashListStr += v.Hash.String()
		if k < len(task.AccountTask)-1 {
			pendHashListStr += "|"
		}
	}
	return pendHashListStr
}
