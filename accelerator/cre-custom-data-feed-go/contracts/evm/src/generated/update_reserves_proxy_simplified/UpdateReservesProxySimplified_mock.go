// Code generated — DO NOT EDIT.

//go:build !wasip1

package update_reserves_proxy_simplified

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

// UpdateReservesProxySimplifiedMock is a mock implementation of UpdateReservesProxySimplified for testing.
type UpdateReservesProxySimplifiedMock struct {
	EXPECTEDAUTHOR       func() (common.Address, error)
	EXPECTEDWORKFLOWNAME func() ([10]byte, error)
	ReserveManager       func() (common.Address, error)
}

// NewUpdateReservesProxySimplifiedMock creates a new UpdateReservesProxySimplifiedMock for testing.
func NewUpdateReservesProxySimplifiedMock(address common.Address, clientMock *evmmock.ClientCapability) *UpdateReservesProxySimplifiedMock {
	mock := &UpdateReservesProxySimplifiedMock{}

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
		string(abi.Methods["reserveManager"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.ReserveManager == nil {
				return nil, errors.New("reserveManager method not mocked")
			}
			result, err := mock.ReserveManager()
			if err != nil {
				return nil, err
			}
			return abi.Methods["reserveManager"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
