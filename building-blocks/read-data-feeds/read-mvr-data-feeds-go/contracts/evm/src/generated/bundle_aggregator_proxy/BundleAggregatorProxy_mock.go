// Code generated — DO NOT EDIT.

//go:build !wasip1

package bundle_aggregator_proxy

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

// BundleAggregatorProxyMock is a mock implementation of BundleAggregatorProxy for testing.
type BundleAggregatorProxyMock struct {
	Aggregator            func() (common.Address, error)
	BundleDecimals        func() ([]uint8, error)
	Description           func() (string, error)
	LatestBundle          func() ([]byte, error)
	LatestBundleTimestamp func() (*big.Int, error)
	Owner                 func() (common.Address, error)
	ProposedAggregator    func() (common.Address, error)
	TypeAndVersion        func() (string, error)
	Version               func() (*big.Int, error)
}

// NewBundleAggregatorProxyMock creates a new BundleAggregatorProxyMock for testing.
func NewBundleAggregatorProxyMock(address common.Address, clientMock *evmmock.ClientCapability) *BundleAggregatorProxyMock {
	mock := &BundleAggregatorProxyMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["aggregator"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Aggregator == nil {
				return nil, errors.New("aggregator method not mocked")
			}
			result, err := mock.Aggregator()
			if err != nil {
				return nil, err
			}
			return abi.Methods["aggregator"].Outputs.Pack(result)
		},
		string(abi.Methods["bundleDecimals"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.BundleDecimals == nil {
				return nil, errors.New("bundleDecimals method not mocked")
			}
			result, err := mock.BundleDecimals()
			if err != nil {
				return nil, err
			}
			return abi.Methods["bundleDecimals"].Outputs.Pack(result)
		},
		string(abi.Methods["description"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Description == nil {
				return nil, errors.New("description method not mocked")
			}
			result, err := mock.Description()
			if err != nil {
				return nil, err
			}
			return abi.Methods["description"].Outputs.Pack(result)
		},
		string(abi.Methods["latestBundle"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LatestBundle == nil {
				return nil, errors.New("latestBundle method not mocked")
			}
			result, err := mock.LatestBundle()
			if err != nil {
				return nil, err
			}
			return abi.Methods["latestBundle"].Outputs.Pack(result)
		},
		string(abi.Methods["latestBundleTimestamp"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LatestBundleTimestamp == nil {
				return nil, errors.New("latestBundleTimestamp method not mocked")
			}
			result, err := mock.LatestBundleTimestamp()
			if err != nil {
				return nil, err
			}
			return abi.Methods["latestBundleTimestamp"].Outputs.Pack(result)
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
		string(abi.Methods["proposedAggregator"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.ProposedAggregator == nil {
				return nil, errors.New("proposedAggregator method not mocked")
			}
			result, err := mock.ProposedAggregator()
			if err != nil {
				return nil, err
			}
			return abi.Methods["proposedAggregator"].Outputs.Pack(result)
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
		string(abi.Methods["version"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.Version == nil {
				return nil, errors.New("version method not mocked")
			}
			result, err := mock.Version()
			if err != nil {
				return nil, err
			}
			return abi.Methods["version"].Outputs.Pack(result)
		},
	}

	evmmock.AddContractMock(address, clientMock, funcMap, nil)
	return mock
}
