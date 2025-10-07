// Code generated — DO NOT EDIT.

//go:build !wasip1

package ireceiver_template

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

// IReceiverTemplateMock is a mock implementation of IReceiverTemplate for testing.
type IReceiverTemplateMock struct {
	EXPECTEDAUTHOR       func() (common.Address, error)
	EXPECTEDWORKFLOWNAME func() ([10]byte, error)
}

// NewIReceiverTemplateMock creates a new IReceiverTemplateMock for testing.
func NewIReceiverTemplateMock(address common.Address, clientMock *evmmock.ClientCapability) *IReceiverTemplateMock {
	mock := &IReceiverTemplateMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["EXPECTED_AUTHOR"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.EXPECTEDAUTHOR == nil {
				return nil, errors.New("EXPECTED_AUTHOR method not mocked")
			}
			result, err := mock.EXPECTEDAUTHOR()
			if err != nil {
				return nil, err
			}
			return abi.Methods["EXPECTED_AUTHOR"].Outputs.Pack(result)
		},
		string(abi.Methods["EXPECTED_WORKFLOW_NAME"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.EXPECTEDWORKFLOWNAME == nil {
				return nil, errors.New("EXPECTED_WORKFLOW_NAME method not mocked")
			}
			result, err := mock.EXPECTEDWORKFLOWNAME()
			if err != nil {
				return nil, err
			}
			return abi.Methods["EXPECTED_WORKFLOW_NAME"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
