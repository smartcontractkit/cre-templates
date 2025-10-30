// Code generated — DO NOT EDIT.

package protocol_smart_wallet

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

var ProtocolSmartWalletMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_keystoneForwarders\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_allowedWorkflowOwners\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_pool\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_router\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_link\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowedCcipSenders\",\"inputs\":[{\"name\":\"chainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowedKeystoneForwarders\",\"inputs\":[{\"name\":\"keystoneForwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"allowedWorkflowOwners\",\"inputs\":[{\"name\":\"workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ccipReceive\",\"inputs\":[{\"name\":\"message\",\"type\":\"tuple\",\"internalType\":\"structClient.Any2EVMMessage\",\"components\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"destTokenAmounts\",\"type\":\"tuple[]\",\"internalType\":\"structClient.EVMTokenAmount[]\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositToPool\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getPoolAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRouter\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"metadata\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pool\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractMockPool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeKeystoneForwarder\",\"inputs\":[{\"name\":\"_keystoneForwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeSenderForSourceChain\",\"inputs\":[{\"name\":\"_sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeWorkflowOwner\",\"inputs\":[{\"name\":\"_workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setKeystoneForwarder\",\"inputs\":[{\"name\":\"_keystoneForwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPool\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSenderForSourceChain\",\"inputs\":[{\"name\":\"_sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"_sender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setWorkflowOwner\",\"inputs\":[{\"name\":\"_workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFromPool\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFromPoolAndDepositCrossChain\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structProtocolSmartWallet.RebalanceParams\",\"components\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destinationProtocolSmartWallet\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"KeystoneForwarderRemoved\",\"inputs\":[{\"name\":\"keystoneForwarder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"KeystoneForwarderSet\",\"inputs\":[{\"name\":\"keystoneForwarder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageReceived\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structClient.EVMTokenAmount\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MessageSent\",\"inputs\":[{\"name\":\"messageId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"receiver\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"tokenAmount\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structClient.EVMTokenAmount\",\"components\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"fees\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferRequested\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReportReceived\",\"inputs\":[{\"name\":\"workflowOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"workflowName\",\"type\":\"bytes10\",\"indexed\":true,\"internalType\":\"bytes10\"},{\"name\":\"params\",\"type\":\"tuple\",\"indexed\":false,\"internalType\":\"structProtocolSmartWallet.RebalanceParams\",\"components\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"destinationProtocolSmartWallet\",\"type\":\"address\",\"internalType\":\"address\"}]}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SenderForSourceChainRemoved\",\"inputs\":[{\"name\":\"_sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"_sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SenderForSourceChainSet\",\"inputs\":[{\"name\":\"_sourceChainSelector\",\"type\":\"uint64\",\"indexed\":true,\"internalType\":\"uint64\"},{\"name\":\"_sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WorkflowOwnerRemoved\",\"inputs\":[{\"name\":\"workflowOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"WorkflowOwnerSet\",\"inputs\":[{\"name\":\"workflowOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientFeeTokenAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientTokenAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidKeystoneForwarder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRouter\",\"inputs\":[{\"name\":\"router\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"InvalidSenderAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSourceChain\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidWorkflowOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedTokenAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MustBeKeystoneForwarder\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoSenderOnSourceChain\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"UnauthorizedWorkflowOwner\",\"inputs\":[{\"name\":\"workflowOwner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"WrongSenderForSourceChain\",\"inputs\":[{\"name\":\"sourceChainSelector\",\"type\":\"uint64\",\"internalType\":\"uint64\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}]",
}

// Structs
type ClientAny2EVMMessage struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              []byte
	Data                []byte
	DestTokenAmounts    []ClientEVMTokenAmount
}

type ClientEVMTokenAmount struct {
	Token  common.Address
	Amount *big.Int
}

type RebalanceParams struct {
	Asset                          common.Address
	Amount                         *big.Int
	DestinationChainSelector       uint64
	DestinationProtocolSmartWallet common.Address
}

// Contract Method Inputs
type AllowedCcipSendersInput struct {
	ChainSelector uint64
}

type AllowedKeystoneForwardersInput struct {
	KeystoneForwarder common.Address
}

type AllowedWorkflowOwnersInput struct {
	WorkflowOwner common.Address
}

type CcipReceiveInput struct {
	Message ClientAny2EVMMessage
}

type DepositToPoolInput struct {
	Asset  common.Address
	Amount *big.Int
}

type OnReportInput struct {
	Metadata []byte
	Report   []byte
}

type RemoveKeystoneForwarderInput struct {
	KeystoneForwarder common.Address
}

type RemoveSenderForSourceChainInput struct {
	SourceChainSelector uint64
}

type RemoveWorkflowOwnerInput struct {
	WorkflowOwner common.Address
}

type SetKeystoneForwarderInput struct {
	KeystoneForwarder common.Address
}

type SetPoolInput struct {
	Pool common.Address
}

type SetSenderForSourceChainInput struct {
	SourceChainSelector uint64
	Sender              common.Address
}

type SetWorkflowOwnerInput struct {
	WorkflowOwner common.Address
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
}

type TransferOwnershipInput struct {
	To common.Address
}

type WithdrawFromPoolInput struct {
	Asset  common.Address
	Amount *big.Int
	To     common.Address
}

type WithdrawFromPoolAndDepositCrossChainInput struct {
	Params RebalanceParams
}

// Contract Method Outputs

// Errors
type InsufficientFeeTokenAmount struct {
}

type InsufficientTokenAmount struct {
}

type InvalidKeystoneForwarder struct {
}

type InvalidRouter struct {
	Router common.Address
}

type InvalidSenderAddress struct {
}

type InvalidSourceChain struct {
}

type InvalidWorkflowOwner struct {
}

type MismatchedTokenAmount struct {
}

type MustBeKeystoneForwarder struct {
}

type NoSenderOnSourceChain struct {
	SourceChainSelector uint64
}

type UnauthorizedWorkflowOwner struct {
	WorkflowOwner common.Address
}

type WrongSenderForSourceChain struct {
	SourceChainSelector uint64
}

