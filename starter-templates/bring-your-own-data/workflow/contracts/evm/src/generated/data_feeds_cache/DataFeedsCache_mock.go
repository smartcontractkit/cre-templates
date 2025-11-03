// Code generated — DO NOT EDIT.

//go:build !wasip1

package data_feeds_cache

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

// DataFeedsCacheMock is a mock implementation of DataFeedsCache for testing.
type DataFeedsCacheMock struct {
	BundleDecimals           func() ([]uint8, error)
	CheckFeedPermission      func(CheckFeedPermissionInput) (bool, error)
	Decimals                 func() (uint8, error)
	Description              func() (string, error)
	GetAnswer                func(GetAnswerInput) (*big.Int, error)
	GetBundleDecimals        func(GetBundleDecimalsInput) ([]uint8, error)
	GetDataIdForProxy        func(GetDataIdForProxyInput) ([16]byte, error)
	GetDescription           func(GetDescriptionInput) (string, error)
	GetFeedMetadata          func(GetFeedMetadataInput) ([]WorkflowMetadata, error)
	GetLatestAnswer          func(GetLatestAnswerInput) (*big.Int, error)
	GetLatestBundle          func(GetLatestBundleInput) ([]byte, error)
	GetLatestBundleTimestamp func(GetLatestBundleTimestampInput) (*big.Int, error)
	GetLatestRoundData       func(GetLatestRoundDataInput) (GetLatestRoundDataOutput, error)
	GetLatestTimestamp       func(GetLatestTimestampInput) (*big.Int, error)
	GetRoundData             func(GetRoundDataInput) (GetRoundDataOutput, error)
	GetTimestamp             func(GetTimestampInput) (*big.Int, error)
	IsFeedAdmin              func(IsFeedAdminInput) (bool, error)
	LatestAnswer             func() (*big.Int, error)
	LatestBundle             func() ([]byte, error)
	LatestBundleTimestamp    func() (*big.Int, error)
	LatestRound              func() (*big.Int, error)
	LatestRoundData          func() (LatestRoundDataOutput, error)
	LatestTimestamp          func() (*big.Int, error)
	Owner                    func() (common.Address, error)
	TypeAndVersion           func() (string, error)
	Version                  func() (*big.Int, error)
}

