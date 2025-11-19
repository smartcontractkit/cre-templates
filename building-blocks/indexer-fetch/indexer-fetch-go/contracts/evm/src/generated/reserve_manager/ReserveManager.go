// Code generated — DO NOT EDIT.

package reserve_manager

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

var ReserveManagerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"lastTotalMinted\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastTotalReserve\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"event\",\"name\":\"RequestReserveUpdate\",\"inputs\":[{\"name\":\"u\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structReserveManager.UpdateReserves\",\"components\":[{\"name\":\"totalMinted\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalReserve\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false}]",
}

// Structs
type UpdateReserves struct {
	TotalMinted  *big.Int
	TotalReserve *big.Int
}

// Contract Method Inputs
type OnReportInput struct {
	Arg0   []byte
	Report []byte
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
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

type RequestReserveUpdateTopics struct {
}

type RequestReserveUpdateDecoded struct {
	U UpdateReserves
}

// Main Binding Type for ReserveManager
type ReserveManager struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   ReserveManagerCodec
}

type ReserveManagerCodec interface {
	EncodeLastTotalMintedMethodCall() ([]byte, error)
	DecodeLastTotalMintedMethodOutput(data []byte) (*big.Int, error)
	EncodeLastTotalReserveMethodCall() ([]byte, error)
	DecodeLastTotalReserveMethodOutput(data []byte) (*big.Int, error)
	EncodeOnReportMethodCall(in OnReportInput) ([]byte, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
	EncodeUpdateReservesStruct(in UpdateReserves) ([]byte, error)
	RequestReserveUpdateLogHash() []byte
	EncodeRequestReserveUpdateTopics(evt abi.Event, values []RequestReserveUpdateTopics) ([]*evm.TopicValues, error)
	DecodeRequestReserveUpdate(log *evm.Log) (*RequestReserveUpdateDecoded, error)
}

func NewReserveManager(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*ReserveManager, error) {
	parsed, err := abi.JSON(strings.NewReader(ReserveManagerMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &ReserveManager{
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

func NewCodec() (ReserveManagerCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(ReserveManagerMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeLastTotalMintedMethodCall() ([]byte, error) {
	return c.abi.Pack("lastTotalMinted")
}

func (c *Codec) DecodeLastTotalMintedMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["lastTotalMinted"].Outputs.Unpack(data)
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

func (c *Codec) EncodeLastTotalReserveMethodCall() ([]byte, error) {
	return c.abi.Pack("lastTotalReserve")
}

func (c *Codec) DecodeLastTotalReserveMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["lastTotalReserve"].Outputs.Unpack(data)
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

func (c *Codec) EncodeOnReportMethodCall(in OnReportInput) ([]byte, error) {
	return c.abi.Pack("onReport", in.Arg0, in.Report)
}

func (c *Codec) EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error) {
	return c.abi.Pack("supportsInterface", in.InterfaceId)
}

func (c *Codec) DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["supportsInterface"].Outputs.Unpack(data)
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

func (c *Codec) EncodeUpdateReservesStruct(in UpdateReserves) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "totalMinted", Type: "uint256"},
			{Name: "totalReserve", Type: "uint256"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for UpdateReserves: %w", err)
	}
	args := abi.Arguments{
		{Name: "updateReserves", Type: tupleType},
	}

	return args.Pack(in)
}

func (c *Codec) RequestReserveUpdateLogHash() []byte {
	return c.abi.Events["RequestReserveUpdate"].ID.Bytes()
}

func (c *Codec) EncodeRequestReserveUpdateTopics(
	evt abi.Event,
	values []RequestReserveUpdateTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
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

// DecodeRequestReserveUpdate decodes a log into a RequestReserveUpdate struct.
func (c *Codec) DecodeRequestReserveUpdate(log *evm.Log) (*RequestReserveUpdateDecoded, error) {
	event := new(RequestReserveUpdateDecoded)
	if err := c.abi.UnpackIntoInterface(event, "RequestReserveUpdate", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["RequestReserveUpdate"].Inputs {
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

func (c ReserveManager) LastTotalMinted(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLastTotalMintedMethodCall()
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
		return c.Codec.DecodeLastTotalMintedMethodOutput(response.Data)
	})

}

func (c ReserveManager) LastTotalReserve(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLastTotalReserveMethodCall()
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
		return c.Codec.DecodeLastTotalReserveMethodOutput(response.Data)
	})

}

func (c ReserveManager) WriteReportFromUpdateReserves(
	runtime cre.Runtime,
	input UpdateReserves,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeUpdateReservesStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c ReserveManager) WriteReport(
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

func (c *ReserveManager) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

// RequestReserveUpdateTrigger wraps the raw log trigger and provides decoded RequestReserveUpdateDecoded data
type RequestReserveUpdateTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *ReserveManager // Keep reference for decoding
}

// Adapt method that decodes the log into RequestReserveUpdate data
func (t *RequestReserveUpdateTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[RequestReserveUpdateDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeRequestReserveUpdate(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RequestReserveUpdate log: %w", err)
	}

	return &bindings.DecodedLog[RequestReserveUpdateDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ReserveManager) LogTriggerRequestReserveUpdateLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []RequestReserveUpdateTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[RequestReserveUpdateDecoded]], error) {
	event := c.ABI.Events["RequestReserveUpdate"]
	topics, err := c.Codec.EncodeRequestReserveUpdateTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for RequestReserveUpdate: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &RequestReserveUpdateTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ReserveManager) FilterLogsRequestReserveUpdate(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.RequestReserveUpdateLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}
