// Code generated — DO NOT EDIT.

package ierc20

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb2 "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

var (
	_ = bytes.Equal
	_ = errors.New
	_ = fmt.Sprintf
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = emptypb.Empty{}
	_ = pb.NewBigIntFromInt
	_ = pb2.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX
	_ = bindings.FilterOptions{}
	_ = evm.FilterLogTriggerRequest{}
	_ = cre.ResponseBufferTooSmall
	_ = rpc.API{}
	_ = json.Unmarshal
	_ = reflect.Bool
)

var IERC20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Structs

// Contract Method Inputs
type AllowanceInput struct {
	Owner   common.Address
	Spender common.Address
}

type ApproveInput struct {
	Spender common.Address
	Amount  *big.Int
}

type BalanceOfInput struct {
	Account common.Address
}

type TransferInput struct {
	Recipient common.Address
	Amount    *big.Int
}

type TransferFromInput struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
}

// Contract Method Outputs

// Errors

// Events
// The <Event>Topics struct should be used as a filter (for log triggers).
// Note: It is only possible to filter on indexed fields.
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type ApprovalTopics struct {
	Owner   common.Address
	Spender common.Address
}

type ApprovalDecoded struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

type TransferTopics struct {
	From common.Address
	To   common.Address
}

type TransferDecoded struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// Main Binding Type for IERC20
type IERC20 struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   IERC20Codec
}

type IERC20Codec interface {
	EncodeAllowanceMethodCall(in AllowanceInput) ([]byte, error)
	DecodeAllowanceMethodOutput(data []byte) (*big.Int, error)
	EncodeApproveMethodCall(in ApproveInput) ([]byte, error)
	DecodeApproveMethodOutput(data []byte) (bool, error)
	EncodeBalanceOfMethodCall(in BalanceOfInput) ([]byte, error)
	DecodeBalanceOfMethodOutput(data []byte) (*big.Int, error)
	EncodeTotalSupplyMethodCall() ([]byte, error)
	DecodeTotalSupplyMethodOutput(data []byte) (*big.Int, error)
	EncodeTransferMethodCall(in TransferInput) ([]byte, error)
	DecodeTransferMethodOutput(data []byte) (bool, error)
	EncodeTransferFromMethodCall(in TransferFromInput) ([]byte, error)
	DecodeTransferFromMethodOutput(data []byte) (bool, error)
	ApprovalLogHash() []byte
	EncodeApprovalTopics(evt abi.Event, values []ApprovalTopics) ([]*evm.TopicValues, error)
	DecodeApproval(log *evm.Log) (*ApprovalDecoded, error)
	TransferLogHash() []byte
	EncodeTransferTopics(evt abi.Event, values []TransferTopics) ([]*evm.TopicValues, error)
	DecodeTransfer(log *evm.Log) (*TransferDecoded, error)
}

func NewIERC20(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*IERC20, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20MetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &IERC20{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   codec,
	}, nil
}

type Codec struct {
	abi *abi.ABI
}

