// Code generated — DO NOT EDIT.

package bundle_aggregator_proxy

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

var BundleAggregatorProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"AggregatorNotProposed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previous\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"latest\",\"type\":\"address\"}],\"name\":\"AggregatorConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"current\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proposed\",\"type\":\"address\"}],\"name\":\"AggregatorProposed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"aggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bundleDecimals\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"decimals\",\"type\":\"uint8[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorAddress\",\"type\":\"address\"}],\"name\":\"confirmAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"aggregatorDescription\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestBundle\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"bundle\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestBundleTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"aggregatorAddress\",\"type\":\"address\"}],\"name\":\"proposeAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proposedAggregator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proposedAggregatorAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"aggregatorVersion\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Structs

// Contract Method Inputs
type ConfirmAggregatorInput struct {
	AggregatorAddress common.Address
}

type ProposeAggregatorInput struct {
	AggregatorAddress common.Address
}

type TransferOwnershipInput struct {
	To common.Address
}

// Contract Method Outputs

// Errors
type AggregatorNotProposed struct {
	Aggregator common.Address
}

// Events
// The <Event>Topics struct should be used as a filter (for log triggers).
// Note: It is only possible to filter on indexed fields.
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type AggregatorConfirmedTopics struct {
	Previous common.Address
	Latest   common.Address
}

type AggregatorConfirmedDecoded struct {
	Previous common.Address
	Latest   common.Address
}

type AggregatorProposedTopics struct {
	Current  common.Address
	Proposed common.Address
}

type AggregatorProposedDecoded struct {
	Current  common.Address
	Proposed common.Address
}

type OwnershipTransferRequestedTopics struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferRequestedDecoded struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferredTopics struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferredDecoded struct {
	From common.Address
	To   common.Address
}

// Main Binding Type for BundleAggregatorProxy
type BundleAggregatorProxy struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   BundleAggregatorProxyCodec
}

