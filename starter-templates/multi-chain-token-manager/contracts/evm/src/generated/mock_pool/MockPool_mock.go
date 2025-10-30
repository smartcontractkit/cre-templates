// Code generated — DO NOT EDIT.

//go:build !wasip1

package mock_pool

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

// MockPoolMock is a mock implementation of MockPool for testing.
type MockPoolMock struct {
	BalanceOf            func(BalanceOfInput) (*big.Int, error)
	CurrentLiquidityRate func(CurrentLiquidityRateInput) (*big.Int, error)
	GetReserveData       func(GetReserveDataInput) (DataTypesReserveDataLegacy, error)
	Owner                func() (common.Address, error)
}

// NewMockPoolMock creates a new MockPoolMock for testing.
func NewMockPoolMock(address common.Address, clientMock *evmmock.ClientCapability) *MockPoolMock {
	mock := &MockPoolMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["balanceOf"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.BalanceOf == nil {
				return nil, errors.New("balanceOf method not mocked")
			}
			inputs := abi.Methods["balanceOf"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := BalanceOfInput{
				Arg0: values[0].(common.Address),
				Arg1: values[1].(common.Address),
			}

			result, err := mock.BalanceOf(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["balanceOf"].Outputs.Pack(result)
		},
		string(abi.Methods["currentLiquidityRate"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.CurrentLiquidityRate == nil {
				return nil, errors.New("currentLiquidityRate method not mocked")
			}
			inputs := abi.Methods["currentLiquidityRate"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := CurrentLiquidityRateInput{
				Arg0: values[0].(common.Address),
			}

			result, err := mock.CurrentLiquidityRate(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["currentLiquidityRate"].Outputs.Pack(result)
		},
		string(abi.Methods["getReserveData"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetReserveData == nil {
				return nil, errors.New("getReserveData method not mocked")
			}
			inputs := abi.Methods["getReserveData"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetReserveDataInput{
				Asset: values[0].(common.Address),
			}

			result, err := mock.GetReserveData(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getReserveData"].Outputs.Pack(result)
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
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
