// Code generated — DO NOT EDIT.

//go:build !wasip1

package protocol_smart_wallet

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	evmmock "github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/mock"
)

var (
	_ = errors.New
	_ = fmt.Errorf
	_ = big.NewInt
	_ = common.Big1
)

// ProtocolSmartWalletMock is a mock implementation of ProtocolSmartWallet for testing.
type ProtocolSmartWalletMock struct {
	AllowedCcipSenders        func(AllowedCcipSendersInput) (common.Address, error)
	AllowedKeystoneForwarders func(AllowedKeystoneForwardersInput) (bool, error)
	AllowedWorkflowOwners     func(AllowedWorkflowOwnersInput) (bool, error)
	GetPoolAddress            func() (common.Address, error)
	GetRouter                 func() (common.Address, error)
	Owner                     func() (common.Address, error)
	Pool                      func() (common.Address, error)
	SupportsInterface         func(SupportsInterfaceInput) (bool, error)
}

// NewProtocolSmartWalletMock creates a new ProtocolSmartWalletMock for testing.
func NewProtocolSmartWalletMock(address common.Address, clientMock *evmmock.ClientCapability) *ProtocolSmartWalletMock {
	mock := &ProtocolSmartWalletMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["allowedCcipSenders"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.AllowedCcipSenders == nil {
				return nil, errors.New("allowedCcipSenders method not mocked")
			}
			inputs := abi.Methods["allowedCcipSenders"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := AllowedCcipSendersInput{
				ChainSelector: values[0].(uint64),
			}

			result, err := mock.AllowedCcipSenders(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["allowedCcipSenders"].Outputs.Pack(result)
		},
		string(abi.Methods["allowedKeystoneForwarders"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.AllowedKeystoneForwarders == nil {
				return nil, errors.New("allowedKeystoneForwarders method not mocked")
			}
			inputs := abi.Methods["allowedKeystoneForwarders"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := AllowedKeystoneForwardersInput{
				KeystoneForwarder: values[0].(common.Address),
			}

			result, err := mock.AllowedKeystoneForwarders(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["allowedKeystoneForwarders"].Outputs.Pack(result)
		},
		string(abi.Methods["allowedWorkflowOwners"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.AllowedWorkflowOwners == nil {
				return nil, errors.New("allowedWorkflowOwners method not mocked")
			}
			inputs := abi.Methods["allowedWorkflowOwners"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := AllowedWorkflowOwnersInput{
				WorkflowOwner: values[0].(common.Address),
			}

			result, err := mock.AllowedWorkflowOwners(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["allowedWorkflowOwners"].Outputs.Pack(result)
		},
		string(abi.Methods["getPoolAddress"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetPoolAddress == nil {
				return nil, errors.New("getPoolAddress method not mocked")
			}
			result, err := mock.GetPoolAddress()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getPoolAddress"].Outputs.Pack(result)
		},
		string(abi.Methods["getRouter"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetRouter == nil {
				return nil, errors.New("getRouter method not mocked")
			}
			result, err := mock.GetRouter()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getRouter"].Outputs.Pack(result)
		},
		string(abi.Methods["owner"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Owner == nil {
				return nil, errors.New("owner method not mocked")
			}
			result, err := mock.Owner()
			if err != nil {
				return nil, err
			}
			return abi.Methods["owner"].Outputs.Pack(result)
		},
		string(abi.Methods["pool"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Pool == nil {
				return nil, errors.New("pool method not mocked")
			}
			result, err := mock.Pool()
			if err != nil {
				return nil, err
			}
			return abi.Methods["pool"].Outputs.Pack(result)
		},
		string(abi.Methods["supportsInterface"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.SupportsInterface == nil {
				return nil, errors.New("supportsInterface method not mocked")
			}
			inputs := abi.Methods["supportsInterface"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := SupportsInterfaceInput{
				InterfaceId: values[0].([4]byte),
			}

			result, err := mock.SupportsInterface(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["supportsInterface"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
