// Code generated — DO NOT EDIT.

//go:build !wasip1

package price_feed_aggregator

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

// PriceFeedAggregatorMock is a mock implementation of PriceFeedAggregator for testing.
type PriceFeedAggregatorMock struct {
	AccessController        func() (common.Address, error)
	Aggregator              func() (common.Address, error)
	Decimals                func() (uint8, error)
	Description             func() (string, error)
	GetAnswer               func(GetAnswerInput) (*big.Int, error)
	GetRoundData            func(GetRoundDataInput) (GetRoundDataOutput, error)
	GetTimestamp            func(GetTimestampInput) (*big.Int, error)
	LatestAnswer            func() (*big.Int, error)
	LatestRound             func() (*big.Int, error)
	LatestRoundData         func() (LatestRoundDataOutput, error)
	LatestTimestamp         func() (*big.Int, error)
	Owner                   func() (common.Address, error)
	PhaseAggregators        func(PhaseAggregatorsInput) (common.Address, error)
	PhaseId                 func() (uint16, error)
	ProposedAggregator      func() (common.Address, error)
	ProposedGetRoundData    func(ProposedGetRoundDataInput) (ProposedGetRoundDataOutput, error)
	ProposedLatestRoundData func() (ProposedLatestRoundDataOutput, error)
	Version                 func() (*big.Int, error)
}

// NewPriceFeedAggregatorMock creates a new PriceFeedAggregatorMock for testing.
func NewPriceFeedAggregatorMock(address common.Address, clientMock *evmmock.ClientCapability) *PriceFeedAggregatorMock {
	mock := &PriceFeedAggregatorMock{}

	codec, err := NewCodec()
	if err != nil {
		panic("failed to create codec for mock: " + err.Error())
	}

	abi := codec.(*Codec).abi
	_ = abi

	funcMap := map[string]func([]byte) ([]byte, error){
		string(abi.Methods["accessController"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.AccessController == nil {
				return nil, errors.New("accessController method not mocked")
			}
			result, err := mock.AccessController()
			if err != nil {
				return nil, err
			}
			return abi.Methods["accessController"].Outputs.Pack(result)
		},
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
				result.RoundId,
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
				result.RoundId,
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
		string(abi.Methods["phaseAggregators"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PhaseAggregators == nil {
				return nil, errors.New("phaseAggregators method not mocked")
			}
			inputs := abi.Methods["phaseAggregators"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := PhaseAggregatorsInput{
				Arg0: values[0].(uint16),
			}

			result, err := mock.PhaseAggregators(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["phaseAggregators"].Outputs.Pack(result)
		},
		string(abi.Methods["phaseId"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.PhaseId == nil {
				return nil, errors.New("phaseId method not mocked")
			}
			result, err := mock.PhaseId()
			if err != nil {
				return nil, err
			}
			return abi.Methods["phaseId"].Outputs.Pack(result)
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
		string(abi.Methods["proposedGetRoundData"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.ProposedGetRoundData == nil {
				return nil, errors.New("proposedGetRoundData method not mocked")
			}
			inputs := abi.Methods["proposedGetRoundData"].Inputs

			values, err := inputs.Unpack(payload)
			if err != nil {
				return nil, errors.New("Failed to unpack payload")
			}
			if len(values) != 1 {
				return nil, errors.New("expected 1 input value")
			}

			args := ProposedGetRoundDataInput{
				RoundId: values[0].(*big.Int),
			}

			result, err := mock.ProposedGetRoundData(args)
			if err != nil {
				return nil, err
			}
			return abi.Methods["proposedGetRoundData"].Outputs.Pack(
				result.RoundId,
				result.Answer,
				result.StartedAt,
				result.UpdatedAt,
				result.AnsweredInRound,
			)
		},
		string(abi.Methods["proposedLatestRoundData"].ID[:4]): func(payload []byte) ([]byte, error) {
			if mock.ProposedLatestRoundData == nil {
				return nil, errors.New("proposedLatestRoundData method not mocked")
			}
			result, err := mock.ProposedLatestRoundData()
			if err != nil {
				return nil, err
			}
			return abi.Methods["proposedLatestRoundData"].Outputs.Pack(
				result.RoundId,
				result.Answer,
				result.StartedAt,
				result.UpdatedAt,
				result.AnsweredInRound,
			)
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