type ZeroAddress struct {
	Index *big.Int
}

// Events
// The <Event> struct should be used as a filter (for log triggers).
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type KeystoneForwarderRemoved struct {
	KeystoneForwarder common.Address
}

type KeystoneForwarderRemovedDecoded struct {
	KeystoneForwarder common.Address
}

type KeystoneForwarderSet struct {
	KeystoneForwarder common.Address
}

type KeystoneForwarderSetDecoded struct {
	KeystoneForwarder common.Address
}

type MessageReceived struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              common.Address
	TokenAmount         ClientEVMTokenAmount
}

type MessageReceivedDecoded struct {
	MessageId           [32]byte
	SourceChainSelector uint64
	Sender              common.Address
	TokenAmount         ClientEVMTokenAmount
}

type MessageSent struct {
	MessageId                [32]byte
	DestinationChainSelector uint64
	Receiver                 common.Address
	TokenAmount              ClientEVMTokenAmount
	Fees                     *big.Int
}

type MessageSentDecoded struct {
	MessageId                [32]byte
	DestinationChainSelector uint64
	Receiver                 common.Address
	TokenAmount              ClientEVMTokenAmount
	Fees                     *big.Int
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

type ReportReceived struct {
	WorkflowOwner common.Address
	WorkflowName  [10]byte
	Params        RebalanceParams
}

type ReportReceivedDecoded struct {
	WorkflowOwner common.Address
	WorkflowName  [10]byte
	Params        RebalanceParams
}

type SenderForSourceChainRemoved struct {
	SourceChainSelector uint64
	Sender              common.Address
}

type SenderForSourceChainRemovedDecoded struct {
	SourceChainSelector uint64
	Sender              common.Address
}

type SenderForSourceChainSet struct {
	SourceChainSelector uint64
	Sender              common.Address
}

type SenderForSourceChainSetDecoded struct {
	SourceChainSelector uint64
	Sender              common.Address
}

type WorkflowOwnerRemoved struct {
	WorkflowOwner common.Address
}

type WorkflowOwnerRemovedDecoded struct {
	WorkflowOwner common.Address
}

type WorkflowOwnerSet struct {
	WorkflowOwner common.Address
}

type WorkflowOwnerSetDecoded struct {
	WorkflowOwner common.Address
}

// Main Binding Type for ProtocolSmartWallet
type ProtocolSmartWallet struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   ProtocolSmartWalletCodec
}

