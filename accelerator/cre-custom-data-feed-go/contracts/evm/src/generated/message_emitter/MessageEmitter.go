// Code generated — DO NOT EDIT.

package message_emitter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
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
)

var MessageEmitterMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"MessageEmitted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"emitMessage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"}],\"name\":\"getLastMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"emitter\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Structs

// Contract Method Inputs
type EmitMessageInput struct {
	Message string
}

type GetLastMessageInput struct {
	Emitter common.Address
}

type GetMessageInput struct {
	Emitter   common.Address
	Timestamp *big.Int
}

// Contract Method Outputs

// Errors

// Events
// The <Event> struct should be used as a filter (for log triggers).
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type MessageEmitted struct {
	Emitter   common.Address
	Timestamp *big.Int
	Message   string
}

type MessageEmittedDecoded struct {
	Emitter   common.Address
	Timestamp *big.Int
	Message   string
}

// Main Binding Type for MessageEmitter
type MessageEmitter struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   MessageEmitterCodec
}

type MessageEmitterCodec interface {
	EncodeEmitMessageMethodCall(in EmitMessageInput) ([]byte, error)
	EncodeGetLastMessageMethodCall(in GetLastMessageInput) ([]byte, error)
	DecodeGetLastMessageMethodOutput(data []byte) (string, error)
	EncodeGetMessageMethodCall(in GetMessageInput) ([]byte, error)
	DecodeGetMessageMethodOutput(data []byte) (string, error)
	EncodeTypeAndVersionMethodCall() ([]byte, error)
	DecodeTypeAndVersionMethodOutput(data []byte) (string, error)
	MessageEmittedLogHash() []byte
	EncodeMessageEmittedTopics(evt abi.Event, values []MessageEmitted) ([]*evm.TopicValues, error)
	DecodeMessageEmitted(log *evm.Log) (*MessageEmittedDecoded, error)
}

func NewMessageEmitter(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*MessageEmitter, error) {
	parsed, err := abi.JSON(strings.NewReader(MessageEmitterMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &MessageEmitter{
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

func NewCodec() (MessageEmitterCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(MessageEmitterMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeEmitMessageMethodCall(in EmitMessageInput) ([]byte, error) {
	return c.abi.Pack("emitMessage", in.Message)
}

func (c *Codec) EncodeGetLastMessageMethodCall(in GetLastMessageInput) ([]byte, error) {
	return c.abi.Pack("getLastMessage", in.Emitter)
}

func (c *Codec) DecodeGetLastMessageMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["getLastMessage"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetMessageMethodCall(in GetMessageInput) ([]byte, error) {
	return c.abi.Pack("getMessage", in.Emitter, in.Timestamp)
}

func (c *Codec) DecodeGetMessageMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["getMessage"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTypeAndVersionMethodCall() ([]byte, error) {
	return c.abi.Pack("typeAndVersion")
}

func (c *Codec) DecodeTypeAndVersionMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["typeAndVersion"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) MessageEmittedLogHash() []byte {
	return c.abi.Events["MessageEmitted"].ID.Bytes()
}

func (c *Codec) EncodeMessageEmittedTopics(
	evt abi.Event,
	values []MessageEmitted,
) ([]*evm.TopicValues, error) {
	var emitterRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Emitter)
		if err != nil {
			return nil, err
		}
		emitterRule = append(emitterRule, fieldVal)
	}
	var timestampRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Timestamp)
		if err != nil {
			return nil, err
		}
		timestampRule = append(timestampRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		emitterRule,
		timestampRule,
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
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeMessageEmitted decodes a log into a MessageEmitted struct.
func (c *Codec) DecodeMessageEmitted(log *evm.Log) (*MessageEmittedDecoded, error) {
	event := new(MessageEmittedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "MessageEmitted", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["MessageEmitted"].Inputs {
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

func (c MessageEmitter) GetLastMessage(
	runtime cre.Runtime,
	args GetLastMessageInput,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeGetLastMessageMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeGetLastMessageMethodOutput(response.Data)
	})

}

func (c MessageEmitter) GetMessage(
	runtime cre.Runtime,
	args GetMessageInput,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeGetMessageMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeGetMessageMethodOutput(response.Data)
	})

}

func (c MessageEmitter) TypeAndVersion(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeTypeAndVersionMethodCall()
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeTypeAndVersionMethodOutput(response.Data)
	})

}

func (c MessageEmitter) WriteReport(
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

func (c *MessageEmitter) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

// MessageEmittedTrigger wraps the raw log trigger and provides decoded MessageEmittedDecoded data
type MessageEmittedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *MessageEmitter // Keep reference for decoding
}

// Adapt method that decodes the log into MessageEmitted data
func (t *MessageEmittedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[MessageEmittedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeMessageEmitted(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MessageEmitted log: %w", err)
	}

	return &bindings.DecodedLog[MessageEmittedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *MessageEmitter) LogTriggerMessageEmittedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []MessageEmitted) (cre.Trigger[*evm.Log, *bindings.DecodedLog[MessageEmittedDecoded]], error) {
	event := c.ABI.Events["MessageEmitted"]
	topics, err := c.Codec.EncodeMessageEmittedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for MessageEmitted: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &MessageEmittedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *MessageEmitter) FilterLogsMessageEmitted(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.MessageEmittedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}