// NewDataFeedsCacheMock creates a new DataFeedsCacheMock for testing.
func NewDataFeedsCacheMock(address common.Address, clientMock *evmmock.ClientCapability) *DataFeedsCacheMock {
	mock := &DataFeedsCacheMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
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
		string(abi.Methods["checkFeedPermission"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.CheckFeedPermission == nil {
				return nil, errors.New("checkFeedPermission method not mocked")
			}
			inputs := abi.Methods["checkFeedPermission"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 2 {
				return nil, errors.New("expected 2 input values")
			}

			args := CheckFeedPermissionInput{
				DataId:           values[0].([16]byte),
				WorkflowMetadata: values[1].(WorkflowMetadata),
			}

			result, err := mock.CheckFeedPermission(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["checkFeedPermission"].Outputs.Pack(result)
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
		string(abi.Methods["getAnswer"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetAnswer == nil {
				return nil, errors.New("getAnswer method not mocked")
			}
			inputs := abi.Methods["getAnswer"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetAnswerInput{
				RoundId: values[0].(*big.Int),
			}

			result, err := mock.GetAnswer(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getAnswer"].Outputs.Pack(result)
		},
		string(abi.Methods["getBundleDecimals"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetBundleDecimals == nil {
				return nil, errors.New("getBundleDecimals method not mocked")
			}
			inputs := abi.Methods["getBundleDecimals"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetBundleDecimalsInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetBundleDecimals(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getBundleDecimals"].Outputs.Pack(result)
		},
		string(abi.Methods["getDataIdForProxy"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetDataIdForProxy == nil {
				return nil, errors.New("getDataIdForProxy method not mocked")
			}
			inputs := abi.Methods["getDataIdForProxy"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetDataIdForProxyInput{
				Proxy: values[0].(common.Address),
			}

			result, err := mock.GetDataIdForProxy(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getDataIdForProxy"].Outputs.Pack(result)
		},
		string(abi.Methods["getDescription"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetDescription == nil {
				return nil, errors.New("getDescription method not mocked")
			}
			inputs := abi.Methods["getDescription"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetDescriptionInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetDescription(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getDescription"].Outputs.Pack(result)
		},
		string(abi.Methods["getFeedMetadata"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetFeedMetadata == nil {
				return nil, errors.New("getFeedMetadata method not mocked")
			}
			inputs := abi.Methods["getFeedMetadata"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 3 {
				return nil, errors.New("expected 3 input values")
			}

			args := GetFeedMetadataInput{
				DataId:     values[0].([16]byte),
				StartIndex: values[1].(*big.Int),
				MaxCount:   values[2].(*big.Int),
			}

			result, err := mock.GetFeedMetadata(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getFeedMetadata"].Outputs.Pack(result)
		},
		string(abi.Methods["getLatestAnswer"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetLatestAnswer == nil {
				return nil, errors.New("getLatestAnswer method not mocked")
			}
			inputs := abi.Methods["getLatestAnswer"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetLatestAnswerInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetLatestAnswer(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getLatestAnswer"].Outputs.Pack(result)
		},
		string(abi.Methods["getLatestBundle"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetLatestBundle == nil {
				return nil, errors.New("getLatestBundle method not mocked")
			}
			inputs := abi.Methods["getLatestBundle"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetLatestBundleInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetLatestBundle(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getLatestBundle"].Outputs.Pack(result)
		},
		string(abi.Methods["getLatestBundleTimestamp"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetLatestBundleTimestamp == nil {
				return nil, errors.New("getLatestBundleTimestamp method not mocked")
			}
			inputs := abi.Methods["getLatestBundleTimestamp"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetLatestBundleTimestampInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetLatestBundleTimestamp(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getLatestBundleTimestamp"].Outputs.Pack(result)
		},
		string(abi.Methods["getLatestRoundData"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetLatestRoundData == nil {
				return nil, errors.New("getLatestRoundData method not mocked")
			}
			inputs := abi.Methods["getLatestRoundData"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetLatestRoundDataInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetLatestRoundData(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getLatestRoundData"].Outputs.Pack(
				result.Id,
				result.Answer,
				result.StartedAt,
				result.UpdatedAt,
				result.AnsweredInRound,
			)
		},
		string(abi.Methods["getLatestTimestamp"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetLatestTimestamp == nil {
				return nil, errors.New("getLatestTimestamp method not mocked")
			}
			inputs := abi.Methods["getLatestTimestamp"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetLatestTimestampInput{
				DataId: values[0].([16]byte),
			}

			result, err := mock.GetLatestTimestamp(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getLatestTimestamp"].Outputs.Pack(result)
		},
		string(abi.Methods["getRoundData"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetRoundData == nil {
				return nil, errors.New("getRoundData method not mocked")
			}
			inputs := abi.Methods["getRoundData"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetRoundDataInput{
				RoundId: values[0].(*big.Int),
			}

			result, err := mock.GetRoundData(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getRoundData"].Outputs.Pack(
				result.Id,
				result.Answer,
				result.StartedAt,
				result.UpdatedAt,
				result.AnsweredInRound,
			)
		},
		string(abi.Methods["getTimestamp"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.GetTimestamp == nil {
				return nil, errors.New("getTimestamp method not mocked")
			}
			inputs := abi.Methods["getTimestamp"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := GetTimestampInput{
				RoundId: values[0].(*big.Int),
			}

			result, err := mock.GetTimestamp(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["getTimestamp"].Outputs.Pack(result)
		},
		string(abi.Methods["isFeedAdmin"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.IsFeedAdmin == nil {
				return nil, errors.New("isFeedAdmin method not mocked")
			}
			inputs := abi.Methods["isFeedAdmin"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := IsFeedAdminInput{
				FeedAdmin: values[0].(common.Address),
			}

			result, err := mock.IsFeedAdmin(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["isFeedAdmin"].Outputs.Pack(result)
		},
		string(abi.Methods["latestAnswer"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LatestAnswer == nil {
				return nil, errors.New("latestAnswer method not mocked")
			}
			result, err := mock.LatestAnswer()
			if err != nil {
				return nil, err
			}
			return abi.Methods["latestAnswer"].Outputs.Pack(result)
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
		string(abi.Methods["latestRound"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LatestRound == nil {
				return nil, errors.New("latestRound method not mocked")
			}
			result, err := mock.LatestRound()
			if err != nil {
				return nil, err
			}
			return abi.Methods["latestRound"].Outputs.Pack(result)
		},
		string(abi.Methods["latestRoundData"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LatestRoundData == nil {
				return nil, errors.New("latestRoundData method not mocked")
			}
			result, err := mock.LatestRoundData()
			if err != nil {
				return nil, err
			}
			return abi.Methods["latestRoundData"].Outputs.Pack(
				result.Id,
				result.Answer,
				result.StartedAt,
				result.UpdatedAt,
				result.AnsweredInRound,
			)
		},
		string(abi.Methods["latestTimestamp"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.LatestTimestamp == nil {
				return nil, errors.New("latestTimestamp method not mocked")
			}
			result, err := mock.LatestTimestamp()
			if err != nil {
				return nil, err
			}
			return abi.Methods["latestTimestamp"].Outputs.Pack(result)
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
