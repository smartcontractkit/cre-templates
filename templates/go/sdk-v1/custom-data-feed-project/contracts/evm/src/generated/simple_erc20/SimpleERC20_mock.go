// Code generated — DO NOT EDIT.

//go:build !wasip1

package simple_erc20

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

// SimpleERC20Mock is a mock implementation of SimpleERC20 for testing.
type SimpleERC20Mock struct {
	Allowance   func(AllowanceInput) (*big.Int, error)
	BalanceOf   func(BalanceOfInput) (*big.Int, error)
	Decimals    func() (uint8, error)
	Name        func() (string, error)
	Symbol      func() (string, error)
	TotalSupply func() (*big.Int, error)
}

// NewSimpleERC20Mock creates a new SimpleERC20Mock for testing.
func NewSimpleERC20Mock(address common.Address, clientMock *evmmock.ClientCapability) *SimpleERC20Mock {
	mock := &SimpleERC20Mock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["allowance"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Allowance == nil {
				return nil, errors.New("allowance method not mocked")
			}
			inputs := abi.Methods["allowance"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := AllowanceInput{
				Owner:   values[0].(common.Address),
				Spender: values[1].(common.Address),
			}

			result, err := mock.Allowance(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["allowance"].Outputs.Pack(result)
		},
		string(abi.Methods["balanceOf"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.BalanceOf == nil {
				return nil, errors.New("balanceOf method not mocked")
			}
			inputs := abi.Methods["balanceOf"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := BalanceOfInput{
				Account: values[0].(common.Address),
			}

			result, err := mock.BalanceOf(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["balanceOf"].Outputs.Pack(result)
		},
		string(abi.Methods["decimals"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Decimals == nil {
				return nil, errors.New("decimals method not mocked")
			}
			result, err := mock.Decimals()
			if err != nil {
				return nil, err
			}
			return abi.Methods["decimals"].Outputs.Pack(result)
		},
		string(abi.Methods["name"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Name == nil {
				return nil, errors.New("name method not mocked")
			}
			result, err := mock.Name()
			if err != nil {
				return nil, err
			}
			return abi.Methods["name"].Outputs.Pack(result)
		},
		string(abi.Methods["symbol"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Symbol == nil {
				return nil, errors.New("symbol method not mocked")
			}
			result, err := mock.Symbol()
			if err != nil {
				return nil, err
			}
			return abi.Methods["symbol"].Outputs.Pack(result)
		},
		string(abi.Methods["totalSupply"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.TotalSupply == nil {
				return nil, errors.New("totalSupply method not mocked")
			}
			result, err := mock.TotalSupply()
			if err != nil {
				return nil, err
			}
			return abi.Methods["totalSupply"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
