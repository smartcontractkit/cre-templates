// Code generated — DO NOT EDIT.

//go:build !wasip1

package ireserve_manager

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

// IReserveManagerMock is a mock implementation of IReserveManager for testing.
type IReserveManagerMock struct {
}

// NewIReserveManagerMock creates a new IReserveManagerMock for testing.
func NewIReserveManagerMock(address common.Address, clientMock *evmmock.ClientCapability) *IReserveManagerMock {
	mock := &IReserveManagerMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
