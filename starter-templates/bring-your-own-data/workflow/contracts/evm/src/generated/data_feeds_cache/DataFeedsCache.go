// Code generated — DO NOT EDIT.

package data_feeds_cache

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb2 "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

var (
	_ = bytes.Equal
	_ = errors.New
	_ = fmt.Sprintf
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = emptypb.Empty{}
	_ = pb.NewBigIntFromInt
	_ = pb2.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX
	_ = bindings.FilterOptions{}
	_ = evm.FilterLogTriggerRequest{}
	_ = cre.ResponseBufferTooSmall
	_ = rpc.API{}
	_ = json.Unmarshal
)

var DataFeedsCacheMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"AddressInsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ArrayLengthMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyConfig\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"ErrorSendingNative\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedInnerCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"FeedNotConfigured\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredBalance\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"InvalidAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidDataId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"InvalidWorkflowName\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"NoMappingForSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"}],\"name\":\"UnauthorizedCaller\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"int256\",\"name\":\"current\",\"type\":\"int256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"}],\"name\":\"AnswerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint8[]\",\"name\":\"decimals\",\"type\":\"uint8[]\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedWorkflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\"}],\"indexed\":false,\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"name\":\"workflowMetadata\",\"type\":\"tuple[]\"}],\"name\":\"BundleFeedConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"bundle\",\"type\":\"bytes\"}],\"name\":\"BundleReportUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedWorkflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\"}],\"indexed\":false,\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"name\":\"workflowMetadata\",\"type\":\"tuple[]\"}],\"name\":\"DecimalFeedConfigSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint224\",\"name\":\"answer\",\"type\":\"uint224\"}],\"name\":\"DecimalReportUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"feedAdmin\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bool\",\"name\":\"isAdmin\",\"type\":\"bool\"}],\"name\":\"FeedAdminSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"FeedConfigRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"workflowOwner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes10\",\"name\":\"workflowName\",\"type\":\"bytes10\"}],\"name\":\"InvalidUpdatePermission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"startedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"}],\"name\":\"NewRound\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"ProxyDataIdRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"ProxyDataIdUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reportTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"latestTimestamp\",\"type\":\"uint256\"}],\"name\":\"StaleBundleReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reportTimestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"latestTimestamp\",\"type\":\"uint256\"}],\"name\":\"StaleDecimalReport\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokenRecovered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"acceptOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bundleDecimals\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"bundleFeedDecimals\",\"type\":\"uint8[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedWorkflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\"}],\"internalType\":\"structDataFeedsCache.WorkflowMetadata\",\"name\":\"workflowMetadata\",\"type\":\"tuple\"}],\"name\":\"checkFeedPermission\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"hasPermission\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"feedDecimals\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"description\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"feedDescription\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getBundleDecimals\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"bundleFeedDecimals\",\"type\":\"uint8[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"name\":\"getDataIdForProxy\",\"outputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"feedDecimals\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getDescription\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"feedDescription\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"startIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxCount\",\"type\":\"uint256\"}],\"name\":\"getFeedMetadata\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedWorkflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\"}],\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"name\":\"workflowMetadata\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getLatestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getLatestBundle\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"bundle\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getLatestBundleTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getLatestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"id\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"dataId\",\"type\":\"bytes16\"}],\"name\":\"getLatestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint80\",\"name\":\"roundId\",\"type\":\"uint80\"}],\"name\":\"getRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"id\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"roundId\",\"type\":\"uint256\"}],\"name\":\"getTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feedAdmin\",\"type\":\"address\"}],\"name\":\"isFeedAdmin\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestAnswer\",\"outputs\":[{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestBundle\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"bundle\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestBundleTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRound\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestRoundData\",\"outputs\":[{\"internalType\":\"uint80\",\"name\":\"id\",\"type\":\"uint80\"},{\"internalType\":\"int256\",\"name\":\"answer\",\"type\":\"int256\"},{\"internalType\":\"uint256\",\"name\":\"startedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"updatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint80\",\"name\":\"answeredInRound\",\"type\":\"uint80\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"latestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"report\",\"type\":\"bytes\"}],\"name\":\"onReport\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"recoverTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"proxies\",\"type\":\"address[]\"}],\"name\":\"removeDataIdMappingsForProxies\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"dataIds\",\"type\":\"bytes16[]\"}],\"name\":\"removeFeedConfigs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"dataIds\",\"type\":\"bytes16[]\"},{\"internalType\":\"string[]\",\"name\":\"descriptions\",\"type\":\"string[]\"},{\"internalType\":\"uint8[][]\",\"name\":\"decimalsMatrix\",\"type\":\"uint8[][]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedWorkflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\"}],\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"name\":\"workflowMetadata\",\"type\":\"tuple[]\"}],\"name\":\"setBundleFeedConfigs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes16[]\",\"name\":\"dataIds\",\"type\":\"bytes16[]\"},{\"internalType\":\"string[]\",\"name\":\"descriptions\",\"type\":\"string[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"allowedSender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"allowedWorkflowOwner\",\"type\":\"address\"},{\"internalType\":\"bytes10\",\"name\":\"allowedWorkflowName\",\"type\":\"bytes10\"}],\"internalType\":\"structDataFeedsCache.WorkflowMetadata[]\",\"name\":\"workflowMetadata\",\"type\":\"tuple[]\"}],\"name\":\"setDecimalFeedConfigs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"feedAdmin\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"isAdmin\",\"type\":\"bool\"}],\"name\":\"setFeedAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"typeAndVersion\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"proxies\",\"type\":\"address[]\"},{\"internalType\":\"bytes16[]\",\"name\":\"dataIds\",\"type\":\"bytes16[]\"}],\"name\":\"updateDataIdMappingsForProxies\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Structs
type WorkflowMetadata struct {
	AllowedSender        common.Address
	AllowedWorkflowOwner common.Address
	AllowedWorkflowName  [10]byte
}

// Contract Method Inputs
type CheckFeedPermissionInput struct {
	DataId           [16]byte
	WorkflowMetadata WorkflowMetadata
}

type GetAnswerInput struct {
	RoundId *big.Int
}

type GetBundleDecimalsInput struct {
	DataId [16]byte
}

type GetDataIdForProxyInput struct {
	Proxy common.Address
}

type GetDecimalsInput struct {
	DataId [16]byte
}

type GetDescriptionInput struct {
	DataId [16]byte
}

type GetFeedMetadataInput struct {
	DataId     [16]byte
	StartIndex *big.Int
	MaxCount   *big.Int
}

type GetLatestAnswerInput struct {
	DataId [16]byte
}

type GetLatestBundleInput struct {
	DataId [16]byte
}

type GetLatestBundleTimestampInput struct {
	DataId [16]byte
}

type GetLatestRoundDataInput struct {
	DataId [16]byte
}

type GetLatestTimestampInput struct {
	DataId [16]byte
}

type GetRoundDataInput struct {
	RoundId *big.Int
}

type GetTimestampInput struct {
	RoundId *big.Int
}

type IsFeedAdminInput struct {
	FeedAdmin common.Address
}

type OnReportInput struct {
	Metadata []byte
	Report   []byte
}

type RecoverTokensInput struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
}

type RemoveDataIdMappingsForProxiesInput struct {
	Proxies []common.Address
}

type RemoveFeedConfigsInput struct {
	DataIds [][16]byte
}

type SetBundleFeedConfigsInput struct {
	DataIds          [][16]byte
	Descriptions     []string
	DecimalsMatrix   [][]uint8
	WorkflowMetadata []WorkflowMetadata
}

type SetDecimalFeedConfigsInput struct {
	DataIds          [][16]byte
	Descriptions     []string
	WorkflowMetadata []WorkflowMetadata
}

type SetFeedAdminInput struct {
	FeedAdmin common.Address
	IsAdmin   bool
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
}

type TransferOwnershipInput struct {
	To common.Address
}

type UpdateDataIdMappingsForProxiesInput struct {
	Proxies []common.Address
	DataIds [][16]byte
}

