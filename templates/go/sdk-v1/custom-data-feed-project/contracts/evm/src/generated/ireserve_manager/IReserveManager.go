// Code generated — DO NOT EDIT.

package ireserve_manager

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

var IReserveManagerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"RequestReserveUpdate\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"totalMinted\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalReserve\",\"type\":\"uint256\"}],\"internalType\":\"structUpdateReserves\",\"name\":\"updateReserves\",\"type\":\"tuple\"}],\"name\":\"updateReserves\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Structs
type UpdateReserves struct {
	TotalMinted  *big.Int
	TotalReserve *big.Int
}

// Contract Method Inputs
type UpdateReservesInput struct {
	UpdateReserves UpdateReserves
}

// Contract Method Outputs

// Errors

// Events
type RequestReserveUpdate struct {
	RequestId *big.Int
}

// Main Binding Type for IReserveManager
type IReserveManager struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   IReserveManagerCodec
}

type IReserveManagerCodec interface {
	EncodeUpdateReservesMethodCall(in UpdateReservesInput) ([]byte, error)
	EncodeUpdateReservesStruct(in UpdateReserves) ([]byte, error)
	RequestReserveUpdateLogHash() []byte
	EncodeRequestReserveUpdateTopics(evt abi.Event, values []RequestReserveUpdate) ([]*evm.TopicValues, error)
	DecodeRequestReserveUpdate(log *evm.Log) (*RequestReserveUpdate, error)
}

func NewIReserveManager(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*IReserveManager, error) {
	parsed, err := abi.JSON(strings.NewReader(IReserveManagerMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &IReserveManager{
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

func NewCodec() (IReserveManagerCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(IReserveManagerMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeUpdateReservesMethodCall(in UpdateReservesInput) ([]byte, error) {
	return c.abi.Pack("updateReserves", in.UpdateReserves)
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
	values []RequestReserveUpdate,
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
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeRequestReserveUpdate decodes a log into a RequestReserveUpdate struct.
func (c *Codec) DecodeRequestReserveUpdate(log *evm.Log) (*RequestReserveUpdate, error) {
	event := new(RequestReserveUpdate)
	if err := c.abi.UnpackIntoInterface(event, "RequestReserveUpdate", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["RequestReserveUpdate"].Inputs {
		if arg.Indexed {
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

func (c IReserveManager) WriteReportFromUpdateReserves(
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

func (c IReserveManager) WriteReport(
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

func (c *IReserveManager) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	default:
		return nil, errors.New("unknown error selector")
	}
}

func (c *IReserveManager) LogTriggerRequestReserveUpdateLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []RequestReserveUpdate) (cre.Trigger[*evm.Log, *evm.Log], error) {
	event := c.ABI.Events["RequestReserveUpdate"]
	topics, err := c.Codec.EncodeRequestReserveUpdateTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for RequestReserveUpdate: %w", err)
	}

	return evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	}), nil
}

func (c *IReserveManager) FilterLogsRequestReserveUpdate(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
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
