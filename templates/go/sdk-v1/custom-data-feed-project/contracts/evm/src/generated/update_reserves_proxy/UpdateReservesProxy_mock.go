// Code generated — DO NOT EDIT.

//go:build !wasip1

package update_reserves_proxy

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

// UpdateReservesProxyMock is a mock implementation of UpdateReservesProxy for testing.
type UpdateReservesProxyMock struct {
	ReserveManager func() (common.Address, error)
}

// NewUpdateReservesProxyMock creates a new UpdateReservesProxyMock for testing.
func NewUpdateReservesProxyMock(address common.Address, clientMock *evmmock.ClientCapability) *UpdateReservesProxyMock {
	mock := &UpdateReservesProxyMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
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
