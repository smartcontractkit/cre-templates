// Code generated — DO NOT EDIT.

//go:build !wasip1

package itype_and_version

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

// ITypeAndVersionMock is a mock implementation of ITypeAndVersion for testing.
type ITypeAndVersionMock struct {
}

// NewITypeAndVersionMock creates a new ITypeAndVersionMock for testing.
func NewITypeAndVersionMock(address common.Address, clientMock *evmmock.ClientCapability) *ITypeAndVersionMock {
	mock := &ITypeAndVersionMock{}

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