// Contract Method Outputs
type GetLatestRoundDataOutput struct {
	Id              *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

type GetRoundDataOutput struct {
	Id              *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

type LatestRoundDataOutput struct {
	Id              *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

// Errors
type AddressEmptyCode struct {
	Target common.Address
}

type AddressInsufficientBalance struct {
	Account common.Address
}

type ArrayLengthMismatch struct {
}

type EmptyConfig struct {
}

type ErrorSendingNative struct {
	To     common.Address
	Amount *big.Int
	Data   []byte
}

type FailedInnerCall struct {
}

type FeedNotConfigured struct {
	DataId [16]byte
}

type InsufficientBalance struct {
	Balance         *big.Int
	RequiredBalance *big.Int
}

type InvalidAddress struct {
	Addr common.Address
}

type InvalidDataId struct {
}

type InvalidWorkflowName struct {
	WorkflowName [10]byte
}

type NoMappingForSender struct {
	Proxy common.Address
}

type SafeERC20FailedOperation struct {
	Token common.Address
}

type UnauthorizedCaller struct {
	Caller common.Address
}

// Events
// The <Event> struct should be used as a filter (for log triggers).
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type AnswerUpdated struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
}

type AnswerUpdatedDecoded struct {
	Current   *big.Int
	RoundId   *big.Int
	UpdatedAt *big.Int
}

type BundleFeedConfigSet struct {
	DataId           [16]byte
	Decimals         []uint8
	Description      string
	WorkflowMetadata []WorkflowMetadata
}

type BundleFeedConfigSetDecoded struct {
	DataId           [16]byte
	Decimals         []uint8
	Description      string
	WorkflowMetadata []WorkflowMetadata
}

type BundleReportUpdated struct {
	DataId    [16]byte
	Timestamp *big.Int
	Bundle    []byte
}

type BundleReportUpdatedDecoded struct {
	DataId    [16]byte
	Timestamp *big.Int
	Bundle    []byte
}

type DecimalFeedConfigSet struct {
	DataId           [16]byte
	Decimals         uint8
	Description      string
	WorkflowMetadata []WorkflowMetadata
}

type DecimalFeedConfigSetDecoded struct {
	DataId           [16]byte
	Decimals         uint8
	Description      string
	WorkflowMetadata []WorkflowMetadata
}

type DecimalReportUpdated struct {
	DataId    [16]byte
	RoundId   *big.Int
	Timestamp *big.Int
	Answer    *big.Int
}

type DecimalReportUpdatedDecoded struct {
	DataId    [16]byte
	RoundId   *big.Int
	Timestamp *big.Int
	Answer    *big.Int
}

type FeedAdminSet struct {
	FeedAdmin common.Address
	IsAdmin   bool
}

type FeedAdminSetDecoded struct {
	FeedAdmin common.Address
	IsAdmin   bool
}

type FeedConfigRemoved struct {
	DataId [16]byte
}

type FeedConfigRemovedDecoded struct {
	DataId [16]byte
}

type InvalidUpdatePermission struct {
	DataId        [16]byte
	Sender        common.Address
	WorkflowOwner common.Address
	WorkflowName  [10]byte
}

type InvalidUpdatePermissionDecoded struct {
	DataId        [16]byte
	Sender        common.Address
	WorkflowOwner common.Address
	WorkflowName  [10]byte
}

type NewRound struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
}

type NewRoundDecoded struct {
	RoundId   *big.Int
	StartedBy common.Address
	StartedAt *big.Int
}

type OwnershipTransferRequested struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferRequestedDecoded struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferred struct {
	From common.Address
	To   common.Address
}

type OwnershipTransferredDecoded struct {
	From common.Address
	To   common.Address
}

type ProxyDataIdRemoved struct {
	Proxy  common.Address
	DataId [16]byte
}

type ProxyDataIdRemovedDecoded struct {
	Proxy  common.Address
	DataId [16]byte
}

type ProxyDataIdUpdated struct {
	Proxy  common.Address
	DataId [16]byte
}

type ProxyDataIdUpdatedDecoded struct {
	Proxy  common.Address
	DataId [16]byte
}

type StaleBundleReport struct {
	DataId          [16]byte
	ReportTimestamp *big.Int
	LatestTimestamp *big.Int
}

type StaleBundleReportDecoded struct {
	DataId          [16]byte
	ReportTimestamp *big.Int
	LatestTimestamp *big.Int
}

type StaleDecimalReport struct {
	DataId          [16]byte
	ReportTimestamp *big.Int
	LatestTimestamp *big.Int
}

type StaleDecimalReportDecoded struct {
	DataId          [16]byte
	ReportTimestamp *big.Int
	LatestTimestamp *big.Int
}

type TokenRecovered struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
}

type TokenRecoveredDecoded struct {
	Token  common.Address
	To     common.Address
	Amount *big.Int
}

// Main Binding Type for DataFeedsCache
type DataFeedsCache struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   DataFeedsCacheCodec
}