type BundleAggregatorProxyCodec interface {
	EncodeAcceptOwnershipMethodCall() ([]byte, error)
	EncodeAggregatorMethodCall() ([]byte, error)
	DecodeAggregatorMethodOutput(data []byte) (common.Address, error)
	EncodeBundleDecimalsMethodCall() ([]byte, error)
	DecodeBundleDecimalsMethodOutput(data []byte) ([]uint8, error)
	EncodeConfirmAggregatorMethodCall(in ConfirmAggregatorInput) ([]byte, error)
	EncodeDescriptionMethodCall() ([]byte, error)
	DecodeDescriptionMethodOutput(data []byte) (string, error)
	EncodeLatestBundleMethodCall() ([]byte, error)
	DecodeLatestBundleMethodOutput(data []byte) ([]byte, error)
	EncodeLatestBundleTimestampMethodCall() ([]byte, error)
	DecodeLatestBundleTimestampMethodOutput(data []byte) (*big.Int, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodeProposeAggregatorMethodCall(in ProposeAggregatorInput) ([]byte, error)
	EncodeProposedAggregatorMethodCall() ([]byte, error)
	DecodeProposedAggregatorMethodOutput(data []byte) (common.Address, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	EncodeTypeAndVersionMethodCall() ([]byte, error)
	DecodeTypeAndVersionMethodOutput(data []byte) (string, error)
	EncodeVersionMethodCall() ([]byte, error)
	DecodeVersionMethodOutput(data []byte) (*big.Int, error)
	AggregatorConfirmedLogHash() []byte
	EncodeAggregatorConfirmedTopics(evt abi.Event, values []AggregatorConfirmedTopics) ([]*evm.TopicValues, error)
	DecodeAggregatorConfirmed(log *evm.Log) (*AggregatorConfirmedDecoded, error)
	AggregatorProposedLogHash() []byte
	EncodeAggregatorProposedTopics(evt abi.Event, values []AggregatorProposedTopics) ([]*evm.TopicValues, error)
	DecodeAggregatorProposed(log *evm.Log) (*AggregatorProposedDecoded, error)
	OwnershipTransferRequestedLogHash() []byte
	EncodeOwnershipTransferRequestedTopics(evt abi.Event, values []OwnershipTransferRequestedTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferRequested(log *evm.Log) (*OwnershipTransferRequestedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferredTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
}

func NewBundleAggregatorProxy(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*BundleAggregatorProxy, error) {
	parsed, err := abi.JSON(strings.NewReader(BundleAggregatorProxyMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &BundleAggregatorProxy{
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

func NewCodec() (BundleAggregatorProxyCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(BundleAggregatorProxyMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeAcceptOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("acceptOwnership")
}

func (c *Codec) EncodeAggregatorMethodCall() ([]byte, error) {
	return c.abi.Pack("aggregator")
}

func (c *Codec) DecodeAggregatorMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["aggregator"].Outputs.Unpack(data)
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

func (c *Codec) EncodeBundleDecimalsMethodCall() ([]byte, error) {
	return c.abi.Pack("bundleDecimals")
}

func (c *Codec) DecodeBundleDecimalsMethodOutput(data []byte) ([]uint8, error) {
	vals, err := c.abi.Methods["bundleDecimals"].Outputs.Unpack(data)
	if err != nil {
		return *new([]uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]uint8), fmt.Errorf("failed to unmarshal to []uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeConfirmAggregatorMethodCall(in ConfirmAggregatorInput) ([]byte, error) {
	return c.abi.Pack("confirmAggregator", in.AggregatorAddress)
}

func (c *Codec) EncodeDescriptionMethodCall() ([]byte, error) {
	return c.abi.Pack("description")
}

func (c *Codec) DecodeDescriptionMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["description"].Outputs.Unpack(data)
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

func (c *Codec) EncodeLatestBundleMethodCall() ([]byte, error) {
	return c.abi.Pack("latestBundle")
}

func (c *Codec) DecodeLatestBundleMethodOutput(data []byte) ([]byte, error) {
	vals, err := c.abi.Methods["latestBundle"].Outputs.Unpack(data)
	if err != nil {
		return *new([]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]byte), fmt.Errorf("failed to unmarshal to []byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeLatestBundleTimestampMethodCall() ([]byte, error) {
	return c.abi.Pack("latestBundleTimestamp")
}

func (c *Codec) DecodeLatestBundleTimestampMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["latestBundleTimestamp"].Outputs.Unpack(data)
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

func (c *Codec) EncodeProposeAggregatorMethodCall(in ProposeAggregatorInput) ([]byte, error) {
	return c.abi.Pack("proposeAggregator", in.AggregatorAddress)
}

func (c *Codec) EncodeProposedAggregatorMethodCall() ([]byte, error) {
	return c.abi.Pack("proposedAggregator")
}

func (c *Codec) DecodeProposedAggregatorMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["proposedAggregator"].Outputs.Unpack(data)
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

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.To)
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

func (c *Codec) EncodeVersionMethodCall() ([]byte, error) {
	return c.abi.Pack("version")
}

func (c *Codec) DecodeVersionMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["version"].Outputs.Unpack(data)
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

func (c *Codec) AggregatorConfirmedLogHash() []byte {
	return c.abi.Events["AggregatorConfirmed"].ID.Bytes()
}

func (c *Codec) EncodeAggregatorConfirmedTopics(
	evt abi.Event,
	values []AggregatorConfirmedTopics,
) ([]*evm.TopicValues, error) {
	var previousRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Previous).IsZero() {
			previousRule = append(previousRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Previous)
		if err != nil {
			return nil, err
		}
		previousRule = append(previousRule, fieldVal)
	}
	var latestRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Latest).IsZero() {
			latestRule = append(latestRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Latest)
		if err != nil {
			return nil, err
		}
		latestRule = append(latestRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousRule,
		latestRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeAggregatorConfirmed decodes a log into a AggregatorConfirmed struct.
func (c *Codec) DecodeAggregatorConfirmed(log *evm.Log) (*AggregatorConfirmedDecoded, error) {
	event := new(AggregatorConfirmedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "AggregatorConfirmed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["AggregatorConfirmed"].Inputs {
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

func (c *Codec) AggregatorProposedLogHash() []byte {
	return c.abi.Events["AggregatorProposed"].ID.Bytes()
}

func (c *Codec) EncodeAggregatorProposedTopics(
	evt abi.Event,
	values []AggregatorProposedTopics,
) ([]*evm.TopicValues, error) {
	var currentRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Current).IsZero() {
			currentRule = append(currentRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Current)
		if err != nil {
			return nil, err
		}
		currentRule = append(currentRule, fieldVal)
	}
	var proposedRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Proposed).IsZero() {
			proposedRule = append(proposedRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Proposed)
		if err != nil {
			return nil, err
		}
		proposedRule = append(proposedRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		currentRule,
		proposedRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeAggregatorProposed decodes a log into a AggregatorProposed struct.
func (c *Codec) DecodeAggregatorProposed(log *evm.Log) (*AggregatorProposedDecoded, error) {
	event := new(AggregatorProposedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "AggregatorProposed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["AggregatorProposed"].Inputs {
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
	values []OwnershipTransferRequestedTopics,
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

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
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
	values []OwnershipTransferredTopics,
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

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
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

func (c BundleAggregatorProxy) Aggregator(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeAggregatorMethodCall()
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
		return c.Codec.DecodeAggregatorMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) BundleDecimals(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[]uint8] {
	calldata, err := c.Codec.EncodeBundleDecimalsMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[]uint8](*new([]uint8), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) ([]uint8, error) {
		return c.Codec.DecodeBundleDecimalsMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) Description(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeDescriptionMethodCall()
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
		return c.Codec.DecodeDescriptionMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) LatestBundle(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[]byte] {
	calldata, err := c.Codec.EncodeLatestBundleMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[]byte](*new([]byte), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) ([]byte, error) {
		return c.Codec.DecodeLatestBundleMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) LatestBundleTimestamp(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLatestBundleTimestampMethodCall()
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
		return c.Codec.DecodeLatestBundleTimestampMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) Owner(
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

func (c BundleAggregatorProxy) ProposedAggregator(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeProposedAggregatorMethodCall()
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
		return c.Codec.DecodeProposedAggregatorMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) TypeAndVersion(
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

func (c BundleAggregatorProxy) Version(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeVersionMethodCall()
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
		return c.Codec.DecodeVersionMethodOutput(response.Data)
	})

}

func (c BundleAggregatorProxy) WriteReport(
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

// DecodeAggregatorNotProposedError decodes a AggregatorNotProposed error from revert data.
func (c *BundleAggregatorProxy) DecodeAggregatorNotProposedError(data []byte) (*AggregatorNotProposed, error) {
	args := c.ABI.Errors["AggregatorNotProposed"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	aggregator, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for aggregator in AggregatorNotProposed error")
	}

	return &AggregatorNotProposed{
		Aggregator: aggregator,
	}, nil
}

// Error implements the error interface for AggregatorNotProposed.
func (e *AggregatorNotProposed) Error() string {
	return fmt.Sprintf("AggregatorNotProposed error: aggregator=%v;", e.Aggregator)
}

func (c *BundleAggregatorProxy) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["AggregatorNotProposed"].ID.Bytes()[:4]):
		return c.DecodeAggregatorNotProposedError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// AggregatorConfirmedTrigger wraps the raw log trigger and provides decoded AggregatorConfirmedDecoded data
type AggregatorConfirmedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                        // Embed the raw trigger
	contract                        *BundleAggregatorProxy // Keep reference for decoding
}

// Adapt method that decodes the log into AggregatorConfirmed data
func (t *AggregatorConfirmedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[AggregatorConfirmedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeAggregatorConfirmed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode AggregatorConfirmed log: %w", err)
	}

	return &bindings.DecodedLog[AggregatorConfirmedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *BundleAggregatorProxy) LogTriggerAggregatorConfirmedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []AggregatorConfirmedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[AggregatorConfirmedDecoded]], error) {
	event := c.ABI.Events["AggregatorConfirmed"]
	topics, err := c.Codec.EncodeAggregatorConfirmedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for AggregatorConfirmed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &AggregatorConfirmedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *BundleAggregatorProxy) FilterLogsAggregatorConfirmed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.AggregatorConfirmedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// AggregatorProposedTrigger wraps the raw log trigger and provides decoded AggregatorProposedDecoded data
type AggregatorProposedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                        // Embed the raw trigger
	contract                        *BundleAggregatorProxy // Keep reference for decoding
}

// Adapt method that decodes the log into AggregatorProposed data
func (t *AggregatorProposedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[AggregatorProposedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeAggregatorProposed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode AggregatorProposed log: %w", err)
	}

	return &bindings.DecodedLog[AggregatorProposedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *BundleAggregatorProxy) LogTriggerAggregatorProposedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []AggregatorProposedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[AggregatorProposedDecoded]], error) {
	event := c.ABI.Events["AggregatorProposed"]
	topics, err := c.Codec.EncodeAggregatorProposedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for AggregatorProposed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &AggregatorProposedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *BundleAggregatorProxy) FilterLogsAggregatorProposed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.AggregatorProposedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// OwnershipTransferRequestedTrigger wraps the raw log trigger and provides decoded OwnershipTransferRequestedDecoded data
type OwnershipTransferRequestedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                        // Embed the raw trigger
	contract                        *BundleAggregatorProxy // Keep reference for decoding
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

func (c *BundleAggregatorProxy) LogTriggerOwnershipTransferRequestedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferRequestedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferRequestedDecoded]], error) {
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

func (c *BundleAggregatorProxy) FilterLogsOwnershipTransferRequested(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
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
	}), nil
}

// OwnershipTransferredTrigger wraps the raw log trigger and provides decoded OwnershipTransferredDecoded data
type OwnershipTransferredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                        // Embed the raw trigger
	contract                        *BundleAggregatorProxy // Keep reference for decoding
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

func (c *BundleAggregatorProxy) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferredTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
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

func (c *BundleAggregatorProxy) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
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
	}), nil
}
