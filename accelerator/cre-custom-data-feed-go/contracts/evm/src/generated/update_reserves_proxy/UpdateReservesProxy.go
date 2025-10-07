// Code generated — DO NOT EDIT.

package update_reserves_proxy

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

var UpdateReservesProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_reserveManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"UnauthorizedWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"}],\"name\":\"UnauthorizedWorkflowOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserveManager\",\"outputs\":[{\"internalType\":\"contractIReserveManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
}

// Structs

// Contract Method Inputs
type OnReportInput struct {
	Metadata []byte
	Report   []byte
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
}

// Contract Method Outputs

// Errors
type UnauthorizedWorkflowName struct {
	WorkflowName [10]byte
}

type UnauthorizedWorkflowOwner struct {
	WorkflowOwner common.Address
}

// Events

// Main Binding Type for UpdateReservesProxy
type UpdateReservesProxy struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   UpdateReservesProxyCodec
}

type UpdateReservesProxyCodec interface {
	EncodeOnReportMethodCall(in OnReportInput) ([]byte, error)
	EncodeReserveManagerMethodCall() ([]byte, error)
	DecodeReserveManagerMethodOutput(data []byte) (common.Address, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
}

func NewUpdateReservesProxy(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*UpdateReservesProxy, error) {
	parsed, err := abi.JSON(strings.NewReader(UpdateReservesProxyMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &UpdateReservesProxy{
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

func NewCodec() (UpdateReservesProxyCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(UpdateReservesProxyMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeOnReportMethodCall(in OnReportInput) ([]byte, error) {
	return c.abi.Pack("onReport", in.Metadata, in.Report)
}

func (c *Codec) EncodeReserveManagerMethodCall() ([]byte, error) {
	return c.abi.Pack("reserveManager")
}

func (c *Codec) DecodeReserveManagerMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["reserveManager"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
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

func (c UpdateReservesProxy) ReserveManager(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeReserveManagerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: pb.NewBigIntFromInt(big.NewInt(rpc.FinalizedBlockNumber.Int64())),
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
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeReserveManagerMethodOutput(response.Data)
	})

}

func (c UpdateReservesProxy) WriteReport(
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

// DecodeUnauthorizedWorkflowNameError decodes a UnauthorizedWorkflowName error from revert data.
func (c *UpdateReservesProxy) DecodeUnauthorizedWorkflowNameError(data []byte) (*UnauthorizedWorkflowName, error) {
	args := c.ABI.Errors["UnauthorizedWorkflowName"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	workflowName, ok0 := values[0].([10]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for workflowName in UnauthorizedWorkflowName error")
	}

	return &UnauthorizedWorkflowName{
		WorkflowName: workflowName,
	}, nil
}

// Error implements the error interface for UnauthorizedWorkflowName.
func (e *UnauthorizedWorkflowName) Error() string {
	return fmt.Sprintf("UnauthorizedWorkflowName error: workflowName=%v;", e.WorkflowName)
}

// DecodeUnauthorizedWorkflowOwnerError decodes a UnauthorizedWorkflowOwner error from revert data.
func (c *UpdateReservesProxy) DecodeUnauthorizedWorkflowOwnerError(data []byte) (*UnauthorizedWorkflowOwner, error) {
	args := c.ABI.Errors["UnauthorizedWorkflowOwner"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	workflowOwner, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for workflowOwner in UnauthorizedWorkflowOwner error")
	}

	return &UnauthorizedWorkflowOwner{
		WorkflowOwner: workflowOwner,
	}, nil
}

// Error implements the error interface for UnauthorizedWorkflowOwner.
func (e *UnauthorizedWorkflowOwner) Error() string {
	return fmt.Sprintf("UnauthorizedWorkflowOwner error: workflowOwner=%v;", e.WorkflowOwner)
}

func (c *UpdateReservesProxy) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["UnauthorizedWorkflowName"].ID.Bytes()[:4]):
		return c.DecodeUnauthorizedWorkflowNameError(data)
	case common.Bytes2Hex(c.ABI.Errors["UnauthorizedWorkflowOwner"].ID.Bytes()[:4]):
		return c.DecodeUnauthorizedWorkflowOwnerError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}
