// Code generated — DO NOT EDIT.

//go:build !wasip1

package reserve_manager

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

// ReserveManagerMock is a mock implementation of ReserveManager for testing.
type ReserveManagerMock struct {
	LastTotalMinted  func() (*big.Int, error)
	LastTotalReserve func() (*big.Int, error)
}

// NewReserveManagerMock creates a new ReserveManagerMock for testing.
func NewReserveManagerMock(address common.Address, clientMock *evmmock.ClientCapability) *ReserveManagerMock {
	mock := &ReserveManagerMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["lastTotalMinted"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LastTotalMinted == nil {
				return nil, errors.New("lastTotalMinted method not mocked")
			}
			result, err := mock.LastTotalMinted()
			if err != nil {
				return nil, err
			}
			return abi.Methods["lastTotalMinted"].Outputs.Pack(result)
		},
		string(abi.Methods["lastTotalReserve"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LastTotalReserve == nil {
				return nil, errors.New("lastTotalReserve method not mocked")
			}
			result, err := mock.LastTotalReserve()
			if err != nil {
				return nil, err
			}
			return abi.Methods["lastTotalReserve"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
