// Code generated — DO NOT EDIT.

//go:build !wasip1

package randomness_consumer

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

// RandomnessConsumerMock is a mock implementation of RandomnessConsumer for testing.
type RandomnessConsumerMock struct {
	GetExpectedAuthor       func() (common.Address, error)
	GetExpectedWorkflowId   func() ([32]byte, error)
	GetExpectedWorkflowName func() ([10]byte, error)
	GetForwarderAddress     func() (common.Address, error)
	NextRequestId           func() (*big.Int, error)
	Owner                   func() (common.Address, error)
	PendingRequests         func(PendingRequestsInput) (bool, error)
	RandomWords             func(RandomWordsInput) (*big.Int, error)
	SupportsInterface       func(SupportsInterfaceInput) (bool, error)
}

// NewRandomnessConsumerMock creates a new RandomnessConsumerMock for testing.
func NewRandomnessConsumerMock(address common.Address, clientMock *evmmock.ClientCapability) *RandomnessConsumerMock {
	mock := &RandomnessConsumerMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["getExpectedAuthor"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetExpectedAuthor == nil {
				return nil, errors.New("getExpectedAuthor method not mocked")
			}
			result, err := mock.GetExpectedAuthor()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getExpectedAuthor"].Outputs.Pack(result)
		},
		string(abi.Methods["getExpectedWorkflowId"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetExpectedWorkflowId == nil {
				return nil, errors.New("getExpectedWorkflowId method not mocked")
			}
			result, err := mock.GetExpectedWorkflowId()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getExpectedWorkflowId"].Outputs.Pack(result)
		},
		string(abi.Methods["getExpectedWorkflowName"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetExpectedWorkflowName == nil {
				return nil, errors.New("getExpectedWorkflowName method not mocked")
			}
			result, err := mock.GetExpectedWorkflowName()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getExpectedWorkflowName"].Outputs.Pack(result)
		},
		string(abi.Methods["getForwarderAddress"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetForwarderAddress == nil {
				return nil, errors.New("getForwarderAddress method not mocked")
			}
			result, err := mock.GetForwarderAddress()
			if err != nil {
				return nil, err
			}
			return abi.Methods["getForwarderAddress"].Outputs.Pack(result)
		},
		string(abi.Methods["nextRequestId"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.NextRequestId == nil {
				return nil, errors.New("nextRequestId method not mocked")
			}
			result, err := mock.NextRequestId()
			if err != nil {
				return nil, err
			}
			return abi.Methods["nextRequestId"].Outputs.Pack(result)
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
		string(abi.Methods["pendingRequests"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PendingRequests == nil {
				return nil, errors.New("pendingRequests method not mocked")
			}
			inputs := abi.Methods["pendingRequests"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := PendingRequestsInput{
				Arg0: values[0].(*big.Int),
			}

			result, err := mock.PendingRequests(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["pendingRequests"].Outputs.Pack(result)
		},
		string(abi.Methods["randomWords"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.RandomWords == nil {
				return nil, errors.New("randomWords method not mocked")
			}
			inputs := abi.Methods["randomWords"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := RandomWordsInput{
				Arg0: values[0].(*big.Int),
			}

			result, err := mock.RandomWords(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["randomWords"].Outputs.Pack(result)
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
