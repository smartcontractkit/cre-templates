// Code generated — DO NOT EDIT.

package mock_pool

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

var MockPoolMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"currentLiquidityRate\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getReserveData\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structDataTypes.ReserveDataLegacy\",\"components\":[{\"name\":\"configuration\",\"type\":\"tuple\",\"internalType\":\"structDataTypes.ReserveConfigurationMap\",\"components\":[{\"name\":\"data\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"liquidityIndex\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"currentLiquidityRate\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"variableBorrowIndex\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"currentVariableBorrowRate\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"currentStableBorrowRate\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"lastUpdateTimestamp\",\"type\":\"uint40\",\"internalType\":\"uint40\"},{\"name\":\"id\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"aTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stableDebtTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"variableDebtTokenAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"interestRateStrategyAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"accruedToTreasury\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"unbacked\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"isolationModeTotalDebt\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setCurrentLiquidityRate\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rate\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supply\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"onBehalfOf\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CurrentLiquidityRateUpdated\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"currentLiquidityRate\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Supply\",\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"onBehalfOf\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"referralCode\",\"type\":\"uint16\",\"indexed\":true,\"internalType\":\"uint16\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Withdraw\",\"inputs\":[{\"name\":\"reserve\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AmountZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoBalance\",\"inputs\":[]}]",
}

// Structs
type DataTypesReserveConfigurationMap struct {
	Data *big.Int
}

type DataTypesReserveDataLegacy struct {
	Configuration               DataTypesReserveConfigurationMap
	LiquidityIndex              *big.Int
	CurrentLiquidityRate        *big.Int
	VariableBorrowIndex         *big.Int
	CurrentVariableBorrowRate   *big.Int
	CurrentStableBorrowRate     *big.Int
	LastUpdateTimestamp         *big.Int
	Id                          uint16
	ATokenAddress               common.Address
	StableDebtTokenAddress      common.Address
	VariableDebtTokenAddress    common.Address
	InterestRateStrategyAddress common.Address
	AccruedToTreasury           *big.Int
	Unbacked                    *big.Int
	IsolationModeTotalDebt      *big.Int
}

// Contract Method Inputs
type BalanceOfInput struct {
	Arg0 common.Address
	Arg1 common.Address
}

type CurrentLiquidityRateInput struct {
	Arg0 common.Address
}

type GetReserveDataInput struct {
	Asset common.Address
}

type SetCurrentLiquidityRateInput struct {
	Asset common.Address
	Rate  *big.Int
}

type SupplyInput struct {
	Asset      common.Address
	Amount     *big.Int
	OnBehalfOf common.Address
	Arg3       uint16
}

type TransferOwnershipInput struct {
	To common.Address
}

type WithdrawInput struct {
	Asset  common.Address
	Amount *big.Int
	To     common.Address
}

// Contract Method Outputs

// Errors
type AmountZero struct {
}

type InsufficientBalance struct {
}

type NoBalance struct {
}

// Events
// The <Event> struct should be used as a filter (for log triggers).
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type CurrentLiquidityRateUpdated struct {
	Asset                common.Address
	CurrentLiquidityRate *big.Int
}

type CurrentLiquidityRateUpdatedDecoded struct {
	Asset                common.Address
	CurrentLiquidityRate *big.Int
}

type OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferRequestedDecoded struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferred struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferredDecoded struct {
	From common.Address
	To   common.Address
}

type Supply struct {
	Reserve      common.Address
	User         common.Address
	OnBehalfOf   common.Address
	Amount       *big.Int
	ReferralCode uint16
}

type SupplyDecoded struct {
	Reserve      common.Address
	User         common.Address
	OnBehalfOf   common.Address
	Amount       *big.Int
	ReferralCode uint16
}

type Withdraw struct {
	Reserve common.Address
	User    common.Address
	To      common.Address
	Amount  *big.Int
}

type WithdrawDecoded struct {
	Reserve common.Address
	User    common.Address
	To      common.Address
	Amount  *big.Int
}

// Main Binding Type for MockPool
type MockPool struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   MockPoolCodec
}

