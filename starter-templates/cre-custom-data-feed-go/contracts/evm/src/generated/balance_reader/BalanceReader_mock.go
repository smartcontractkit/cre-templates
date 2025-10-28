// Code generated — DO NOT EDIT.

//go:build !wasip1

package balance_reader

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

// BalanceReaderMock is a mock implementation of BalanceReader for testing.
type BalanceReaderMock struct {
	GetNativeBalances func(GetNativeBalancesInput) ([]*big.Int, error)
	TypeAndVersion    func() (string, error)
}

// NewBalanceReaderMock creates a new BalanceReaderMock for testing.
func NewBalanceReaderMock(address common.Address, clientMock *evmmock.ClientCapability) *BalanceReaderMock {
	mock := &BalanceReaderMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["getNativeBalances"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetNativeBalances == nil {
				return nil, errors.New("getNativeBalances method not mocked")
			}
			inputs := abi.Methods["getNativeBalances"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetNativeBalancesInput{
				Addresses: values[0].([]common.Address),
			}

			result, err := mock.GetNativeBalances(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getNativeBalances"].Outputs.Pack(result)
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