func NewCodec() (IERC20Codec, error) {
	parsed, err := abi.JSON(strings.NewReader(IERC20MetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeAllowanceMethodCall(in AllowanceInput) ([]byte, error) {
	return c.abi.Pack("allowance", in.Owner, in.Spender)
}

func (c *Codec) DecodeAllowanceMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["allowance"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeApproveMethodCall(in ApproveInput) ([]byte, error) {
	return c.abi.Pack("approve", in.Spender, in.Amount)
}

func (c *Codec) DecodeApproveMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["approve"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeBalanceOfMethodCall(in BalanceOfInput) ([]byte, error) {
	return c.abi.Pack("balanceOf", in.Account)
}

func (c *Codec) DecodeBalanceOfMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["balanceOf"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTotalSupplyMethodCall() ([]byte, error) {
	return c.abi.Pack("totalSupply")
}

func (c *Codec) DecodeTotalSupplyMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["totalSupply"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTransferMethodCall(in TransferInput) ([]byte, error) {
	return c.abi.Pack("transfer", in.Recipient, in.Amount)
}

func (c *Codec) DecodeTransferMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["transfer"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTransferFromMethodCall(in TransferFromInput) ([]byte, error) {
	return c.abi.Pack("transferFrom", in.Sender, in.Recipient, in.Amount)
}

func (c *Codec) DecodeTransferFromMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["transferFrom"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) ApprovalLogHash() []byte {
	return c.abi.Events["Approval"].ID.Bytes()
}

func (c *Codec) EncodeApprovalTopics(
	evt abi.Event,
	values []ApprovalTopics,
) ([]*evm.TopicValues, error) {
	var ownerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Owner).IsZero() {
			ownerRule = append(ownerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Owner)
		if err != nil {
			return nil, err
		}
		ownerRule = append(ownerRule, fieldVal)
	}
	var spenderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Spender).IsZero() {
			spenderRule = append(spenderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Spender)
		if err != nil {
			return nil, err
		}
		spenderRule = append(spenderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		ownerRule,
		spenderRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			// don't include empty bytes if hashed value is 0x0
			if reflect.ValueOf(h).IsZero() {
				bs[j] = []byte{}
			} else {
				bs[j] = h.Bytes()
			}
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeApproval decodes a log into a Approval struct.
func (c *Codec) DecodeApproval(log *evm.Log) (*ApprovalDecoded, error) {
	event := new(ApprovalDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Approval", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Approval"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) TransferLogHash() []byte {
	return c.abi.Events["Transfer"].ID.Bytes()
}

func (c *Codec) EncodeTransferTopics(
	evt abi.Event,
	values []TransferTopics,
) ([]*evm.TopicValues, error) {
	var fromRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.From).IsZero() {
			fromRule = append(fromRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.From)
		if err != nil {
			return nil, err
		}
		fromRule = append(fromRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.To).IsZero() {
			toRule = append(toRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.To)
		if err != nil {
			return nil, err
		}
		toRule = append(toRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		fromRule,
		toRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			// don't include empty bytes if hashed value is 0x0
			if reflect.ValueOf(h).IsZero() {
				bs[j] = []byte{}
			} else {
				bs[j] = h.Bytes()
			}
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeTransfer decodes a log into a Transfer struct.
func (c *Codec) DecodeTransfer(log *evm.Log) (*TransferDecoded, error) {
	event := new(TransferDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Transfer", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Transfer"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c IERC20) Allowance(
	runtime cre.Runtime,
	args AllowanceInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeAllowanceMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeAllowanceMethodOutput(response.Data)
	})

}

func (c IERC20) BalanceOf(
	runtime cre.Runtime,
	args BalanceOfInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeBalanceOfMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeBalanceOfMethodOutput(response.Data)
	})

}

func (c IERC20) TotalSupply(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTotalSupplyMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTotalSupplyMethodOutput(response.Data)
	})

}

func (c IERC20) WriteReport(
	runtime cre.Runtime,
	report *cre.Report,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  c.Address.Bytes(),
		Report:    report,
		GasConfig: gasConfig,
	})
}

func (c *IERC20) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

// ApprovalTrigger wraps the raw log trigger and provides decoded ApprovalDecoded data
type ApprovalTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *IERC20 // Keep reference for decoding
}

// Adapt method that decodes the log into Approval data
func (t *ApprovalTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ApprovalDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeApproval(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Approval log: %w", err)
	}

	return &bindings.DecodedLog[ApprovalDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *IERC20) LogTriggerApprovalLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ApprovalTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ApprovalDecoded]], error) {
	event := c.ABI.Events["Approval"]
	topics, err := c.Codec.EncodeApprovalTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Approval: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ApprovalTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *IERC20) FilterLogsApproval(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ApprovalLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// TransferTrigger wraps the raw log trigger and provides decoded TransferDecoded data
type TransferTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]         // Embed the raw trigger
	contract                        *IERC20 // Keep reference for decoding
}

// Adapt method that decodes the log into Transfer data
func (t *TransferTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[TransferDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeTransfer(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Transfer log: %w", err)
	}

	return &bindings.DecodedLog[TransferDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *IERC20) LogTriggerTransferLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []TransferTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[TransferDecoded]], error) {
	event := c.ABI.Events["Transfer"]
	topics, err := c.Codec.EncodeTransferTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Transfer: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &TransferTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *IERC20) FilterLogsTransfer(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.TransferLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}