type DataFeedsCacheCodec interface {
	EncodeAcceptOwnershipMethodCall() ([]byte, error)
	EncodeBundleDecimalsMethodCall() ([]byte, error)
	DecodeBundleDecimalsMethodOutput(data []byte) ([]uint8, error)
	EncodeCheckFeedPermissionMethodCall(in CheckFeedPermissionInput) ([]byte, error)
	DecodeCheckFeedPermissionMethodOutput(data []byte) (bool, error)
	EncodeDecimalsMethodCall() ([]byte, error)
	DecodeDecimalsMethodOutput(data []byte) (uint8, error)
	EncodeDescriptionMethodCall() ([]byte, error)
	DecodeDescriptionMethodOutput(data []byte) (string, error)
	EncodeGetAnswerMethodCall(in GetAnswerInput) ([]byte, error)
	DecodeGetAnswerMethodOutput(data []byte) (*big.Int, error)
	EncodeGetBundleDecimalsMethodCall(in GetBundleDecimalsInput) ([]byte, error)
	DecodeGetBundleDecimalsMethodOutput(data []byte) ([]uint8, error)
	EncodeGetDataIdForProxyMethodCall(in GetDataIdForProxyInput) ([]byte, error)
	DecodeGetDataIdForProxyMethodOutput(data []byte) ([16]byte, error)
	EncodeGetDecimalsMethodCall(in GetDecimalsInput) ([]byte, error)
	DecodeGetDecimalsMethodOutput(data []byte) (uint8, error)
	EncodeGetDescriptionMethodCall(in GetDescriptionInput) ([]byte, error)
	DecodeGetDescriptionMethodOutput(data []byte) (string, error)
	EncodeGetFeedMetadataMethodCall(in GetFeedMetadataInput) ([]byte, error)
	DecodeGetFeedMetadataMethodOutput(data []byte) ([]WorkflowMetadata, error)
	EncodeGetLatestAnswerMethodCall(in GetLatestAnswerInput) ([]byte, error)
	DecodeGetLatestAnswerMethodOutput(data []byte) (*big.Int, error)
	EncodeGetLatestBundleMethodCall(in GetLatestBundleInput) ([]byte, error)
	DecodeGetLatestBundleMethodOutput(data []byte) ([]byte, error)
	EncodeGetLatestBundleTimestampMethodCall(in GetLatestBundleTimestampInput) ([]byte, error)
	DecodeGetLatestBundleTimestampMethodOutput(data []byte) (*big.Int, error)
	EncodeGetLatestRoundDataMethodCall(in GetLatestRoundDataInput) ([]byte, error)
	DecodeGetLatestRoundDataMethodOutput(data []byte) (GetLatestRoundDataOutput, error)
	EncodeGetLatestTimestampMethodCall(in GetLatestTimestampInput) ([]byte, error)
	DecodeGetLatestTimestampMethodOutput(data []byte) (*big.Int, error)
	EncodeGetRoundDataMethodCall(in GetRoundDataInput) ([]byte, error)
	DecodeGetRoundDataMethodOutput(data []byte) (GetRoundDataOutput, error)
	EncodeGetTimestampMethodCall(in GetTimestampInput) ([]byte, error)
	DecodeGetTimestampMethodOutput(data []byte) (*big.Int, error)
	EncodeIsFeedAdminMethodCall(in IsFeedAdminInput) ([]byte, error)
	DecodeIsFeedAdminMethodOutput(data []byte) (bool, error)
	EncodeLatestAnswerMethodCall() ([]byte, error)
	DecodeLatestAnswerMethodOutput(data []byte) (*big.Int, error)
	EncodeLatestBundleMethodCall() ([]byte, error)
	DecodeLatestBundleMethodOutput(data []byte) ([]byte, error)
	EncodeLatestBundleTimestampMethodCall() ([]byte, error)
	DecodeLatestBundleTimestampMethodOutput(data []byte) (*big.Int, error)
	EncodeLatestRoundMethodCall() ([]byte, error)
	DecodeLatestRoundMethodOutput(data []byte) (*big.Int, error)
	EncodeLatestRoundDataMethodCall() ([]byte, error)
	DecodeLatestRoundDataMethodOutput(data []byte) (LatestRoundDataOutput, error)
	EncodeLatestTimestampMethodCall() ([]byte, error)
	DecodeLatestTimestampMethodOutput(data []byte) (*big.Int, error)
	EncodeOnReportMethodCall(in OnReportInput) ([]byte, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodeRecoverTokensMethodCall(in RecoverTokensInput) ([]byte, error)
	EncodeRemoveDataIdMappingsForProxiesMethodCall(in RemoveDataIdMappingsForProxiesInput) ([]byte, error)
	EncodeRemoveFeedConfigsMethodCall(in RemoveFeedConfigsInput) ([]byte, error)
	EncodeSetBundleFeedConfigsMethodCall(in SetBundleFeedConfigsInput) ([]byte, error)
	EncodeSetDecimalFeedConfigsMethodCall(in SetDecimalFeedConfigsInput) ([]byte, error)
	EncodeSetFeedAdminMethodCall(in SetFeedAdminInput) ([]byte, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	EncodeTypeAndVersionMethodCall() ([]byte, error)
	DecodeTypeAndVersionMethodOutput(data []byte) (string, error)
	EncodeUpdateDataIdMappingsForProxiesMethodCall(in UpdateDataIdMappingsForProxiesInput) ([]byte, error)
	EncodeVersionMethodCall() ([]byte, error)
	DecodeVersionMethodOutput(data []byte) (*big.Int, error)
	EncodeWorkflowMetadataStruct(in WorkflowMetadata) ([]byte, error)
	AnswerUpdatedLogHash() []byte
	EncodeAnswerUpdatedTopics(evt abi.Event, values []AnswerUpdated) ([]*evm.TopicValues, error)
	DecodeAnswerUpdated(log *evm.Log) (*AnswerUpdatedDecoded, error)
	BundleFeedConfigSetLogHash() []byte
	EncodeBundleFeedConfigSetTopics(evt abi.Event, values []BundleFeedConfigSet) ([]*evm.TopicValues, error)
	DecodeBundleFeedConfigSet(log *evm.Log) (*BundleFeedConfigSetDecoded, error)
	BundleReportUpdatedLogHash() []byte
	EncodeBundleReportUpdatedTopics(evt abi.Event, values []BundleReportUpdated) ([]*evm.TopicValues, error)
	DecodeBundleReportUpdated(log *evm.Log) (*BundleReportUpdatedDecoded, error)
	DecimalFeedConfigSetLogHash() []byte
	EncodeDecimalFeedConfigSetTopics(evt abi.Event, values []DecimalFeedConfigSet) ([]*evm.TopicValues, error)
	DecodeDecimalFeedConfigSet(log *evm.Log) (*DecimalFeedConfigSetDecoded, error)
	DecimalReportUpdatedLogHash() []byte
	EncodeDecimalReportUpdatedTopics(evt abi.Event, values []DecimalReportUpdated) ([]*evm.TopicValues, error)
	DecodeDecimalReportUpdated(log *evm.Log) (*DecimalReportUpdatedDecoded, error)
	FeedAdminSetLogHash() []byte
	EncodeFeedAdminSetTopics(evt abi.Event, values []FeedAdminSet) ([]*evm.TopicValues, error)
	DecodeFeedAdminSet(log *evm.Log) (*FeedAdminSetDecoded, error)
	FeedConfigRemovedLogHash() []byte
	EncodeFeedConfigRemovedTopics(evt abi.Event, values []FeedConfigRemoved) ([]*evm.TopicValues, error)
	DecodeFeedConfigRemoved(log *evm.Log) (*FeedConfigRemovedDecoded, error)
	InvalidUpdatePermissionLogHash() []byte
	EncodeInvalidUpdatePermissionTopics(evt abi.Event, values []InvalidUpdatePermission) ([]*evm.TopicValues, error)
	DecodeInvalidUpdatePermission(log *evm.Log) (*InvalidUpdatePermissionDecoded, error)
	NewRoundLogHash() []byte
	EncodeNewRoundTopics(evt abi.Event, values []NewRound) ([]*evm.TopicValues, error)
	DecodeNewRound(log *evm.Log) (*NewRoundDecoded, error)
	OwnershipTransferRequestedLogHash() []byte
	EncodeOwnershipTransferRequestedTopics(evt abi.Event, values []OwnershipTransferRequested) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferRequested(log *evm.Log) (*OwnershipTransferRequestedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferred) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
	ProxyDataIdRemovedLogHash() []byte
	EncodeProxyDataIdRemovedTopics(evt abi.Event, values []ProxyDataIdRemoved) ([]*evm.TopicValues, error)
	DecodeProxyDataIdRemoved(log *evm.Log) (*ProxyDataIdRemovedDecoded, error)
	ProxyDataIdUpdatedLogHash() []byte
	EncodeProxyDataIdUpdatedTopics(evt abi.Event, values []ProxyDataIdUpdated) ([]*evm.TopicValues, error)
	DecodeProxyDataIdUpdated(log *evm.Log) (*ProxyDataIdUpdatedDecoded, error)
	StaleBundleReportLogHash() []byte
	EncodeStaleBundleReportTopics(evt abi.Event, values []StaleBundleReport) ([]*evm.TopicValues, error)
	DecodeStaleBundleReport(log *evm.Log) (*StaleBundleReportDecoded, error)
	StaleDecimalReportLogHash() []byte
	EncodeStaleDecimalReportTopics(evt abi.Event, values []StaleDecimalReport) ([]*evm.TopicValues, error)
	DecodeStaleDecimalReport(log *evm.Log) (*StaleDecimalReportDecoded, error)
	TokenRecoveredLogHash() []byte
	EncodeTokenRecoveredTopics(evt abi.Event, values []TokenRecovered) ([]*evm.TopicValues, error)
	DecodeTokenRecovered(log *evm.Log) (*TokenRecoveredDecoded, error)
}

func NewDataFeedsCache(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*DataFeedsCache, error) {
	parsed, err := abi.JSON(strings.NewReader(DataFeedsCacheMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &DataFeedsCache{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   codec,
	}, nil
}

type Codec struct {
	abi *abi.ABI
}

func NewCodec() (DataFeedsCacheCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(DataFeedsCacheMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeAcceptOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("acceptOwnership")
}

func (c *Codec) EncodeBundleDecimalsMethodCall() ([]byte, error) {
	return c.abi.Pack("bundleDecimals")
}

func (c *Codec) DecodeBundleDecimalsMethodOutput(data []byte) ([]uint8, error) {
	vals, err := c.abi.Methods["bundleDecimals"].Outputs.Unpack(data)
	if err != nil {
		return *new([]uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]uint8), fmt.Errorf("failed to unmarshal to []uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeCheckFeedPermissionMethodCall(in CheckFeedPermissionInput) ([]byte, error) {
	return c.abi.Pack("checkFeedPermission", in.DataId, in.WorkflowMetadata)
}

func (c *Codec) DecodeCheckFeedPermissionMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["checkFeedPermission"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeDecimalsMethodCall() ([]byte, error) {
	return c.abi.Pack("decimals")
}

func (c *Codec) DecodeDecimalsMethodOutput(data []byte) (uint8, error) {
	vals, err := c.abi.Methods["decimals"].Outputs.Unpack(data)
	if err != nil {
		return *new(uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(uint8), fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeDescriptionMethodCall() ([]byte, error) {
	return c.abi.Pack("description")
}

func (c *Codec) DecodeDescriptionMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["description"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetAnswerMethodCall(in GetAnswerInput) ([]byte, error) {
	return c.abi.Pack("getAnswer", in.RoundId)
}

func (c *Codec) DecodeGetAnswerMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getAnswer"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetBundleDecimalsMethodCall(in GetBundleDecimalsInput) ([]byte, error) {
	return c.abi.Pack("getBundleDecimals", in.DataId)
}

func (c *Codec) DecodeGetBundleDecimalsMethodOutput(data []byte) ([]uint8, error) {
	vals, err := c.abi.Methods["getBundleDecimals"].Outputs.Unpack(data)
	if err != nil {
		return *new([]uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]uint8), fmt.Errorf("failed to unmarshal to []uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetDataIdForProxyMethodCall(in GetDataIdForProxyInput) ([]byte, error) {
	return c.abi.Pack("getDataIdForProxy", in.Proxy)
}

func (c *Codec) DecodeGetDataIdForProxyMethodOutput(data []byte) ([16]byte, error) {
	vals, err := c.abi.Methods["getDataIdForProxy"].Outputs.Unpack(data)
	if err != nil {
		return *new([16]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([16]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result [16]byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([16]byte), fmt.Errorf("failed to unmarshal to [16]byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetDecimalsMethodCall(in GetDecimalsInput) ([]byte, error) {
	return c.abi.Pack("getDecimals", in.DataId)
}

func (c *Codec) DecodeGetDecimalsMethodOutput(data []byte) (uint8, error) {
	vals, err := c.abi.Methods["getDecimals"].Outputs.Unpack(data)
	if err != nil {
		return *new(uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(uint8), fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetDescriptionMethodCall(in GetDescriptionInput) ([]byte, error) {
	return c.abi.Pack("getDescription", in.DataId)
}

func (c *Codec) DecodeGetDescriptionMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["getDescription"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetFeedMetadataMethodCall(in GetFeedMetadataInput) ([]byte, error) {
	return c.abi.Pack("getFeedMetadata", in.DataId, in.StartIndex, in.MaxCount)
}

func (c *Codec) DecodeGetFeedMetadataMethodOutput(data []byte) ([]WorkflowMetadata, error) {
	vals, err := c.abi.Methods["getFeedMetadata"].Outputs.Unpack(data)
	if err != nil {
		return *new([]WorkflowMetadata), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]WorkflowMetadata), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []WorkflowMetadata
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]WorkflowMetadata), fmt.Errorf("failed to unmarshal to []WorkflowMetadata: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetLatestAnswerMethodCall(in GetLatestAnswerInput) ([]byte, error) {
	return c.abi.Pack("getLatestAnswer", in.DataId)
}

func (c *Codec) DecodeGetLatestAnswerMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getLatestAnswer"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetLatestBundleMethodCall(in GetLatestBundleInput) ([]byte, error) {
	return c.abi.Pack("getLatestBundle", in.DataId)
}

func (c *Codec) DecodeGetLatestBundleMethodOutput(data []byte) ([]byte, error) {
	vals, err := c.abi.Methods["getLatestBundle"].Outputs.Unpack(data)
	if err != nil {
		return *new([]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]byte), fmt.Errorf("failed to unmarshal to []byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetLatestBundleTimestampMethodCall(in GetLatestBundleTimestampInput) ([]byte, error) {
	return c.abi.Pack("getLatestBundleTimestamp", in.DataId)
}

func (c *Codec) DecodeGetLatestBundleTimestampMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getLatestBundleTimestamp"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetLatestRoundDataMethodCall(in GetLatestRoundDataInput) ([]byte, error) {
	return c.abi.Pack("getLatestRoundData", in.DataId)
}

func (c *Codec) DecodeGetLatestRoundDataMethodOutput(data []byte) (GetLatestRoundDataOutput, error) {
	vals, err := c.abi.Methods["getLatestRoundData"].Outputs.Unpack(data)
	if err != nil {
		return GetLatestRoundDataOutput{}, err
	}
	if len(vals) != 5 {
		return GetLatestRoundDataOutput{}, fmt.Errorf("expected 5 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 *big.Int
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 *big.Int
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return GetLatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return GetLatestRoundDataOutput{
		Id:              result0,
		Answer:          result1,
		StartedAt:       result2,
		UpdatedAt:       result3,
		AnsweredInRound: result4,
	}, nil
}

func (c *Codec) EncodeGetLatestTimestampMethodCall(in GetLatestTimestampInput) ([]byte, error) {
	return c.abi.Pack("getLatestTimestamp", in.DataId)
}

func (c *Codec) DecodeGetLatestTimestampMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getLatestTimestamp"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetRoundDataMethodCall(in GetRoundDataInput) ([]byte, error) {
	return c.abi.Pack("getRoundData", in.RoundId)
}

func (c *Codec) DecodeGetRoundDataMethodOutput(data []byte) (GetRoundDataOutput, error) {
	vals, err := c.abi.Methods["getRoundData"].Outputs.Unpack(data)
	if err != nil {
		return GetRoundDataOutput{}, err
	}
	if len(vals) != 5 {
		return GetRoundDataOutput{}, fmt.Errorf("expected 5 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 *big.Int
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 *big.Int
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return GetRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return GetRoundDataOutput{
		Id:              result0,
		Answer:          result1,
		StartedAt:       result2,
		UpdatedAt:       result3,
		AnsweredInRound: result4,
	}, nil
}

func (c *Codec) EncodeGetTimestampMethodCall(in GetTimestampInput) ([]byte, error) {
	return c.abi.Pack("getTimestamp", in.RoundId)
}

func (c *Codec) DecodeGetTimestampMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getTimestamp"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeIsFeedAdminMethodCall(in IsFeedAdminInput) ([]byte, error) {
	return c.abi.Pack("isFeedAdmin", in.FeedAdmin)
}

func (c *Codec) DecodeIsFeedAdminMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["isFeedAdmin"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeLatestAnswerMethodCall() ([]byte, error) {
	return c.abi.Pack("latestAnswer")
}

func (c *Codec) DecodeLatestAnswerMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["latestAnswer"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeLatestBundleMethodCall() ([]byte, error) {
	return c.abi.Pack("latestBundle")
}

func (c *Codec) DecodeLatestBundleMethodOutput(data []byte) ([]byte, error) {
	vals, err := c.abi.Methods["latestBundle"].Outputs.Unpack(data)
	if err != nil {
		return *new([]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]byte), fmt.Errorf("failed to unmarshal to []byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeLatestBundleTimestampMethodCall() ([]byte, error) {
	return c.abi.Pack("latestBundleTimestamp")
}

func (c *Codec) DecodeLatestBundleTimestampMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["latestBundleTimestamp"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeLatestRoundMethodCall() ([]byte, error) {
	return c.abi.Pack("latestRound")
}

func (c *Codec) DecodeLatestRoundMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["latestRound"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeLatestRoundDataMethodCall() ([]byte, error) {
	return c.abi.Pack("latestRoundData")
}

func (c *Codec) DecodeLatestRoundDataMethodOutput(data []byte) (LatestRoundDataOutput, error) {
	vals, err := c.abi.Methods["latestRoundData"].Outputs.Unpack(data)
	if err != nil {
		return LatestRoundDataOutput{}, err
	}
	if len(vals) != 5 {
		return LatestRoundDataOutput{}, fmt.Errorf("expected 5 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 *big.Int
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 *big.Int
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return LatestRoundDataOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return LatestRoundDataOutput{
		Id:              result0,
		Answer:          result1,
		StartedAt:       result2,
		UpdatedAt:       result3,
		AnsweredInRound: result4,
	}, nil
}

func (c *Codec) EncodeLatestTimestampMethodCall() ([]byte, error) {
	return c.abi.Pack("latestTimestamp")
}

func (c *Codec) DecodeLatestTimestampMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["latestTimestamp"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeOnReportMethodCall(in OnReportInput) ([]byte, error) {
	return c.abi.Pack("onReport", in.Metadata, in.Report)
}

func (c *Codec) EncodeOwnerMethodCall() ([]byte, error) {
	return c.abi.Pack("owner")
}

func (c *Codec) DecodeOwnerMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["owner"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeRecoverTokensMethodCall(in RecoverTokensInput) ([]byte, error) {
	return c.abi.Pack("recoverTokens", in.Token, in.To, in.Amount)
}

func (c *Codec) EncodeRemoveDataIdMappingsForProxiesMethodCall(in RemoveDataIdMappingsForProxiesInput) ([]byte, error) {
	return c.abi.Pack("removeDataIdMappingsForProxies", in.Proxies)
}

func (c *Codec) EncodeRemoveFeedConfigsMethodCall(in RemoveFeedConfigsInput) ([]byte, error) {
	return c.abi.Pack("removeFeedConfigs", in.DataIds)
}

func (c *Codec) EncodeSetBundleFeedConfigsMethodCall(in SetBundleFeedConfigsInput) ([]byte, error) {
	return c.abi.Pack("setBundleFeedConfigs", in.DataIds, in.Descriptions, in.DecimalsMatrix, in.WorkflowMetadata)
}

func (c *Codec) EncodeSetDecimalFeedConfigsMethodCall(in SetDecimalFeedConfigsInput) ([]byte, error) {
	return c.abi.Pack("setDecimalFeedConfigs", in.DataIds, in.Descriptions, in.WorkflowMetadata)
}

func (c *Codec) EncodeSetFeedAdminMethodCall(in SetFeedAdminInput) ([]byte, error) {
	return c.abi.Pack("setFeedAdmin", in.FeedAdmin, in.IsAdmin)
}

func (c *Codec) EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error) {
	return c.abi.Pack("supportsInterface", in.InterfaceId)
}

func (c *Codec) DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["supportsInterface"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.To)
}

func (c *Codec) EncodeTypeAndVersionMethodCall() ([]byte, error) {
	return c.abi.Pack("typeAndVersion")
}

func (c *Codec) DecodeTypeAndVersionMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["typeAndVersion"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeUpdateDataIdMappingsForProxiesMethodCall(in UpdateDataIdMappingsForProxiesInput) ([]byte, error) {
	return c.abi.Pack("updateDataIdMappingsForProxies", in.Proxies, in.DataIds)
}

func (c *Codec) EncodeVersionMethodCall() ([]byte, error) {
	return c.abi.Pack("version")
}

func (c *Codec) DecodeVersionMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["version"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeWorkflowMetadataStruct(in WorkflowMetadata) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "allowedSender", Type: "address"},
			{Name: "allowedWorkflowOwner", Type: "address"},
			{Name: "allowedWorkflowName", Type: "bytes10"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for WorkflowMetadata: %w", err)
	}
	args := abi.Arguments{
		{Name: "workflowMetadata", Type: tupleType},
	}

	return args.Pack(in)
}

func (c *Codec) AnswerUpdatedLogHash() []byte {
	return c.abi.Events["AnswerUpdated"].ID.Bytes()
}

func (c *Codec) EncodeAnswerUpdatedTopics(
	evt abi.Event,
	values []AnswerUpdated,
) ([]*evm.TopicValues, error) {
	var currentRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Current)
		if err != nil {
			return nil, err
		}
		currentRule = append(currentRule, fieldVal)
	}
	var roundIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.RoundId)
		if err != nil {
			return nil, err
		}
		roundIdRule = append(roundIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		currentRule,
		roundIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeAnswerUpdated decodes a log into a AnswerUpdated struct.
func (c *Codec) DecodeAnswerUpdated(log *evm.Log) (*AnswerUpdatedDecoded, error) {
	event := new(AnswerUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "AnswerUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["AnswerUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) BundleFeedConfigSetLogHash() []byte {
	return c.abi.Events["BundleFeedConfigSet"].ID.Bytes()
}

func (c *Codec) EncodeBundleFeedConfigSetTopics(
	evt abi.Event,
	values []BundleFeedConfigSet,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeBundleFeedConfigSet decodes a log into a BundleFeedConfigSet struct.
func (c *Codec) DecodeBundleFeedConfigSet(log *evm.Log) (*BundleFeedConfigSetDecoded, error) {
	event := new(BundleFeedConfigSetDecoded)
	if err := c.abi.UnpackIntoInterface(event, "BundleFeedConfigSet", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["BundleFeedConfigSet"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) BundleReportUpdatedLogHash() []byte {
	return c.abi.Events["BundleReportUpdated"].ID.Bytes()
}

func (c *Codec) EncodeBundleReportUpdatedTopics(
	evt abi.Event,
	values []BundleReportUpdated,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}
	var timestampRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Timestamp)
		if err != nil {
			return nil, err
		}
		timestampRule = append(timestampRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
		timestampRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeBundleReportUpdated decodes a log into a BundleReportUpdated struct.
func (c *Codec) DecodeBundleReportUpdated(log *evm.Log) (*BundleReportUpdatedDecoded, error) {
	event := new(BundleReportUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "BundleReportUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["BundleReportUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) DecimalFeedConfigSetLogHash() []byte {
	return c.abi.Events["DecimalFeedConfigSet"].ID.Bytes()
}

func (c *Codec) EncodeDecimalFeedConfigSetTopics(
	evt abi.Event,
	values []DecimalFeedConfigSet,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeDecimalFeedConfigSet decodes a log into a DecimalFeedConfigSet struct.
func (c *Codec) DecodeDecimalFeedConfigSet(log *evm.Log) (*DecimalFeedConfigSetDecoded, error) {
	event := new(DecimalFeedConfigSetDecoded)
	if err := c.abi.UnpackIntoInterface(event, "DecimalFeedConfigSet", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["DecimalFeedConfigSet"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) DecimalReportUpdatedLogHash() []byte {
	return c.abi.Events["DecimalReportUpdated"].ID.Bytes()
}

func (c *Codec) EncodeDecimalReportUpdatedTopics(
	evt abi.Event,
	values []DecimalReportUpdated,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}
	var roundIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.RoundId)
		if err != nil {
			return nil, err
		}
		roundIdRule = append(roundIdRule, fieldVal)
	}
	var timestampRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.Timestamp)
		if err != nil {
			return nil, err
		}
		timestampRule = append(timestampRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
		roundIdRule,
		timestampRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeDecimalReportUpdated decodes a log into a DecimalReportUpdated struct.
func (c *Codec) DecodeDecimalReportUpdated(log *evm.Log) (*DecimalReportUpdatedDecoded, error) {
	event := new(DecimalReportUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "DecimalReportUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["DecimalReportUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) FeedAdminSetLogHash() []byte {
	return c.abi.Events["FeedAdminSet"].ID.Bytes()
}

func (c *Codec) EncodeFeedAdminSetTopics(
	evt abi.Event,
	values []FeedAdminSet,
) ([]*evm.TopicValues, error) {
	var feedAdminRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.FeedAdmin)
		if err != nil {
			return nil, err
		}
		feedAdminRule = append(feedAdminRule, fieldVal)
	}
	var isAdminRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.IsAdmin)
		if err != nil {
			return nil, err
		}
		isAdminRule = append(isAdminRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		feedAdminRule,
		isAdminRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeFeedAdminSet decodes a log into a FeedAdminSet struct.
func (c *Codec) DecodeFeedAdminSet(log *evm.Log) (*FeedAdminSetDecoded, error) {
	event := new(FeedAdminSetDecoded)
	if err := c.abi.UnpackIntoInterface(event, "FeedAdminSet", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["FeedAdminSet"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) FeedConfigRemovedLogHash() []byte {
	return c.abi.Events["FeedConfigRemoved"].ID.Bytes()
}

func (c *Codec) EncodeFeedConfigRemovedTopics(
	evt abi.Event,
	values []FeedConfigRemoved,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeFeedConfigRemoved decodes a log into a FeedConfigRemoved struct.
func (c *Codec) DecodeFeedConfigRemoved(log *evm.Log) (*FeedConfigRemovedDecoded, error) {
	event := new(FeedConfigRemovedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "FeedConfigRemoved", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["FeedConfigRemoved"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) InvalidUpdatePermissionLogHash() []byte {
	return c.abi.Events["InvalidUpdatePermission"].ID.Bytes()
}

func (c *Codec) EncodeInvalidUpdatePermissionTopics(
	evt abi.Event,
	values []InvalidUpdatePermission,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeInvalidUpdatePermission decodes a log into a InvalidUpdatePermission struct.
func (c *Codec) DecodeInvalidUpdatePermission(log *evm.Log) (*InvalidUpdatePermissionDecoded, error) {
	event := new(InvalidUpdatePermissionDecoded)
	if err := c.abi.UnpackIntoInterface(event, "InvalidUpdatePermission", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["InvalidUpdatePermission"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) NewRoundLogHash() []byte {
	return c.abi.Events["NewRound"].ID.Bytes()
}

func (c *Codec) EncodeNewRoundTopics(
	evt abi.Event,
	values []NewRound,
) ([]*evm.TopicValues, error) {
	var roundIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.RoundId)
		if err != nil {
			return nil, err
		}
		roundIdRule = append(roundIdRule, fieldVal)
	}
	var startedByRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.StartedBy)
		if err != nil {
			return nil, err
		}
		startedByRule = append(startedByRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		roundIdRule,
		startedByRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeNewRound decodes a log into a NewRound struct.
func (c *Codec) DecodeNewRound(log *evm.Log) (*NewRoundDecoded, error) {
	event := new(NewRoundDecoded)
	if err := c.abi.UnpackIntoInterface(event, "NewRound", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["NewRound"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) OwnershipTransferRequestedLogHash() []byte {
	return c.abi.Events["OwnershipTransferRequested"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferRequestedTopics(
	evt abi.Event,
	values []OwnershipTransferRequested,
) ([]*evm.TopicValues, error) {
	var fromRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.From)
		if err != nil {
			return nil, err
		}
		fromRule = append(fromRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.To)
		if err != nil {
			return nil, err
		}
		toRule = append(toRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		fromRule,
		toRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeOwnershipTransferRequested decodes a log into a OwnershipTransferRequested struct.
func (c *Codec) DecodeOwnershipTransferRequested(log *evm.Log) (*OwnershipTransferRequestedDecoded, error) {
	event := new(OwnershipTransferRequestedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferRequested", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferRequested"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) OwnershipTransferredLogHash() []byte {
	return c.abi.Events["OwnershipTransferred"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferredTopics(
	evt abi.Event,
	values []OwnershipTransferred,
) ([]*evm.TopicValues, error) {
	var fromRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.From)
		if err != nil {
			return nil, err
		}
		fromRule = append(fromRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.To)
		if err != nil {
			return nil, err
		}
		toRule = append(toRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		fromRule,
		toRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeOwnershipTransferred decodes a log into a OwnershipTransferred struct.
func (c *Codec) DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error) {
	event := new(OwnershipTransferredDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferred", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferred"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ProxyDataIdRemovedLogHash() []byte {
	return c.abi.Events["ProxyDataIdRemoved"].ID.Bytes()
}

func (c *Codec) EncodeProxyDataIdRemovedTopics(
	evt abi.Event,
	values []ProxyDataIdRemoved,
) ([]*evm.TopicValues, error) {
	var proxyRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Proxy)
		if err != nil {
			return nil, err
		}
		proxyRule = append(proxyRule, fieldVal)
	}
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		proxyRule,
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeProxyDataIdRemoved decodes a log into a ProxyDataIdRemoved struct.
func (c *Codec) DecodeProxyDataIdRemoved(log *evm.Log) (*ProxyDataIdRemovedDecoded, error) {
	event := new(ProxyDataIdRemovedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ProxyDataIdRemoved", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ProxyDataIdRemoved"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ProxyDataIdUpdatedLogHash() []byte {
	return c.abi.Events["ProxyDataIdUpdated"].ID.Bytes()
}

func (c *Codec) EncodeProxyDataIdUpdatedTopics(
	evt abi.Event,
	values []ProxyDataIdUpdated,
) ([]*evm.TopicValues, error) {
	var proxyRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Proxy)
		if err != nil {
			return nil, err
		}
		proxyRule = append(proxyRule, fieldVal)
	}
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		proxyRule,
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeProxyDataIdUpdated decodes a log into a ProxyDataIdUpdated struct.
func (c *Codec) DecodeProxyDataIdUpdated(log *evm.Log) (*ProxyDataIdUpdatedDecoded, error) {
	event := new(ProxyDataIdUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ProxyDataIdUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ProxyDataIdUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) StaleBundleReportLogHash() []byte {
	return c.abi.Events["StaleBundleReport"].ID.Bytes()
}

func (c *Codec) EncodeStaleBundleReportTopics(
	evt abi.Event,
	values []StaleBundleReport,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeStaleBundleReport decodes a log into a StaleBundleReport struct.
func (c *Codec) DecodeStaleBundleReport(log *evm.Log) (*StaleBundleReportDecoded, error) {
	event := new(StaleBundleReportDecoded)
	if err := c.abi.UnpackIntoInterface(event, "StaleBundleReport", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["StaleBundleReport"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) StaleDecimalReportLogHash() []byte {
	return c.abi.Events["StaleDecimalReport"].ID.Bytes()
}

func (c *Codec) EncodeStaleDecimalReportTopics(
	evt abi.Event,
	values []StaleDecimalReport,
) ([]*evm.TopicValues, error) {
	var dataIdRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.DataId)
		if err != nil {
			return nil, err
		}
		dataIdRule = append(dataIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		dataIdRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeStaleDecimalReport decodes a log into a StaleDecimalReport struct.
func (c *Codec) DecodeStaleDecimalReport(log *evm.Log) (*StaleDecimalReportDecoded, error) {
	event := new(StaleDecimalReportDecoded)
	if err := c.abi.UnpackIntoInterface(event, "StaleDecimalReport", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["StaleDecimalReport"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) TokenRecoveredLogHash() []byte {
	return c.abi.Events["TokenRecovered"].ID.Bytes()
}

func (c *Codec) EncodeTokenRecoveredTopics(
	evt abi.Event,
	values []TokenRecovered,
) ([]*evm.TopicValues, error) {
	var tokenRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Token)
		if err != nil {
			return nil, err
		}
		tokenRule = append(tokenRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.To)
		if err != nil {
			return nil, err
		}
		toRule = append(toRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenRule,
		toRule,
	)
	if err != nil {
		return nil, err
	}

	topics := make([]*evm.TopicValues, len(rawTopics)+1)
	topics[0] = &evm.TopicValues{
		Values: [][]byte{evt.ID.Bytes()},
	}
	for i, hashList := range rawTopics {
		bs := make([][]byte, len(hashList))
		for j, h := range hashList {
			bs[j] = h.Bytes()
		}
		topics[i+1] = &evm.TopicValues{Values: bs}
	}
	return topics, nil
}

// DecodeTokenRecovered decodes a log into a TokenRecovered struct.
func (c *Codec) DecodeTokenRecovered(log *evm.Log) (*TokenRecoveredDecoded, error) {
	event := new(TokenRecoveredDecoded)
	if err := c.abi.UnpackIntoInterface(event, "TokenRecovered", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["TokenRecovered"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c DataFeedsCache) BundleDecimals(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[]uint8] {
	calldata, err := c.Codec.EncodeBundleDecimalsMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[]uint8](*new([]uint8), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([]uint8, error) {
		return c.Codec.DecodeBundleDecimalsMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) CheckFeedPermission(
	runtime cre.Runtime,
	args CheckFeedPermissionInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeCheckFeedPermissionMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeCheckFeedPermissionMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) Decimals(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[uint8] {
	calldata, err := c.Codec.EncodeDecimalsMethodCall()
	if err != nil {
		return cre.PromiseFromResult[uint8](*new(uint8), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (uint8, error) {
		return c.Codec.DecodeDecimalsMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) Description(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeDescriptionMethodCall()
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeDescriptionMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetAnswer(
	runtime cre.Runtime,
	args GetAnswerInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetAnswerMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetAnswerMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetBundleDecimals(
	runtime cre.Runtime,
	args GetBundleDecimalsInput,
	blockNumber *big.Int,
) cre.Promise[[]uint8] {
	calldata, err := c.Codec.EncodeGetBundleDecimalsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[[]uint8](*new([]uint8), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([]uint8, error) {
		return c.Codec.DecodeGetBundleDecimalsMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetDataIdForProxy(
	runtime cre.Runtime,
	args GetDataIdForProxyInput,
	blockNumber *big.Int,
) cre.Promise[[16]byte] {
	calldata, err := c.Codec.EncodeGetDataIdForProxyMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[[16]byte](*new([16]byte), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([16]byte, error) {
		return c.Codec.DecodeGetDataIdForProxyMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetDescription(
	runtime cre.Runtime,
	args GetDescriptionInput,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeGetDescriptionMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeGetDescriptionMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetFeedMetadata(
	runtime cre.Runtime,
	args GetFeedMetadataInput,
	blockNumber *big.Int,
) cre.Promise[[]WorkflowMetadata] {
	calldata, err := c.Codec.EncodeGetFeedMetadataMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[[]WorkflowMetadata](*new([]WorkflowMetadata), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([]WorkflowMetadata, error) {
		return c.Codec.DecodeGetFeedMetadataMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetLatestAnswer(
	runtime cre.Runtime,
	args GetLatestAnswerInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetLatestAnswerMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetLatestAnswerMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetLatestBundle(
	runtime cre.Runtime,
	args GetLatestBundleInput,
	blockNumber *big.Int,
) cre.Promise[[]byte] {
	calldata, err := c.Codec.EncodeGetLatestBundleMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[[]byte](*new([]byte), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([]byte, error) {
		return c.Codec.DecodeGetLatestBundleMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetLatestBundleTimestamp(
	runtime cre.Runtime,
	args GetLatestBundleTimestampInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetLatestBundleTimestampMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetLatestBundleTimestampMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetLatestRoundData(
	runtime cre.Runtime,
	args GetLatestRoundDataInput,
	blockNumber *big.Int,
) cre.Promise[GetLatestRoundDataOutput] {
	calldata, err := c.Codec.EncodeGetLatestRoundDataMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[GetLatestRoundDataOutput](GetLatestRoundDataOutput{}, err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (GetLatestRoundDataOutput, error) {
		return c.Codec.DecodeGetLatestRoundDataMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetLatestTimestamp(
	runtime cre.Runtime,
	args GetLatestTimestampInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetLatestTimestampMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetLatestTimestampMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetRoundData(
	runtime cre.Runtime,
	args GetRoundDataInput,
	blockNumber *big.Int,
) cre.Promise[GetRoundDataOutput] {
	calldata, err := c.Codec.EncodeGetRoundDataMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[GetRoundDataOutput](GetRoundDataOutput{}, err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (GetRoundDataOutput, error) {
		return c.Codec.DecodeGetRoundDataMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) GetTimestamp(
	runtime cre.Runtime,
	args GetTimestampInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetTimestampMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetTimestampMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) IsFeedAdmin(
	runtime cre.Runtime,
	args IsFeedAdminInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeIsFeedAdminMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeIsFeedAdminMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) LatestAnswer(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLatestAnswerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeLatestAnswerMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) LatestBundle(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[]byte] {
	calldata, err := c.Codec.EncodeLatestBundleMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[]byte](*new([]byte), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) ([]byte, error) {
		return c.Codec.DecodeLatestBundleMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) LatestBundleTimestamp(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLatestBundleTimestampMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeLatestBundleTimestampMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) LatestRound(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLatestRoundMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeLatestRoundMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) LatestRoundData(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[LatestRoundDataOutput] {
	calldata, err := c.Codec.EncodeLatestRoundDataMethodCall()
	if err != nil {
		return cre.PromiseFromResult[LatestRoundDataOutput](LatestRoundDataOutput{}, err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (LatestRoundDataOutput, error) {
		return c.Codec.DecodeLatestRoundDataMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) LatestTimestamp(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeLatestTimestampMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeLatestTimestampMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) Owner(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeOwnerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeOwnerMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) TypeAndVersion(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeTypeAndVersionMethodCall()
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeTypeAndVersionMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) Version(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeVersionMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeVersionMethodOutput(response.Data)
	})

}

func (c DataFeedsCache) WriteReportFromWorkflowMetadata(
	runtime cre.Runtime,
	input WorkflowMetadata,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeWorkflowMetadataStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c DataFeedsCache) WriteReport(
	runtime cre.Runtime,
	report *cre.Report,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  c.Address.Bytes(),
		Report:    report,
		GasConfig: gasConfig,
	})
}

// DecodeAddressEmptyCodeError decodes a AddressEmptyCode error from revert data.
func (c *DataFeedsCache) DecodeAddressEmptyCodeError(data []byte) (*AddressEmptyCode, error) {
	args := c.ABI.Errors["AddressEmptyCode"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	target, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for target in AddressEmptyCode error")
	}

	return &AddressEmptyCode{
		Target: target,
	}, nil
}

// Error implements the error interface for AddressEmptyCode.
func (e *AddressEmptyCode) Error() string {
	return fmt.Sprintf("AddressEmptyCode error: target=%v;", e.Target)
}

// DecodeAddressInsufficientBalanceError decodes a AddressInsufficientBalance error from revert data.
func (c *DataFeedsCache) DecodeAddressInsufficientBalanceError(data []byte) (*AddressInsufficientBalance, error) {
	args := c.ABI.Errors["AddressInsufficientBalance"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	account, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for account in AddressInsufficientBalance error")
	}

	return &AddressInsufficientBalance{
		Account: account,
	}, nil
}

// Error implements the error interface for AddressInsufficientBalance.
func (e *AddressInsufficientBalance) Error() string {
	return fmt.Sprintf("AddressInsufficientBalance error: account=%v;", e.Account)
}

// DecodeArrayLengthMismatchError decodes a ArrayLengthMismatch error from revert data.
func (c *DataFeedsCache) DecodeArrayLengthMismatchError(data []byte) (*ArrayLengthMismatch, error) {
	args := c.ABI.Errors["ArrayLengthMismatch"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &ArrayLengthMismatch{}, nil
}

// Error implements the error interface for ArrayLengthMismatch.
func (e *ArrayLengthMismatch) Error() string {
	return fmt.Sprintf("ArrayLengthMismatch error:")
}

// DecodeEmptyConfigError decodes a EmptyConfig error from revert data.
func (c *DataFeedsCache) DecodeEmptyConfigError(data []byte) (*EmptyConfig, error) {
	args := c.ABI.Errors["EmptyConfig"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &EmptyConfig{}, nil
}

// Error implements the error interface for EmptyConfig.
func (e *EmptyConfig) Error() string {
	return fmt.Sprintf("EmptyConfig error:")
}

// DecodeErrorSendingNativeError decodes a ErrorSendingNative error from revert data.
func (c *DataFeedsCache) DecodeErrorSendingNativeError(data []byte) (*ErrorSendingNative, error) {
	args := c.ABI.Errors["ErrorSendingNative"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 3 {
		return nil, fmt.Errorf("expected 3 values, got %d", len(values))
	}

	to, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for to in ErrorSendingNative error")
	}

	amount, ok1 := values[1].(*big.Int)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for amount in ErrorSendingNative error")
	}

	data, ok2 := values[2].([]byte)
	if !ok2 {
		return nil, fmt.Errorf("unexpected type for data in ErrorSendingNative error")
	}

	return &ErrorSendingNative{
		To:     to,
		Amount: amount,
		Data:   data,
	}, nil
}

// Error implements the error interface for ErrorSendingNative.
func (e *ErrorSendingNative) Error() string {
	return fmt.Sprintf("ErrorSendingNative error: to=%v; amount=%v; data=%v;", e.To, e.Amount, e.Data)
}

// DecodeFailedInnerCallError decodes a FailedInnerCall error from revert data.
func (c *DataFeedsCache) DecodeFailedInnerCallError(data []byte) (*FailedInnerCall, error) {
	args := c.ABI.Errors["FailedInnerCall"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &FailedInnerCall{}, nil
}

// Error implements the error interface for FailedInnerCall.
func (e *FailedInnerCall) Error() string {
	return fmt.Sprintf("FailedInnerCall error:")
}

// DecodeFeedNotConfiguredError decodes a FeedNotConfigured error from revert data.
func (c *DataFeedsCache) DecodeFeedNotConfiguredError(data []byte) (*FeedNotConfigured, error) {
	args := c.ABI.Errors["FeedNotConfigured"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	dataId, ok0 := values[0].([16]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for dataId in FeedNotConfigured error")
	}

	return &FeedNotConfigured{
		DataId: dataId,
	}, nil
}

// Error implements the error interface for FeedNotConfigured.
func (e *FeedNotConfigured) Error() string {
	return fmt.Sprintf("FeedNotConfigured error: dataId=%v;", e.DataId)
}

// DecodeInsufficientBalanceError decodes a InsufficientBalance error from revert data.
func (c *DataFeedsCache) DecodeInsufficientBalanceError(data []byte) (*InsufficientBalance, error) {
	args := c.ABI.Errors["InsufficientBalance"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	balance, ok0 := values[0].(*big.Int)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for balance in InsufficientBalance error")
	}

	requiredBalance, ok1 := values[1].(*big.Int)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for requiredBalance in InsufficientBalance error")
	}

	return &InsufficientBalance{
		Balance:         balance,
		RequiredBalance: requiredBalance,
	}, nil
}

// Error implements the error interface for InsufficientBalance.
func (e *InsufficientBalance) Error() string {
	return fmt.Sprintf("InsufficientBalance error: balance=%v; requiredBalance=%v;", e.Balance, e.RequiredBalance)
}

// DecodeInvalidAddressError decodes a InvalidAddress error from revert data.
func (c *DataFeedsCache) DecodeInvalidAddressError(data []byte) (*InvalidAddress, error) {
	args := c.ABI.Errors["InvalidAddress"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	addr, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for addr in InvalidAddress error")
	}

	return &InvalidAddress{
		Addr: addr,
	}, nil
}

// Error implements the error interface for InvalidAddress.
func (e *InvalidAddress) Error() string {
	return fmt.Sprintf("InvalidAddress error: addr=%v;", e.Addr)
}

// DecodeInvalidDataIdError decodes a InvalidDataId error from revert data.
func (c *DataFeedsCache) DecodeInvalidDataIdError(data []byte) (*InvalidDataId, error) {
	args := c.ABI.Errors["InvalidDataId"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidDataId{}, nil
}

// Error implements the error interface for InvalidDataId.
func (e *InvalidDataId) Error() string {
	return fmt.Sprintf("InvalidDataId error:")
}

// DecodeInvalidWorkflowNameError decodes a InvalidWorkflowName error from revert data.
func (c *DataFeedsCache) DecodeInvalidWorkflowNameError(data []byte) (*InvalidWorkflowName, error) {
	args := c.ABI.Errors["InvalidWorkflowName"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	workflowName, ok0 := values[0].([10]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for workflowName in InvalidWorkflowName error")
	}

	return &InvalidWorkflowName{
		WorkflowName: workflowName,
	}, nil
}

// Error implements the error interface for InvalidWorkflowName.
func (e *InvalidWorkflowName) Error() string {
	return fmt.Sprintf("InvalidWorkflowName error: workflowName=%v;", e.WorkflowName)
}

// DecodeNoMappingForSenderError decodes a NoMappingForSender error from revert data.
func (c *DataFeedsCache) DecodeNoMappingForSenderError(data []byte) (*NoMappingForSender, error) {
	args := c.ABI.Errors["NoMappingForSender"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	proxy, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for proxy in NoMappingForSender error")
	}

	return &NoMappingForSender{
		Proxy: proxy,
	}, nil
}

// Error implements the error interface for NoMappingForSender.
func (e *NoMappingForSender) Error() string {
	return fmt.Sprintf("NoMappingForSender error: proxy=%v;", e.Proxy)
}

// DecodeSafeERC20FailedOperationError decodes a SafeERC20FailedOperation error from revert data.
func (c *DataFeedsCache) DecodeSafeERC20FailedOperationError(data []byte) (*SafeERC20FailedOperation, error) {
	args := c.ABI.Errors["SafeERC20FailedOperation"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	token, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for token in SafeERC20FailedOperation error")
	}

	return &SafeERC20FailedOperation{
		Token: token,
	}, nil
}

// Error implements the error interface for SafeERC20FailedOperation.
func (e *SafeERC20FailedOperation) Error() string {
	return fmt.Sprintf("SafeERC20FailedOperation error: token=%v;", e.Token)
}

// DecodeUnauthorizedCallerError decodes a UnauthorizedCaller error from revert data.
func (c *DataFeedsCache) DecodeUnauthorizedCallerError(data []byte) (*UnauthorizedCaller, error) {
	args := c.ABI.Errors["UnauthorizedCaller"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	caller, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for caller in UnauthorizedCaller error")
	}

	return &UnauthorizedCaller{
		Caller: caller,
	}, nil
}

// Error implements the error interface for UnauthorizedCaller.
func (e *UnauthorizedCaller) Error() string {
	return fmt.Sprintf("UnauthorizedCaller error: caller=%v;", e.Caller)
}

func (c *DataFeedsCache) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["AddressEmptyCode"].ID.Bytes()[:4]):
		return c.DecodeAddressEmptyCodeError(data)
	case common.Bytes2Hex(c.ABI.Errors["AddressInsufficientBalance"].ID.Bytes()[:4]):
		return c.DecodeAddressInsufficientBalanceError(data)
	case common.Bytes2Hex(c.ABI.Errors["ArrayLengthMismatch"].ID.Bytes()[:4]):
		return c.DecodeArrayLengthMismatchError(data)
	case common.Bytes2Hex(c.ABI.Errors["EmptyConfig"].ID.Bytes()[:4]):
		return c.DecodeEmptyConfigError(data)
	case common.Bytes2Hex(c.ABI.Errors["ErrorSendingNative"].ID.Bytes()[:4]):
		return c.DecodeErrorSendingNativeError(data)
	case common.Bytes2Hex(c.ABI.Errors["FailedInnerCall"].ID.Bytes()[:4]):
		return c.DecodeFailedInnerCallError(data)
	case common.Bytes2Hex(c.ABI.Errors["FeedNotConfigured"].ID.Bytes()[:4]):
		return c.DecodeFeedNotConfiguredError(data)
	case common.Bytes2Hex(c.ABI.Errors["InsufficientBalance"].ID.Bytes()[:4]):
		return c.DecodeInsufficientBalanceError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidAddress"].ID.Bytes()[:4]):
		return c.DecodeInvalidAddressError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidDataId"].ID.Bytes()[:4]):
		return c.DecodeInvalidDataIdError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidWorkflowName"].ID.Bytes()[:4]):
		return c.DecodeInvalidWorkflowNameError(data)
	case common.Bytes2Hex(c.ABI.Errors["NoMappingForSender"].ID.Bytes()[:4]):
		return c.DecodeNoMappingForSenderError(data)
	case common.Bytes2Hex(c.ABI.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]):
		return c.DecodeSafeERC20FailedOperationError(data)
	case common.Bytes2Hex(c.ABI.Errors["UnauthorizedCaller"].ID.Bytes()[:4]):
		return c.DecodeUnauthorizedCallerError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// AnswerUpdatedTrigger wraps the raw log trigger and provides decoded AnswerUpdatedDecoded data
type AnswerUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into AnswerUpdated data
func (t *AnswerUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[AnswerUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeAnswerUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode AnswerUpdated log: %w", err)
	}

	return &bindings.DecodedLog[AnswerUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerAnswerUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []AnswerUpdated) (cre.Trigger[*evm.Log, *bindings.DecodedLog[AnswerUpdatedDecoded]], error) {
	event := c.ABI.Events["AnswerUpdated"]
	topics, err := c.Codec.EncodeAnswerUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for AnswerUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &AnswerUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsAnswerUpdated(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.AnswerUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// BundleFeedConfigSetTrigger wraps the raw log trigger and provides decoded BundleFeedConfigSetDecoded data
type BundleFeedConfigSetTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into BundleFeedConfigSet data
func (t *BundleFeedConfigSetTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[BundleFeedConfigSetDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeBundleFeedConfigSet(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BundleFeedConfigSet log: %w", err)
	}

	return &bindings.DecodedLog[BundleFeedConfigSetDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerBundleFeedConfigSetLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []BundleFeedConfigSet) (cre.Trigger[*evm.Log, *bindings.DecodedLog[BundleFeedConfigSetDecoded]], error) {
	event := c.ABI.Events["BundleFeedConfigSet"]
	topics, err := c.Codec.EncodeBundleFeedConfigSetTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for BundleFeedConfigSet: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &BundleFeedConfigSetTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsBundleFeedConfigSet(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.BundleFeedConfigSetLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// BundleReportUpdatedTrigger wraps the raw log trigger and provides decoded BundleReportUpdatedDecoded data
type BundleReportUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into BundleReportUpdated data
func (t *BundleReportUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[BundleReportUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeBundleReportUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BundleReportUpdated log: %w", err)
	}

	return &bindings.DecodedLog[BundleReportUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerBundleReportUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []BundleReportUpdated) (cre.Trigger[*evm.Log, *bindings.DecodedLog[BundleReportUpdatedDecoded]], error) {
	event := c.ABI.Events["BundleReportUpdated"]
	topics, err := c.Codec.EncodeBundleReportUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for BundleReportUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &BundleReportUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsBundleReportUpdated(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.BundleReportUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// DecimalFeedConfigSetTrigger wraps the raw log trigger and provides decoded DecimalFeedConfigSetDecoded data
type DecimalFeedConfigSetTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into DecimalFeedConfigSet data
func (t *DecimalFeedConfigSetTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[DecimalFeedConfigSetDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeDecimalFeedConfigSet(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DecimalFeedConfigSet log: %w", err)
	}

	return &bindings.DecodedLog[DecimalFeedConfigSetDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerDecimalFeedConfigSetLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []DecimalFeedConfigSet) (cre.Trigger[*evm.Log, *bindings.DecodedLog[DecimalFeedConfigSetDecoded]], error) {
	event := c.ABI.Events["DecimalFeedConfigSet"]
	topics, err := c.Codec.EncodeDecimalFeedConfigSetTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for DecimalFeedConfigSet: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &DecimalFeedConfigSetTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsDecimalFeedConfigSet(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.DecimalFeedConfigSetLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// DecimalReportUpdatedTrigger wraps the raw log trigger and provides decoded DecimalReportUpdatedDecoded data
type DecimalReportUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into DecimalReportUpdated data
func (t *DecimalReportUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[DecimalReportUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeDecimalReportUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DecimalReportUpdated log: %w", err)
	}

	return &bindings.DecodedLog[DecimalReportUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerDecimalReportUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []DecimalReportUpdated) (cre.Trigger[*evm.Log, *bindings.DecodedLog[DecimalReportUpdatedDecoded]], error) {
	event := c.ABI.Events["DecimalReportUpdated"]
	topics, err := c.Codec.EncodeDecimalReportUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for DecimalReportUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &DecimalReportUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsDecimalReportUpdated(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.DecimalReportUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// FeedAdminSetTrigger wraps the raw log trigger and provides decoded FeedAdminSetDecoded data
type FeedAdminSetTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into FeedAdminSet data
func (t *FeedAdminSetTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[FeedAdminSetDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeFeedAdminSet(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode FeedAdminSet log: %w", err)
	}

	return &bindings.DecodedLog[FeedAdminSetDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerFeedAdminSetLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []FeedAdminSet) (cre.Trigger[*evm.Log, *bindings.DecodedLog[FeedAdminSetDecoded]], error) {
	event := c.ABI.Events["FeedAdminSet"]
	topics, err := c.Codec.EncodeFeedAdminSetTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for FeedAdminSet: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &FeedAdminSetTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsFeedAdminSet(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.FeedAdminSetLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// FeedConfigRemovedTrigger wraps the raw log trigger and provides decoded FeedConfigRemovedDecoded data
type FeedConfigRemovedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into FeedConfigRemoved data
func (t *FeedConfigRemovedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[FeedConfigRemovedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeFeedConfigRemoved(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode FeedConfigRemoved log: %w", err)
	}

	return &bindings.DecodedLog[FeedConfigRemovedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerFeedConfigRemovedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []FeedConfigRemoved) (cre.Trigger[*evm.Log, *bindings.DecodedLog[FeedConfigRemovedDecoded]], error) {
	event := c.ABI.Events["FeedConfigRemoved"]
	topics, err := c.Codec.EncodeFeedConfigRemovedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for FeedConfigRemoved: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &FeedConfigRemovedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsFeedConfigRemoved(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.FeedConfigRemovedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// InvalidUpdatePermissionTrigger wraps the raw log trigger and provides decoded InvalidUpdatePermissionDecoded data
type InvalidUpdatePermissionTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into InvalidUpdatePermission data
func (t *InvalidUpdatePermissionTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[InvalidUpdatePermissionDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeInvalidUpdatePermission(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode InvalidUpdatePermission log: %w", err)
	}

	return &bindings.DecodedLog[InvalidUpdatePermissionDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerInvalidUpdatePermissionLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []InvalidUpdatePermission) (cre.Trigger[*evm.Log, *bindings.DecodedLog[InvalidUpdatePermissionDecoded]], error) {
	event := c.ABI.Events["InvalidUpdatePermission"]
	topics, err := c.Codec.EncodeInvalidUpdatePermissionTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for InvalidUpdatePermission: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &InvalidUpdatePermissionTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsInvalidUpdatePermission(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.InvalidUpdatePermissionLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// NewRoundTrigger wraps the raw log trigger and provides decoded NewRoundDecoded data
type NewRoundTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into NewRound data
func (t *NewRoundTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[NewRoundDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeNewRound(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode NewRound log: %w", err)
	}

	return &bindings.DecodedLog[NewRoundDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerNewRoundLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []NewRound) (cre.Trigger[*evm.Log, *bindings.DecodedLog[NewRoundDecoded]], error) {
	event := c.ABI.Events["NewRound"]
	topics, err := c.Codec.EncodeNewRoundTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for NewRound: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &NewRoundTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsNewRound(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.NewRoundLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// OwnershipTransferRequestedTrigger wraps the raw log trigger and provides decoded OwnershipTransferRequestedDecoded data
type OwnershipTransferRequestedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferRequested data
func (t *OwnershipTransferRequestedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferRequestedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferRequested(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferRequested log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferRequestedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerOwnershipTransferRequestedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferRequested) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferRequestedDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferRequested"]
	topics, err := c.Codec.EncodeOwnershipTransferRequestedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferRequested: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferRequestedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsOwnershipTransferRequested(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferRequestedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// OwnershipTransferredTrigger wraps the raw log trigger and provides decoded OwnershipTransferredDecoded data
type OwnershipTransferredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferred data
func (t *OwnershipTransferredTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferredDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferred(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferred log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferredDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferred) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferred"]
	topics, err := c.Codec.EncodeOwnershipTransferredTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferred: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferredTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferredLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// ProxyDataIdRemovedTrigger wraps the raw log trigger and provides decoded ProxyDataIdRemovedDecoded data
type ProxyDataIdRemovedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into ProxyDataIdRemoved data
func (t *ProxyDataIdRemovedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ProxyDataIdRemovedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeProxyDataIdRemoved(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ProxyDataIdRemoved log: %w", err)
	}

	return &bindings.DecodedLog[ProxyDataIdRemovedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerProxyDataIdRemovedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ProxyDataIdRemoved) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ProxyDataIdRemovedDecoded]], error) {
	event := c.ABI.Events["ProxyDataIdRemoved"]
	topics, err := c.Codec.EncodeProxyDataIdRemovedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ProxyDataIdRemoved: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ProxyDataIdRemovedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsProxyDataIdRemoved(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ProxyDataIdRemovedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// ProxyDataIdUpdatedTrigger wraps the raw log trigger and provides decoded ProxyDataIdUpdatedDecoded data
type ProxyDataIdUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into ProxyDataIdUpdated data
func (t *ProxyDataIdUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ProxyDataIdUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeProxyDataIdUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ProxyDataIdUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ProxyDataIdUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerProxyDataIdUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ProxyDataIdUpdated) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ProxyDataIdUpdatedDecoded]], error) {
	event := c.ABI.Events["ProxyDataIdUpdated"]
	topics, err := c.Codec.EncodeProxyDataIdUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ProxyDataIdUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ProxyDataIdUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsProxyDataIdUpdated(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ProxyDataIdUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// StaleBundleReportTrigger wraps the raw log trigger and provides decoded StaleBundleReportDecoded data
type StaleBundleReportTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into StaleBundleReport data
func (t *StaleBundleReportTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[StaleBundleReportDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeStaleBundleReport(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode StaleBundleReport log: %w", err)
	}

	return &bindings.DecodedLog[StaleBundleReportDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerStaleBundleReportLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []StaleBundleReport) (cre.Trigger[*evm.Log, *bindings.DecodedLog[StaleBundleReportDecoded]], error) {
	event := c.ABI.Events["StaleBundleReport"]
	topics, err := c.Codec.EncodeStaleBundleReportTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for StaleBundleReport: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &StaleBundleReportTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsStaleBundleReport(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.StaleBundleReportLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// StaleDecimalReportTrigger wraps the raw log trigger and provides decoded StaleDecimalReportDecoded data
type StaleDecimalReportTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into StaleDecimalReport data
func (t *StaleDecimalReportTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[StaleDecimalReportDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeStaleDecimalReport(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode StaleDecimalReport log: %w", err)
	}

	return &bindings.DecodedLog[StaleDecimalReportDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerStaleDecimalReportLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []StaleDecimalReport) (cre.Trigger[*evm.Log, *bindings.DecodedLog[StaleDecimalReportDecoded]], error) {
	event := c.ABI.Events["StaleDecimalReport"]
	topics, err := c.Codec.EncodeStaleDecimalReportTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for StaleDecimalReport: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &StaleDecimalReportTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsStaleDecimalReport(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.StaleDecimalReportLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}

// TokenRecoveredTrigger wraps the raw log trigger and provides decoded TokenRecoveredDecoded data
type TokenRecoveredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                 // Embed the raw trigger
	contract                        *DataFeedsCache // Keep reference for decoding
}

// Adapt method that decodes the log into TokenRecovered data
func (t *TokenRecoveredTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[TokenRecoveredDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeTokenRecovered(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode TokenRecovered log: %w", err)
	}

	return &bindings.DecodedLog[TokenRecoveredDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *DataFeedsCache) LogTriggerTokenRecoveredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []TokenRecovered) (cre.Trigger[*evm.Log, *bindings.DecodedLog[TokenRecoveredDecoded]], error) {
	event := c.ABI.Events["TokenRecovered"]
	topics, err := c.Codec.EncodeTokenRecoveredTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for TokenRecovered: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &TokenRecoveredTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *DataFeedsCache) FilterLogsTokenRecovered(runtime cre.Runtime, options *bindings.FilterOptions) cre.Promise[*evm.FilterLogsReply] {
	if options == nil {
		options = &bindings.FilterOptions{
			ToBlock: options.ToBlock,
		}
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.TokenRecoveredLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	})
}
