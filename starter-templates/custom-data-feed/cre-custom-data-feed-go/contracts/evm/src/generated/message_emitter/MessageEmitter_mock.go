// Code generated — DO NOT EDIT.

//go:build !wasip1

package message_emitter

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

// MessageEmitterMock is a mock implementation of MessageEmitter for testing.
type MessageEmitterMock struct {
	GetLastMessage func(GetLastMessageInput) (string, error)
	GetMessage     func(GetMessageInput) (string, error)
	TypeAndVersion func() (string, error)
}

// NewMessageEmitterMock creates a new MessageEmitterMock for testing.
func NewMessageEmitterMock(address common.Address, clientMock *evmmock.ClientCapability) *MessageEmitterMock {
	mock := &MessageEmitterMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["getLastMessage"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetLastMessage == nil {
				return nil, errors.New("getLastMessage method not mocked")
			}
			inputs := abi.Methods["getLastMessage"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetLastMessageInput{
				Emitter: values[0].(common.Address),
			}

			result, err := mock.GetLastMessage(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getLastMessage"].Outputs.Pack(result)
		},
		string(abi.Methods["getMessage"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetMessage == nil {
				return nil, errors.New("getMessage method not mocked")
			}
			inputs := abi.Methods["getMessage"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := GetMessageInput{
				Emitter:   values[0].(common.Address),
				Timestamp: values[1].(*big.Int),
			}

			result, err := mock.GetMessage(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getMessage"].Outputs.Pack(result)
		},
		string(abi.Methods["typeAndVersion"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TypeAndVersion == nil {
				return nil, errors.New("typeAndVersion method not mocked")
			}
			result, err := mock.TypeAndVersion()
			if err != nil {
				return nil, err
			}
			return abi.Methods["typeAndVersion"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