type MockPoolCodec interface {
	EncodeAcceptOwnershipMethodCall() ([]byte, error)
	EncodeBalanceOfMethodCall(in BalanceOfInput) ([]byte, error)
	DecodeBalanceOfMethodOutput(data []byte) (*big.Int, error)
	EncodeCurrentLiquidityRateMethodCall(in CurrentLiquidityRateInput) ([]byte, error)
	DecodeCurrentLiquidityRateMethodOutput(data []byte) (*big.Int, error)
	EncodeGetReserveDataMethodCall(in GetReserveDataInput) ([]byte, error)
	DecodeGetReserveDataMethodOutput(data []byte) (DataTypesReserveDataLegacy, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodeSetCurrentLiquidityRateMethodCall(in SetCurrentLiquidityRateInput) ([]byte, error)
	EncodeSupplyMethodCall(in SupplyInput) ([]byte, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	EncodeWithdrawMethodCall(in WithdrawInput) ([]byte, error)
	DecodeWithdrawMethodOutput(data []byte) (*big.Int, error)
	EncodeDataTypesReserveConfigurationMapStruct(in DataTypesReserveConfigurationMap) ([]byte, error)
	EncodeDataTypesReserveDataLegacyStruct(in DataTypesReserveDataLegacy) ([]byte, error)
	CurrentLiquidityRateUpdatedLogHash() []byte
	EncodeCurrentLiquidityRateUpdatedTopics(evt abi.Event, values []CurrentLiquidityRateUpdated) ([]*evm.TopicValues, error)
	DecodeCurrentLiquidityRateUpdated(log *evm.Log) (*CurrentLiquidityRateUpdatedDecoded, error)
	OwnershipTransferRequestedLogHash() []byte
	EncodeOwnershipTransferRequestedTopics(evt abi.Event, values []OwnershipTransferRequested) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferRequested(log *evm.Log) (*OwnershipTransferRequestedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferred) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
	SupplyLogHash() []byte
	EncodeSupplyTopics(evt abi.Event, values []Supply) ([]*evm.TopicValues, error)
	DecodeSupply(log *evm.Log) (*SupplyDecoded, error)
	WithdrawLogHash() []byte
	EncodeWithdrawTopics(evt abi.Event, values []Withdraw) ([]*evm.TopicValues, error)
	DecodeWithdraw(log *evm.Log) (*WithdrawDecoded, error)
}

func NewMockPool(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*MockPool, error) {
	parsed, err := abi.JSON(strings.NewReader(MockPoolMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &MockPool{
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

func NewCodec() (MockPoolCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(MockPoolMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeAcceptOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("acceptOwnership")
}

func (c *Codec) EncodeBalanceOfMethodCall(in BalanceOfInput) ([]byte, error) {
	return c.abi.Pack("balanceOf", in.Arg0, in.Arg1)
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

func (c *Codec) EncodeCurrentLiquidityRateMethodCall(in CurrentLiquidityRateInput) ([]byte, error) {
	return c.abi.Pack("currentLiquidityRate", in.Arg0)
}

func (c *Codec) DecodeCurrentLiquidityRateMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["currentLiquidityRate"].Outputs.Unpack(data)
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

func (c *Codec) EncodeGetReserveDataMethodCall(in GetReserveDataInput) ([]byte, error) {
	return c.abi.Pack("getReserveData", in.Asset)
}

func (c *Codec) DecodeGetReserveDataMethodOutput(data []byte) (DataTypesReserveDataLegacy, error) {
	vals, err := c.abi.Methods["getReserveData"].Outputs.Unpack(data)
	if err != nil {
		return *new(DataTypesReserveDataLegacy), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(DataTypesReserveDataLegacy), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result DataTypesReserveDataLegacy
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(DataTypesReserveDataLegacy), fmt.Errorf("failed to unmarshal to DataTypesReserveDataLegacy: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeOwnerMethodCall() ([]byte, error) {
	return c.abi.Pack("owner")
}

func (c *Codec) DecodeOwnerMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["owner"].Outputs.Unpack(data)
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

func (c *Codec) EncodeSetCurrentLiquidityRateMethodCall(in SetCurrentLiquidityRateInput) ([]byte, error) {
	return c.abi.Pack("setCurrentLiquidityRate", in.Asset, in.Rate)
}

func (c *Codec) EncodeSupplyMethodCall(in SupplyInput) ([]byte, error) {
	return c.abi.Pack("supply", in.Asset, in.Amount, in.OnBehalfOf, in.Arg3)
}

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.To)
}

func (c *Codec) EncodeWithdrawMethodCall(in WithdrawInput) ([]byte, error) {
	return c.abi.Pack("withdraw", in.Asset, in.Amount, in.To)
}

func (c *Codec) DecodeWithdrawMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["withdraw"].Outputs.Unpack(data)
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

func (c *Codec) EncodeDataTypesReserveConfigurationMapStruct(in DataTypesReserveConfigurationMap) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "data", Type: "uint256"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for DataTypesReserveConfigurationMap: %w", err)
	}
	args := abi.Arguments{
		{Name: "dataTypesReserveConfigurationMap", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeDataTypesReserveDataLegacyStruct(in DataTypesReserveDataLegacy) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "configuration", Type: "(uint256)"},
			{Name: "liquidityIndex", Type: "uint128"},
			{Name: "currentLiquidityRate", Type: "uint128"},
			{Name: "variableBorrowIndex", Type: "uint128"},
			{Name: "currentVariableBorrowRate", Type: "uint128"},
			{Name: "currentStableBorrowRate", Type: "uint128"},
			{Name: "lastUpdateTimestamp", Type: "uint40"},
			{Name: "id", Type: "uint16"},
			{Name: "aTokenAddress", Type: "address"},
			{Name: "stableDebtTokenAddress", Type: "address"},
			{Name: "variableDebtTokenAddress", Type: "address"},
			{Name: "interestRateStrategyAddress", Type: "address"},
			{Name: "accruedToTreasury", Type: "uint128"},
			{Name: "unbacked", Type: "uint128"},
			{Name: "isolationModeTotalDebt", Type: "uint128"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for DataTypesReserveDataLegacy: %w", err)
	}
	args := abi.Arguments{
		{Name: "dataTypesReserveDataLegacy", Type: tupleType},
	}

	return args.Pack(in)
}

func (c *Codec) CurrentLiquidityRateUpdatedLogHash() []byte {
	return c.abi.Events["CurrentLiquidityRateUpdated"].ID.Bytes()
}

func (c *Codec) EncodeCurrentLiquidityRateUpdatedTopics(
	evt abi.Event,
	values []CurrentLiquidityRateUpdated,
) ([]*evm.TopicValues, error) {
	var assetRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Asset)
		if err != nil {
			return nil, err
		}
		assetRule = append(assetRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		assetRule,
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

// DecodeCurrentLiquidityRateUpdated decodes a log into a CurrentLiquidityRateUpdated struct.
func (c *Codec) DecodeCurrentLiquidityRateUpdated(log *evm.Log) (*CurrentLiquidityRateUpdatedDecoded, error) {
	event := new(CurrentLiquidityRateUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "CurrentLiquidityRateUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["CurrentLiquidityRateUpdated"].Inputs {
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

func (c *Codec) OwnershipTransferRequestedLogHash() []byte {
	return c.abi.Events["OwnershipTransferRequested"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferRequestedTopics(
	evt abi.Event,
	values []OwnershipTransferRequested,
) ([]*evm.TopicValues, error) {
	var fromRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.From)
		if err != nil {
			return nil, err
		}
		fromRule = append(fromRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
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
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeOwnershipTransferRequested decodes a log into a OwnershipTransferRequested struct.
func (c *Codec) DecodeOwnershipTransferRequested(log *evm.Log) (*OwnershipTransferRequestedDecoded, error) {
	event := new(OwnershipTransferRequestedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferRequested", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferRequested"].Inputs {
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

func (c *Codec) OwnershipTransferredLogHash() []byte {
	return c.abi.Events["OwnershipTransferred"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferredTopics(
	evt abi.Event,
	values []OwnershipTransferred,
) ([]*evm.TopicValues, error) {
	var fromRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.From)
		if err != nil {
			return nil, err
		}
		fromRule = append(fromRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
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
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeOwnershipTransferred decodes a log into a OwnershipTransferred struct.
func (c *Codec) DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error) {
	event := new(OwnershipTransferredDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferred", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferred"].Inputs {
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

func (c *Codec) SupplyLogHash() []byte {
	return c.abi.Events["Supply"].ID.Bytes()
}

func (c *Codec) EncodeSupplyTopics(
	evt abi.Event,
	values []Supply,
) ([]*evm.TopicValues, error) {
	var reserveRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Reserve)
		if err != nil {
			return nil, err
		}
		reserveRule = append(reserveRule, fieldVal)
	}
	var onBehalfOfRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.OnBehalfOf)
		if err != nil {
			return nil, err
		}
		onBehalfOfRule = append(onBehalfOfRule, fieldVal)
	}
	var referralCodeRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[4], v.ReferralCode)
		if err != nil {
			return nil, err
		}
		referralCodeRule = append(referralCodeRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		reserveRule,
		onBehalfOfRule,
		referralCodeRule,
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

// DecodeSupply decodes a log into a Supply struct.
func (c *Codec) DecodeSupply(log *evm.Log) (*SupplyDecoded, error) {
	event := new(SupplyDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Supply", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Supply"].Inputs {
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

func (c *Codec) WithdrawLogHash() []byte {
	return c.abi.Events["Withdraw"].ID.Bytes()
}

func (c *Codec) EncodeWithdrawTopics(
	evt abi.Event,
	values []Withdraw,
) ([]*evm.TopicValues, error) {
	var reserveRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Reserve)
		if err != nil {
			return nil, err
		}
		reserveRule = append(reserveRule, fieldVal)
	}
	var userRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.User)
		if err != nil {
			return nil, err
		}
		userRule = append(userRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.To)
		if err != nil {
			return nil, err
		}
		toRule = append(toRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		reserveRule,
		userRule,
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
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeWithdraw decodes a log into a Withdraw struct.
func (c *Codec) DecodeWithdraw(log *evm.Log) (*WithdrawDecoded, error) {
	event := new(WithdrawDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Withdraw", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Withdraw"].Inputs {
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

func (c MockPool) BalanceOf(
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

func (c MockPool) CurrentLiquidityRate(
	runtime cre.Runtime,
	args CurrentLiquidityRateInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeCurrentLiquidityRateMethodCall(args)
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
		return c.Codec.DecodeCurrentLiquidityRateMethodOutput(response.Data)
	})

}

func (c MockPool) GetReserveData(
	runtime cre.Runtime,
	args GetReserveDataInput,
	blockNumber *big.Int,
) cre.Promise[DataTypesReserveDataLegacy] {
	calldata, err := c.Codec.EncodeGetReserveDataMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[DataTypesReserveDataLegacy](*new(DataTypesReserveDataLegacy), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (DataTypesReserveDataLegacy, error) {
		return c.Codec.DecodeGetReserveDataMethodOutput(response.Data)
	})

}

func (c MockPool) Owner(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeOwnerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeOwnerMethodOutput(response.Data)
	})

}

func (c MockPool) WriteReportFromDataTypesReserveConfigurationMap(
	runtime cre.Runtime,
	input DataTypesReserveConfigurationMap,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeDataTypesReserveConfigurationMapStruct(input)
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

func (c MockPool) WriteReportFromDataTypesReserveDataLegacy(
	runtime cre.Runtime,
	input DataTypesReserveDataLegacy,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeDataTypesReserveDataLegacyStruct(input)
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

func (c MockPool) WriteReport(
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

// DecodeAmountZeroError decodes a AmountZero error from revert data.
func (c *MockPool) DecodeAmountZeroError(data []byte) (*AmountZero, error) {
	args := c.ABI.Errors["AmountZero"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &AmountZero{}, nil
}

// Error implements the error interface for AmountZero.
func (e *AmountZero) Error() string {
	return fmt.Sprintf("AmountZero error:")
}

// DecodeInsufficientBalanceError decodes a InsufficientBalance error from revert data.
func (c *MockPool) DecodeInsufficientBalanceError(data []byte) (*InsufficientBalance, error) {
	args := c.ABI.Errors["InsufficientBalance"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InsufficientBalance{}, nil
}

// Error implements the error interface for InsufficientBalance.
func (e *InsufficientBalance) Error() string {
	return fmt.Sprintf("InsufficientBalance error:")
}

// DecodeNoBalanceError decodes a NoBalance error from revert data.
func (c *MockPool) DecodeNoBalanceError(data []byte) (*NoBalance, error) {
	args := c.ABI.Errors["NoBalance"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &NoBalance{}, nil
}

// Error implements the error interface for NoBalance.
func (e *NoBalance) Error() string {
	return fmt.Sprintf("NoBalance error:")
}

func (c *MockPool) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["AmountZero"].ID.Bytes()[:4]):
		return c.DecodeAmountZeroError(data)
	case common.Bytes2Hex(c.ABI.Errors["InsufficientBalance"].ID.Bytes()[:4]):
		return c.DecodeInsufficientBalanceError(data)
	case common.Bytes2Hex(c.ABI.Errors["NoBalance"].ID.Bytes()[:4]):
		return c.DecodeNoBalanceError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// CurrentLiquidityRateUpdatedTrigger wraps the raw log trigger and provides decoded CurrentLiquidityRateUpdatedDecoded data
type CurrentLiquidityRateUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *MockPool // Keep reference for decoding
}

// Adapt method that decodes the log into CurrentLiquidityRateUpdated data
func (t *CurrentLiquidityRateUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[CurrentLiquidityRateUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeCurrentLiquidityRateUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CurrentLiquidityRateUpdated log: %w", err)
	}

	return &bindings.DecodedLog[CurrentLiquidityRateUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *MockPool) LogTriggerCurrentLiquidityRateUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []CurrentLiquidityRateUpdated) (cre.Trigger[*evm.Log, *bindings.DecodedLog[CurrentLiquidityRateUpdatedDecoded]], error) {
	event := c.ABI.Events["CurrentLiquidityRateUpdated"]
	topics, err := c.Codec.EncodeCurrentLiquidityRateUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for CurrentLiquidityRateUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &CurrentLiquidityRateUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *MockPool) FilterLogsCurrentLiquidityRateUpdated(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.CurrentLiquidityRateUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// OwnershipTransferRequestedTrigger wraps the raw log trigger and provides decoded OwnershipTransferRequestedDecoded data
type OwnershipTransferRequestedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *MockPool // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferRequested data
func (t *OwnershipTransferRequestedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferRequestedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferRequested(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferRequested log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferRequestedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *MockPool) LogTriggerOwnershipTransferRequestedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferRequested) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferRequestedDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferRequested"]
	topics, err := c.Codec.EncodeOwnershipTransferRequestedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferRequested: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferRequestedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *MockPool) FilterLogsOwnershipTransferRequested(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferRequestedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// OwnershipTransferredTrigger wraps the raw log trigger and provides decoded OwnershipTransferredDecoded data
type OwnershipTransferredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *MockPool // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferred data
func (t *OwnershipTransferredTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferredDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferred(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferred log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferredDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *MockPool) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferred) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferred"]
	topics, err := c.Codec.EncodeOwnershipTransferredTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferred: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferredTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *MockPool) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferredLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// SupplyTrigger wraps the raw log trigger and provides decoded SupplyDecoded data
type SupplyTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *MockPool // Keep reference for decoding
}

// Adapt method that decodes the log into Supply data
func (t *SupplyTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[SupplyDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeSupply(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Supply log: %w", err)
	}

	return &bindings.DecodedLog[SupplyDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *MockPool) LogTriggerSupplyLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []Supply) (cre.Trigger[*evm.Log, *bindings.DecodedLog[SupplyDecoded]], error) {
	event := c.ABI.Events["Supply"]
	topics, err := c.Codec.EncodeSupplyTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Supply: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &SupplyTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *MockPool) FilterLogsSupply(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.SupplyLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// WithdrawTrigger wraps the raw log trigger and provides decoded WithdrawDecoded data
type WithdrawTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *MockPool // Keep reference for decoding
}

// Adapt method that decodes the log into Withdraw data
func (t *WithdrawTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[WithdrawDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeWithdraw(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Withdraw log: %w", err)
	}

	return &bindings.DecodedLog[WithdrawDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *MockPool) LogTriggerWithdrawLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []Withdraw) (cre.Trigger[*evm.Log, *bindings.DecodedLog[WithdrawDecoded]], error) {
	event := c.ABI.Events["Withdraw"]
	topics, err := c.Codec.EncodeWithdrawTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Withdraw: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &WithdrawTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *MockPool) FilterLogsWithdraw(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.WithdrawLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}