type ProtocolSmartWalletCodec interface {
	EncodeAcceptOwnershipMethodCall() ([]byte, error)
	EncodeAllowedCcipSendersMethodCall(in AllowedCcipSendersInput) ([]byte, error)
	DecodeAllowedCcipSendersMethodOutput(data []byte) (common.Address, error)
	EncodeAllowedKeystoneForwardersMethodCall(in AllowedKeystoneForwardersInput) ([]byte, error)
	DecodeAllowedKeystoneForwardersMethodOutput(data []byte) (bool, error)
	EncodeAllowedWorkflowOwnersMethodCall(in AllowedWorkflowOwnersInput) ([]byte, error)
	DecodeAllowedWorkflowOwnersMethodOutput(data []byte) (bool, error)
	EncodeCcipReceiveMethodCall(in CcipReceiveInput) ([]byte, error)
	EncodeDepositToPoolMethodCall(in DepositToPoolInput) ([]byte, error)
	EncodeGetPoolAddressMethodCall() ([]byte, error)
	DecodeGetPoolAddressMethodOutput(data []byte) (common.Address, error)
	EncodeGetRouterMethodCall() ([]byte, error)
	DecodeGetRouterMethodOutput(data []byte) (common.Address, error)
	EncodeOnReportMethodCall(in OnReportInput) ([]byte, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodePoolMethodCall() ([]byte, error)
	DecodePoolMethodOutput(data []byte) (common.Address, error)
	EncodeRemoveKeystoneForwarderMethodCall(in RemoveKeystoneForwarderInput) ([]byte, error)
	EncodeRemoveSenderForSourceChainMethodCall(in RemoveSenderForSourceChainInput) ([]byte, error)
	EncodeRemoveWorkflowOwnerMethodCall(in RemoveWorkflowOwnerInput) ([]byte, error)
	EncodeSetKeystoneForwarderMethodCall(in SetKeystoneForwarderInput) ([]byte, error)
	EncodeSetPoolMethodCall(in SetPoolInput) ([]byte, error)
	EncodeSetSenderForSourceChainMethodCall(in SetSenderForSourceChainInput) ([]byte, error)
	EncodeSetWorkflowOwnerMethodCall(in SetWorkflowOwnerInput) ([]byte, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	EncodeWithdrawFromPoolMethodCall(in WithdrawFromPoolInput) ([]byte, error)
	EncodeWithdrawFromPoolAndDepositCrossChainMethodCall(in WithdrawFromPoolAndDepositCrossChainInput) ([]byte, error)
	EncodeClientAny2EVMMessageStruct(in ClientAny2EVMMessage) ([]byte, error)
	EncodeClientEVMTokenAmountStruct(in ClientEVMTokenAmount) ([]byte, error)
	EncodeRebalanceParamsStruct(in RebalanceParams) ([]byte, error)
	KeystoneForwarderRemovedLogHash() []byte
	EncodeKeystoneForwarderRemovedTopics(evt abi.Event, values []KeystoneForwarderRemoved) ([]*evm.TopicValues, error)
	DecodeKeystoneForwarderRemoved(log *evm.Log) (*KeystoneForwarderRemovedDecoded, error)
	KeystoneForwarderSetLogHash() []byte
	EncodeKeystoneForwarderSetTopics(evt abi.Event, values []KeystoneForwarderSet) ([]*evm.TopicValues, error)
	DecodeKeystoneForwarderSet(log *evm.Log) (*KeystoneForwarderSetDecoded, error)
	MessageReceivedLogHash() []byte
	EncodeMessageReceivedTopics(evt abi.Event, values []MessageReceived) ([]*evm.TopicValues, error)
	DecodeMessageReceived(log *evm.Log) (*MessageReceivedDecoded, error)
	MessageSentLogHash() []byte
	EncodeMessageSentTopics(evt abi.Event, values []MessageSent) ([]*evm.TopicValues, error)
	DecodeMessageSent(log *evm.Log) (*MessageSentDecoded, error)
	OwnershipTransferRequestedLogHash() []byte
	EncodeOwnershipTransferRequestedTopics(evt abi.Event, values []OwnershipTransferRequested) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferRequested(log *evm.Log) (*OwnershipTransferRequestedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferred) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
	ReportReceivedLogHash() []byte
	EncodeReportReceivedTopics(evt abi.Event, values []ReportReceived) ([]*evm.TopicValues, error)
	DecodeReportReceived(log *evm.Log) (*ReportReceivedDecoded, error)
	SenderForSourceChainRemovedLogHash() []byte
	EncodeSenderForSourceChainRemovedTopics(evt abi.Event, values []SenderForSourceChainRemoved) ([]*evm.TopicValues, error)
	DecodeSenderForSourceChainRemoved(log *evm.Log) (*SenderForSourceChainRemovedDecoded, error)
	SenderForSourceChainSetLogHash() []byte
	EncodeSenderForSourceChainSetTopics(evt abi.Event, values []SenderForSourceChainSet) ([]*evm.TopicValues, error)
	DecodeSenderForSourceChainSet(log *evm.Log) (*SenderForSourceChainSetDecoded, error)
	WorkflowOwnerRemovedLogHash() []byte
	EncodeWorkflowOwnerRemovedTopics(evt abi.Event, values []WorkflowOwnerRemoved) ([]*evm.TopicValues, error)
	DecodeWorkflowOwnerRemoved(log *evm.Log) (*WorkflowOwnerRemovedDecoded, error)
	WorkflowOwnerSetLogHash() []byte
	EncodeWorkflowOwnerSetTopics(evt abi.Event, values []WorkflowOwnerSet) ([]*evm.TopicValues, error)
	DecodeWorkflowOwnerSet(log *evm.Log) (*WorkflowOwnerSetDecoded, error)
}

func NewProtocolSmartWallet(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*ProtocolSmartWallet, error) {
	parsed, err := abi.JSON(strings.NewReader(ProtocolSmartWalletMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &ProtocolSmartWallet{
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

func NewCodec() (ProtocolSmartWalletCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(ProtocolSmartWalletMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeAcceptOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("acceptOwnership")
}

func (c *Codec) EncodeAllowedCcipSendersMethodCall(in AllowedCcipSendersInput) ([]byte, error) {
	return c.abi.Pack("allowedCcipSenders", in.ChainSelector)
}

func (c *Codec) DecodeAllowedCcipSendersMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["allowedCcipSenders"].Outputs.Unpack(data)
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

func (c *Codec) EncodeAllowedKeystoneForwardersMethodCall(in AllowedKeystoneForwardersInput) ([]byte, error) {
	return c.abi.Pack("allowedKeystoneForwarders", in.KeystoneForwarder)
}

func (c *Codec) DecodeAllowedKeystoneForwardersMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["allowedKeystoneForwarders"].Outputs.Unpack(data)
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

func (c *Codec) EncodeAllowedWorkflowOwnersMethodCall(in AllowedWorkflowOwnersInput) ([]byte, error) {
	return c.abi.Pack("allowedWorkflowOwners", in.WorkflowOwner)
}

func (c *Codec) DecodeAllowedWorkflowOwnersMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["allowedWorkflowOwners"].Outputs.Unpack(data)
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

func (c *Codec) EncodeCcipReceiveMethodCall(in CcipReceiveInput) ([]byte, error) {
	return c.abi.Pack("ccipReceive", in.Message)
}

func (c *Codec) EncodeDepositToPoolMethodCall(in DepositToPoolInput) ([]byte, error) {
	return c.abi.Pack("depositToPool", in.Asset, in.Amount)
}

func (c *Codec) EncodeGetPoolAddressMethodCall() ([]byte, error) {
	return c.abi.Pack("getPoolAddress")
}

func (c *Codec) DecodeGetPoolAddressMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["getPoolAddress"].Outputs.Unpack(data)
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

func (c *Codec) EncodeGetRouterMethodCall() ([]byte, error) {
	return c.abi.Pack("getRouter")
}

func (c *Codec) DecodeGetRouterMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["getRouter"].Outputs.Unpack(data)
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

func (c *Codec) EncodeOnReportMethodCall(in OnReportInput) ([]byte, error) {
	return c.abi.Pack("onReport", in.Metadata, in.Report)
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

func (c *Codec) EncodePoolMethodCall() ([]byte, error) {
	return c.abi.Pack("pool")
}

func (c *Codec) DecodePoolMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["pool"].Outputs.Unpack(data)
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

func (c *Codec) EncodeRemoveKeystoneForwarderMethodCall(in RemoveKeystoneForwarderInput) ([]byte, error) {
	return c.abi.Pack("removeKeystoneForwarder", in.KeystoneForwarder)
}

func (c *Codec) EncodeRemoveSenderForSourceChainMethodCall(in RemoveSenderForSourceChainInput) ([]byte, error) {
	return c.abi.Pack("removeSenderForSourceChain", in.SourceChainSelector)
}

func (c *Codec) EncodeRemoveWorkflowOwnerMethodCall(in RemoveWorkflowOwnerInput) ([]byte, error) {
	return c.abi.Pack("removeWorkflowOwner", in.WorkflowOwner)
}

func (c *Codec) EncodeSetKeystoneForwarderMethodCall(in SetKeystoneForwarderInput) ([]byte, error) {
	return c.abi.Pack("setKeystoneForwarder", in.KeystoneForwarder)
}

func (c *Codec) EncodeSetPoolMethodCall(in SetPoolInput) ([]byte, error) {
	return c.abi.Pack("setPool", in.Pool)
}

func (c *Codec) EncodeSetSenderForSourceChainMethodCall(in SetSenderForSourceChainInput) ([]byte, error) {
	return c.abi.Pack("setSenderForSourceChain", in.SourceChainSelector, in.Sender)
}

func (c *Codec) EncodeSetWorkflowOwnerMethodCall(in SetWorkflowOwnerInput) ([]byte, error) {
	return c.abi.Pack("setWorkflowOwner", in.WorkflowOwner)
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

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.To)
}

func (c *Codec) EncodeWithdrawFromPoolMethodCall(in WithdrawFromPoolInput) ([]byte, error) {
	return c.abi.Pack("withdrawFromPool", in.Asset, in.Amount, in.To)
}

func (c *Codec) EncodeWithdrawFromPoolAndDepositCrossChainMethodCall(in WithdrawFromPoolAndDepositCrossChainInput) ([]byte, error) {
	return c.abi.Pack("withdrawFromPoolAndDepositCrossChain", in.Params)
}

func (c *Codec) EncodeClientAny2EVMMessageStruct(in ClientAny2EVMMessage) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "messageId", Type: "bytes32"},
			{Name: "sourceChainSelector", Type: "uint64"},
			{Name: "sender", Type: "bytes"},
			{Name: "data", Type: "bytes"},
			{Name: "destTokenAmounts", Type: "(address,uint256)[]"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for ClientAny2EVMMessage: %w", err)
	}
	args := abi.Arguments{
		{Name: "clientAny2EVMMessage", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeClientEVMTokenAmountStruct(in ClientEVMTokenAmount) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "token", Type: "address"},
			{Name: "amount", Type: "uint256"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for ClientEVMTokenAmount: %w", err)
	}
	args := abi.Arguments{
		{Name: "clientEVMTokenAmount", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeRebalanceParamsStruct(in RebalanceParams) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "asset", Type: "address"},
			{Name: "amount", Type: "uint256"},
			{Name: "destinationChainSelector", Type: "uint64"},
			{Name: "destinationProtocolSmartWallet", Type: "address"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for RebalanceParams: %w", err)
	}
	args := abi.Arguments{
		{Name: "rebalanceParams", Type: tupleType},
	}

	return args.Pack(in)
}

func (c *Codec) KeystoneForwarderRemovedLogHash() []byte {
	return c.abi.Events["KeystoneForwarderRemoved"].ID.Bytes()
}

func (c *Codec) EncodeKeystoneForwarderRemovedTopics(
	evt abi.Event,
	values []KeystoneForwarderRemoved,
) ([]*evm.TopicValues, error) {
	var keystoneForwarderRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.KeystoneForwarder)
		if err != nil {
			return nil, err
		}
		keystoneForwarderRule = append(keystoneForwarderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		keystoneForwarderRule,
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

// DecodeKeystoneForwarderRemoved decodes a log into a KeystoneForwarderRemoved struct.
func (c *Codec) DecodeKeystoneForwarderRemoved(log *evm.Log) (*KeystoneForwarderRemovedDecoded, error) {
	event := new(KeystoneForwarderRemovedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "KeystoneForwarderRemoved", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["KeystoneForwarderRemoved"].Inputs {
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

func (c *Codec) KeystoneForwarderSetLogHash() []byte {
	return c.abi.Events["KeystoneForwarderSet"].ID.Bytes()
}

func (c *Codec) EncodeKeystoneForwarderSetTopics(
	evt abi.Event,
	values []KeystoneForwarderSet,
) ([]*evm.TopicValues, error) {
	var keystoneForwarderRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.KeystoneForwarder)
		if err != nil {
			return nil, err
		}
		keystoneForwarderRule = append(keystoneForwarderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		keystoneForwarderRule,
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

// DecodeKeystoneForwarderSet decodes a log into a KeystoneForwarderSet struct.
func (c *Codec) DecodeKeystoneForwarderSet(log *evm.Log) (*KeystoneForwarderSetDecoded, error) {
	event := new(KeystoneForwarderSetDecoded)
	if err := c.abi.UnpackIntoInterface(event, "KeystoneForwarderSet", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["KeystoneForwarderSet"].Inputs {
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

func (c *Codec) MessageReceivedLogHash() []byte {
	return c.abi.Events["MessageReceived"].ID.Bytes()
}

func (c *Codec) EncodeMessageReceivedTopics(
	evt abi.Event,
	values []MessageReceived,
) ([]*evm.TopicValues, error) {
	var messageIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.MessageId)
		if err != nil {
			return nil, err
		}
		messageIdRule = append(messageIdRule, fieldVal)
	}
	var sourceChainSelectorRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.SourceChainSelector)
		if err != nil {
			return nil, err
		}
		sourceChainSelectorRule = append(sourceChainSelectorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		messageIdRule,
		sourceChainSelectorRule,
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

// DecodeMessageReceived decodes a log into a MessageReceived struct.
func (c *Codec) DecodeMessageReceived(log *evm.Log) (*MessageReceivedDecoded, error) {
	event := new(MessageReceivedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "MessageReceived", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["MessageReceived"].Inputs {
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

func (c *Codec) MessageSentLogHash() []byte {
	return c.abi.Events["MessageSent"].ID.Bytes()
}

func (c *Codec) EncodeMessageSentTopics(
	evt abi.Event,
	values []MessageSent,
) ([]*evm.TopicValues, error) {
	var messageIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.MessageId)
		if err != nil {
			return nil, err
		}
		messageIdRule = append(messageIdRule, fieldVal)
	}
	var destinationChainSelectorRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.DestinationChainSelector)
		if err != nil {
			return nil, err
		}
		destinationChainSelectorRule = append(destinationChainSelectorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		messageIdRule,
		destinationChainSelectorRule,
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

// DecodeMessageSent decodes a log into a MessageSent struct.
func (c *Codec) DecodeMessageSent(log *evm.Log) (*MessageSentDecoded, error) {
	event := new(MessageSentDecoded)
	if err := c.abi.UnpackIntoInterface(event, "MessageSent", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["MessageSent"].Inputs {
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

func (c *Codec) ReportReceivedLogHash() []byte {
	return c.abi.Events["ReportReceived"].ID.Bytes()
}

func (c *Codec) EncodeReportReceivedTopics(
	evt abi.Event,
	values []ReportReceived,
) ([]*evm.TopicValues, error) {
	var workflowOwnerRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.WorkflowOwner)
		if err != nil {
			return nil, err
		}
		workflowOwnerRule = append(workflowOwnerRule, fieldVal)
	}
	var workflowNameRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.WorkflowName)
		if err != nil {
			return nil, err
		}
		workflowNameRule = append(workflowNameRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		workflowOwnerRule,
		workflowNameRule,
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

// DecodeReportReceived decodes a log into a ReportReceived struct.
func (c *Codec) DecodeReportReceived(log *evm.Log) (*ReportReceivedDecoded, error) {
	event := new(ReportReceivedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ReportReceived", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ReportReceived"].Inputs {
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

func (c *Codec) SenderForSourceChainRemovedLogHash() []byte {
	return c.abi.Events["SenderForSourceChainRemoved"].ID.Bytes()
}

func (c *Codec) EncodeSenderForSourceChainRemovedTopics(
	evt abi.Event,
	values []SenderForSourceChainRemoved,
) ([]*evm.TopicValues, error) {
	var sourceChainSelectorRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.SourceChainSelector)
		if err != nil {
			return nil, err
		}
		sourceChainSelectorRule = append(sourceChainSelectorRule, fieldVal)
	}
	var senderRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Sender)
		if err != nil {
			return nil, err
		}
		senderRule = append(senderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		sourceChainSelectorRule,
		senderRule,
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

// DecodeSenderForSourceChainRemoved decodes a log into a SenderForSourceChainRemoved struct.
func (c *Codec) DecodeSenderForSourceChainRemoved(log *evm.Log) (*SenderForSourceChainRemovedDecoded, error) {
	event := new(SenderForSourceChainRemovedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "SenderForSourceChainRemoved", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["SenderForSourceChainRemoved"].Inputs {
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

func (c *Codec) SenderForSourceChainSetLogHash() []byte {
	return c.abi.Events["SenderForSourceChainSet"].ID.Bytes()
}

func (c *Codec) EncodeSenderForSourceChainSetTopics(
	evt abi.Event,
	values []SenderForSourceChainSet,
) ([]*evm.TopicValues, error) {
	var sourceChainSelectorRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.SourceChainSelector)
		if err != nil {
			return nil, err
		}
		sourceChainSelectorRule = append(sourceChainSelectorRule, fieldVal)
	}
	var senderRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Sender)
		if err != nil {
			return nil, err
		}
		senderRule = append(senderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		sourceChainSelectorRule,
		senderRule,
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

// DecodeSenderForSourceChainSet decodes a log into a SenderForSourceChainSet struct.
func (c *Codec) DecodeSenderForSourceChainSet(log *evm.Log) (*SenderForSourceChainSetDecoded, error) {
	event := new(SenderForSourceChainSetDecoded)
	if err := c.abi.UnpackIntoInterface(event, "SenderForSourceChainSet", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["SenderForSourceChainSet"].Inputs {
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

func (c *Codec) WorkflowOwnerRemovedLogHash() []byte {
	return c.abi.Events["WorkflowOwnerRemoved"].ID.Bytes()
}

func (c *Codec) EncodeWorkflowOwnerRemovedTopics(
	evt abi.Event,
	values []WorkflowOwnerRemoved,
) ([]*evm.TopicValues, error) {
	var workflowOwnerRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.WorkflowOwner)
		if err != nil {
			return nil, err
		}
		workflowOwnerRule = append(workflowOwnerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		workflowOwnerRule,
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

// DecodeWorkflowOwnerRemoved decodes a log into a WorkflowOwnerRemoved struct.
func (c *Codec) DecodeWorkflowOwnerRemoved(log *evm.Log) (*WorkflowOwnerRemovedDecoded, error) {
	event := new(WorkflowOwnerRemovedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "WorkflowOwnerRemoved", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["WorkflowOwnerRemoved"].Inputs {
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

func (c *Codec) WorkflowOwnerSetLogHash() []byte {
	return c.abi.Events["WorkflowOwnerSet"].ID.Bytes()
}

func (c *Codec) EncodeWorkflowOwnerSetTopics(
	evt abi.Event,
	values []WorkflowOwnerSet,
) ([]*evm.TopicValues, error) {
	var workflowOwnerRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.WorkflowOwner)
		if err != nil {
			return nil, err
		}
		workflowOwnerRule = append(workflowOwnerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		workflowOwnerRule,
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

// DecodeWorkflowOwnerSet decodes a log into a WorkflowOwnerSet struct.
func (c *Codec) DecodeWorkflowOwnerSet(log *evm.Log) (*WorkflowOwnerSetDecoded, error) {
	event := new(WorkflowOwnerSetDecoded)
	if err := c.abi.UnpackIntoInterface(event, "WorkflowOwnerSet", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["WorkflowOwnerSet"].Inputs {
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

func (c ProtocolSmartWallet) AllowedCcipSenders(
	runtime cre.Runtime,
	args AllowedCcipSendersInput,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeAllowedCcipSendersMethodCall(args)
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
		return c.Codec.DecodeAllowedCcipSendersMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) AllowedKeystoneForwarders(
	runtime cre.Runtime,
	args AllowedKeystoneForwardersInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeAllowedKeystoneForwardersMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeAllowedKeystoneForwardersMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) AllowedWorkflowOwners(
	runtime cre.Runtime,
	args AllowedWorkflowOwnersInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeAllowedWorkflowOwnersMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeAllowedWorkflowOwnersMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) GetPoolAddress(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeGetPoolAddressMethodCall()
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
		return c.Codec.DecodeGetPoolAddressMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) GetRouter(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeGetRouterMethodCall()
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
		return c.Codec.DecodeGetRouterMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) Owner(
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

func (c ProtocolSmartWallet) Pool(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodePoolMethodCall()
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
		return c.Codec.DecodePoolMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) SupportsInterface(
	runtime cre.Runtime,
	args SupportsInterfaceInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeSupportsInterfaceMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeSupportsInterfaceMethodOutput(response.Data)
	})

}

func (c ProtocolSmartWallet) WriteReportFromClientAny2EVMMessage(
	runtime cre.Runtime,
	input ClientAny2EVMMessage,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeClientAny2EVMMessageStruct(input)
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

func (c ProtocolSmartWallet) WriteReportFromClientEVMTokenAmount(
	runtime cre.Runtime,
	input ClientEVMTokenAmount,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeClientEVMTokenAmountStruct(input)
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

func (c ProtocolSmartWallet) WriteReportFromRebalanceParams(
	runtime cre.Runtime,
	input RebalanceParams,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeRebalanceParamsStruct(input)
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

func (c ProtocolSmartWallet) WriteReport(
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

// DecodeInsufficientFeeTokenAmountError decodes a InsufficientFeeTokenAmount error from revert data.
func (c *ProtocolSmartWallet) DecodeInsufficientFeeTokenAmountError(data []byte) (*InsufficientFeeTokenAmount, error) {
	args := c.ABI.Errors["InsufficientFeeTokenAmount"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InsufficientFeeTokenAmount{}, nil
}

// Error implements the error interface for InsufficientFeeTokenAmount.
func (e *InsufficientFeeTokenAmount) Error() string {
	return fmt.Sprintf("InsufficientFeeTokenAmount error:")
}

// DecodeInsufficientTokenAmountError decodes a InsufficientTokenAmount error from revert data.
func (c *ProtocolSmartWallet) DecodeInsufficientTokenAmountError(data []byte) (*InsufficientTokenAmount, error) {
	args := c.ABI.Errors["InsufficientTokenAmount"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InsufficientTokenAmount{}, nil
}

// Error implements the error interface for InsufficientTokenAmount.
func (e *InsufficientTokenAmount) Error() string {
	return fmt.Sprintf("InsufficientTokenAmount error:")
}

// DecodeInvalidKeystoneForwarderError decodes a InvalidKeystoneForwarder error from revert data.
func (c *ProtocolSmartWallet) DecodeInvalidKeystoneForwarderError(data []byte) (*InvalidKeystoneForwarder, error) {
	args := c.ABI.Errors["InvalidKeystoneForwarder"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidKeystoneForwarder{}, nil
}

// Error implements the error interface for InvalidKeystoneForwarder.
func (e *InvalidKeystoneForwarder) Error() string {
	return fmt.Sprintf("InvalidKeystoneForwarder error:")
}

// DecodeInvalidRouterError decodes a InvalidRouter error from revert data.
func (c *ProtocolSmartWallet) DecodeInvalidRouterError(data []byte) (*InvalidRouter, error) {
	args := c.ABI.Errors["InvalidRouter"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	router, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for router in InvalidRouter error")
	}

	return &InvalidRouter{
		Router: router,
	}, nil
}

// Error implements the error interface for InvalidRouter.
func (e *InvalidRouter) Error() string {
	return fmt.Sprintf("InvalidRouter error: router=%v;", e.Router)
}

// DecodeInvalidSenderAddressError decodes a InvalidSenderAddress error from revert data.
func (c *ProtocolSmartWallet) DecodeInvalidSenderAddressError(data []byte) (*InvalidSenderAddress, error) {
	args := c.ABI.Errors["InvalidSenderAddress"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidSenderAddress{}, nil
}

// Error implements the error interface for InvalidSenderAddress.
func (e *InvalidSenderAddress) Error() string {
	return fmt.Sprintf("InvalidSenderAddress error:")
}

// DecodeInvalidSourceChainError decodes a InvalidSourceChain error from revert data.
func (c *ProtocolSmartWallet) DecodeInvalidSourceChainError(data []byte) (*InvalidSourceChain, error) {
	args := c.ABI.Errors["InvalidSourceChain"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidSourceChain{}, nil
}

// Error implements the error interface for InvalidSourceChain.
func (e *InvalidSourceChain) Error() string {
	return fmt.Sprintf("InvalidSourceChain error:")
}

// DecodeInvalidWorkflowOwnerError decodes a InvalidWorkflowOwner error from revert data.
func (c *ProtocolSmartWallet) DecodeInvalidWorkflowOwnerError(data []byte) (*InvalidWorkflowOwner, error) {
	args := c.ABI.Errors["InvalidWorkflowOwner"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidWorkflowOwner{}, nil
}

// Error implements the error interface for InvalidWorkflowOwner.
func (e *InvalidWorkflowOwner) Error() string {
	return fmt.Sprintf("InvalidWorkflowOwner error:")
}

// DecodeMismatchedTokenAmountError decodes a MismatchedTokenAmount error from revert data.
func (c *ProtocolSmartWallet) DecodeMismatchedTokenAmountError(data []byte) (*MismatchedTokenAmount, error) {
	args := c.ABI.Errors["MismatchedTokenAmount"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &MismatchedTokenAmount{}, nil
}

// Error implements the error interface for MismatchedTokenAmount.
func (e *MismatchedTokenAmount) Error() string {
	return fmt.Sprintf("MismatchedTokenAmount error:")
}

// DecodeMustBeKeystoneForwarderError decodes a MustBeKeystoneForwarder error from revert data.
func (c *ProtocolSmartWallet) DecodeMustBeKeystoneForwarderError(data []byte) (*MustBeKeystoneForwarder, error) {
	args := c.ABI.Errors["MustBeKeystoneForwarder"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &MustBeKeystoneForwarder{}, nil
}

// Error implements the error interface for MustBeKeystoneForwarder.
func (e *MustBeKeystoneForwarder) Error() string {
	return fmt.Sprintf("MustBeKeystoneForwarder error:")
}

// DecodeNoSenderOnSourceChainError decodes a NoSenderOnSourceChain error from revert data.
func (c *ProtocolSmartWallet) DecodeNoSenderOnSourceChainError(data []byte) (*NoSenderOnSourceChain, error) {
	args := c.ABI.Errors["NoSenderOnSourceChain"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	sourceChainSelector, ok0 := values[0].(uint64)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for sourceChainSelector in NoSenderOnSourceChain error")
	}

	return &NoSenderOnSourceChain{
		SourceChainSelector: sourceChainSelector,
	}, nil
}

// Error implements the error interface for NoSenderOnSourceChain.
func (e *NoSenderOnSourceChain) Error() string {
	return fmt.Sprintf("NoSenderOnSourceChain error: sourceChainSelector=%v;", e.SourceChainSelector)
}

// DecodeUnauthorizedWorkflowOwnerError decodes a UnauthorizedWorkflowOwner error from revert data.
func (c *ProtocolSmartWallet) DecodeUnauthorizedWorkflowOwnerError(data []byte) (*UnauthorizedWorkflowOwner, error) {
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

// DecodeWrongSenderForSourceChainError decodes a WrongSenderForSourceChain error from revert data.
func (c *ProtocolSmartWallet) DecodeWrongSenderForSourceChainError(data []byte) (*WrongSenderForSourceChain, error) {
	args := c.ABI.Errors["WrongSenderForSourceChain"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	sourceChainSelector, ok0 := values[0].(uint64)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for sourceChainSelector in WrongSenderForSourceChain error")
	}

	return &WrongSenderForSourceChain{
		SourceChainSelector: sourceChainSelector,
	}, nil
}

// Error implements the error interface for WrongSenderForSourceChain.
func (e *WrongSenderForSourceChain) Error() string {
	return fmt.Sprintf("WrongSenderForSourceChain error: sourceChainSelector=%v;", e.SourceChainSelector)
}

// DecodeZeroAddressError decodes a ZeroAddress error from revert data.
func (c *ProtocolSmartWallet) DecodeZeroAddressError(data []byte) (*ZeroAddress, error) {
	args := c.ABI.Errors["ZeroAddress"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	index, ok0 := values[0].(*big.Int)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for index in ZeroAddress error")
	}

	return &ZeroAddress{
		Index: index,
	}, nil
}

// Error implements the error interface for ZeroAddress.
func (e *ZeroAddress) Error() string {
	return fmt.Sprintf("ZeroAddress error: index=%v;", e.Index)
}

func (c *ProtocolSmartWallet) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["InsufficientFeeTokenAmount"].ID.Bytes()[:4]):
		return c.DecodeInsufficientFeeTokenAmountError(data)
	case common.Bytes2Hex(c.ABI.Errors["InsufficientTokenAmount"].ID.Bytes()[:4]):
		return c.DecodeInsufficientTokenAmountError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidKeystoneForwarder"].ID.Bytes()[:4]):
		return c.DecodeInvalidKeystoneForwarderError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidRouter"].ID.Bytes()[:4]):
		return c.DecodeInvalidRouterError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidSenderAddress"].ID.Bytes()[:4]):
		return c.DecodeInvalidSenderAddressError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidSourceChain"].ID.Bytes()[:4]):
		return c.DecodeInvalidSourceChainError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidWorkflowOwner"].ID.Bytes()[:4]):
		return c.DecodeInvalidWorkflowOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["MismatchedTokenAmount"].ID.Bytes()[:4]):
		return c.DecodeMismatchedTokenAmountError(data)
	case common.Bytes2Hex(c.ABI.Errors["MustBeKeystoneForwarder"].ID.Bytes()[:4]):
		return c.DecodeMustBeKeystoneForwarderError(data)
	case common.Bytes2Hex(c.ABI.Errors["NoSenderOnSourceChain"].ID.Bytes()[:4]):
		return c.DecodeNoSenderOnSourceChainError(data)
	case common.Bytes2Hex(c.ABI.Errors["UnauthorizedWorkflowOwner"].ID.Bytes()[:4]):
		return c.DecodeUnauthorizedWorkflowOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["WrongSenderForSourceChain"].ID.Bytes()[:4]):
		return c.DecodeWrongSenderForSourceChainError(data)
	case common.Bytes2Hex(c.ABI.Errors["ZeroAddress"].ID.Bytes()[:4]):
		return c.DecodeZeroAddressError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// KeystoneForwarderRemovedTrigger wraps the raw log trigger and provides decoded KeystoneForwarderRemovedDecoded data
type KeystoneForwarderRemovedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into KeystoneForwarderRemoved data
func (t *KeystoneForwarderRemovedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[KeystoneForwarderRemovedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeKeystoneForwarderRemoved(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode KeystoneForwarderRemoved log: %w", err)
	}

	return &bindings.DecodedLog[KeystoneForwarderRemovedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerKeystoneForwarderRemovedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []KeystoneForwarderRemoved) (cre.Trigger[*evm.Log, *bindings.DecodedLog[KeystoneForwarderRemovedDecoded]], error) {
	event := c.ABI.Events["KeystoneForwarderRemoved"]
	topics, err := c.Codec.EncodeKeystoneForwarderRemovedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for KeystoneForwarderRemoved: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &KeystoneForwarderRemovedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsKeystoneForwarderRemoved(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.KeystoneForwarderRemovedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// KeystoneForwarderSetTrigger wraps the raw log trigger and provides decoded KeystoneForwarderSetDecoded data
type KeystoneForwarderSetTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into KeystoneForwarderSet data
func (t *KeystoneForwarderSetTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[KeystoneForwarderSetDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeKeystoneForwarderSet(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode KeystoneForwarderSet log: %w", err)
	}

	return &bindings.DecodedLog[KeystoneForwarderSetDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerKeystoneForwarderSetLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []KeystoneForwarderSet) (cre.Trigger[*evm.Log, *bindings.DecodedLog[KeystoneForwarderSetDecoded]], error) {
	event := c.ABI.Events["KeystoneForwarderSet"]
	topics, err := c.Codec.EncodeKeystoneForwarderSetTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for KeystoneForwarderSet: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &KeystoneForwarderSetTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsKeystoneForwarderSet(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.KeystoneForwarderSetLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// MessageReceivedTrigger wraps the raw log trigger and provides decoded MessageReceivedDecoded data
type MessageReceivedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into MessageReceived data
func (t *MessageReceivedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[MessageReceivedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeMessageReceived(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MessageReceived log: %w", err)
	}

	return &bindings.DecodedLog[MessageReceivedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerMessageReceivedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []MessageReceived) (cre.Trigger[*evm.Log, *bindings.DecodedLog[MessageReceivedDecoded]], error) {
	event := c.ABI.Events["MessageReceived"]
	topics, err := c.Codec.EncodeMessageReceivedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for MessageReceived: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &MessageReceivedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsMessageReceived(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.MessageReceivedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// MessageSentTrigger wraps the raw log trigger and provides decoded MessageSentDecoded data
type MessageSentTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into MessageSent data
func (t *MessageSentTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[MessageSentDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeMessageSent(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode MessageSent log: %w", err)
	}

	return &bindings.DecodedLog[MessageSentDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerMessageSentLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []MessageSent) (cre.Trigger[*evm.Log, *bindings.DecodedLog[MessageSentDecoded]], error) {
	event := c.ABI.Events["MessageSent"]
	topics, err := c.Codec.EncodeMessageSentTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for MessageSent: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &MessageSentTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsMessageSent(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.MessageSentLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// OwnershipTransferRequestedTrigger wraps the raw log trigger and provides decoded OwnershipTransferRequestedDecoded data
type OwnershipTransferRequestedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
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

func (c *ProtocolSmartWallet) LogTriggerOwnershipTransferRequestedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferRequested) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferRequestedDecoded]], error) {
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

func (c *ProtocolSmartWallet) FilterLogsOwnershipTransferRequested(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
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
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
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

func (c *ProtocolSmartWallet) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferred) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
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

func (c *ProtocolSmartWallet) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
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

// ReportReceivedTrigger wraps the raw log trigger and provides decoded ReportReceivedDecoded data
type ReportReceivedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into ReportReceived data
func (t *ReportReceivedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ReportReceivedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeReportReceived(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ReportReceived log: %w", err)
	}

	return &bindings.DecodedLog[ReportReceivedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerReportReceivedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ReportReceived) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ReportReceivedDecoded]], error) {
	event := c.ABI.Events["ReportReceived"]
	topics, err := c.Codec.EncodeReportReceivedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ReportReceived: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ReportReceivedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsReportReceived(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ReportReceivedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// SenderForSourceChainRemovedTrigger wraps the raw log trigger and provides decoded SenderForSourceChainRemovedDecoded data
type SenderForSourceChainRemovedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into SenderForSourceChainRemoved data
func (t *SenderForSourceChainRemovedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[SenderForSourceChainRemovedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeSenderForSourceChainRemoved(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SenderForSourceChainRemoved log: %w", err)
	}

	return &bindings.DecodedLog[SenderForSourceChainRemovedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerSenderForSourceChainRemovedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []SenderForSourceChainRemoved) (cre.Trigger[*evm.Log, *bindings.DecodedLog[SenderForSourceChainRemovedDecoded]], error) {
	event := c.ABI.Events["SenderForSourceChainRemoved"]
	topics, err := c.Codec.EncodeSenderForSourceChainRemovedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for SenderForSourceChainRemoved: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &SenderForSourceChainRemovedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsSenderForSourceChainRemoved(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.SenderForSourceChainRemovedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// SenderForSourceChainSetTrigger wraps the raw log trigger and provides decoded SenderForSourceChainSetDecoded data
type SenderForSourceChainSetTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into SenderForSourceChainSet data
func (t *SenderForSourceChainSetTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[SenderForSourceChainSetDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeSenderForSourceChainSet(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SenderForSourceChainSet log: %w", err)
	}

	return &bindings.DecodedLog[SenderForSourceChainSetDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerSenderForSourceChainSetLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []SenderForSourceChainSet) (cre.Trigger[*evm.Log, *bindings.DecodedLog[SenderForSourceChainSetDecoded]], error) {
	event := c.ABI.Events["SenderForSourceChainSet"]
	topics, err := c.Codec.EncodeSenderForSourceChainSetTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for SenderForSourceChainSet: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &SenderForSourceChainSetTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsSenderForSourceChainSet(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.SenderForSourceChainSetLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// WorkflowOwnerRemovedTrigger wraps the raw log trigger and provides decoded WorkflowOwnerRemovedDecoded data
type WorkflowOwnerRemovedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into WorkflowOwnerRemoved data
func (t *WorkflowOwnerRemovedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[WorkflowOwnerRemovedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeWorkflowOwnerRemoved(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WorkflowOwnerRemoved log: %w", err)
	}

	return &bindings.DecodedLog[WorkflowOwnerRemovedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerWorkflowOwnerRemovedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []WorkflowOwnerRemoved) (cre.Trigger[*evm.Log, *bindings.DecodedLog[WorkflowOwnerRemovedDecoded]], error) {
	event := c.ABI.Events["WorkflowOwnerRemoved"]
	topics, err := c.Codec.EncodeWorkflowOwnerRemovedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for WorkflowOwnerRemoved: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &WorkflowOwnerRemovedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsWorkflowOwnerRemoved(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.WorkflowOwnerRemovedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// WorkflowOwnerSetTrigger wraps the raw log trigger and provides decoded WorkflowOwnerSetDecoded data
type WorkflowOwnerSetTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                      // Embed the raw trigger
	contract                        *ProtocolSmartWallet // Keep reference for decoding
}

// Adapt method that decodes the log into WorkflowOwnerSet data
func (t *WorkflowOwnerSetTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[WorkflowOwnerSetDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeWorkflowOwnerSet(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode WorkflowOwnerSet log: %w", err)
	}

	return &bindings.DecodedLog[WorkflowOwnerSetDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *ProtocolSmartWallet) LogTriggerWorkflowOwnerSetLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []WorkflowOwnerSet) (cre.Trigger[*evm.Log, *bindings.DecodedLog[WorkflowOwnerSetDecoded]], error) {
	event := c.ABI.Events["WorkflowOwnerSet"]
	topics, err := c.Codec.EncodeWorkflowOwnerSetTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for WorkflowOwnerSet: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &WorkflowOwnerSetTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *ProtocolSmartWallet) FilterLogsWorkflowOwnerSet(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.WorkflowOwnerSetLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}
