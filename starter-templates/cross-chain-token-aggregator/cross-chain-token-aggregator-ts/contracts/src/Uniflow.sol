// SPDX-License-Identifier: MIT

pragma solidity 0.8.34;

import {ReceiverTemplate} from "./interfaces/ReceiverTemplate.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/contracts/interfaces/IRouterClient.sol";
import {Client} from "@chainlink/contracts-ccip/contracts/libraries/Client.sol";

contract Uniflow is ReceiverTemplate {
    using SafeERC20 for IERC20;

    // Errors
    error Uniflow__TokenNotConfigured();
    error Uniflow__MinAmountNotFulfilled();
    error Uniflow__InvalidReceiverConfigured();
    error Uniflow__AcrossBridgeDepositFailed();
    error Uniflow__EmptyReport();
    error Uniflow__CCIPChainNotAllowlisted(uint64 chainSelector);
    error Uniflow__InsufficientLinkTokenForCCIPBridge();

    // Struct
    struct TokenConfig {
        address receiver;
        uint256 minAmountToTrigger;
    }

    // Variables
    mapping(address token => TokenConfig) public s_tokenConfig;

    // Mapping to keep track of allowlisted destination chains.
    mapping(uint64 => bool) public s_allowlistedChains;

    IRouterClient public s_ccipRouter;

    IERC20 public s_linkToken;

    // Events
    event UniflowApproval(address user, bool approval);
    event AcrossBridgeInitiated(address token, address receiver, uint256 amount);
    event CCIPBridgeInitiated(bytes32 messageId, address token, address receiver, uint256 amount);

    // modifiers
    modifier onlyConfiguredToken(address token) {
        _onlyConfiguredToken(token);
        _;
    }

    // constructor
    constructor(address _forwarderAddress, address _ccipRouter, address _linkToken)
        ReceiverTemplate(_forwarderAddress)
    {
        s_ccipRouter = IRouterClient(_ccipRouter);
        s_linkToken = IERC20(_linkToken);
    }

    function setupToken(address token, TokenConfig memory tokenConfig) external onlyOwner {
        s_tokenConfig[token] =
            TokenConfig({receiver: tokenConfig.receiver, minAmountToTrigger: tokenConfig.minAmountToTrigger});
    }

    function allowlistDestanationChainForCCIP(uint64 selector, bool enable) external {
        s_allowlistedChains[selector] = enable;
    }

    function updateReceiver(address token, address receiver) external onlyConfiguredToken(token) onlyOwner {
        if (s_tokenConfig[token].receiver == address(0)) {
            revert Uniflow__TokenNotConfigured();
        }
        s_tokenConfig[token].receiver = receiver;
    }

    function _performAcrossBridgeOp(bytes calldata report) internal {
        (
            address receiver,
            uint256 amount,
            address token,
            address approvalContract,
            address depositContract,
            bytes memory depositData
        ) = abi.decode(report, (address, uint256, address, address, address, bytes));

        _onlyConfiguredToken(token);

        TokenConfig memory tokenConfig = s_tokenConfig[token];

        if (amount < tokenConfig.minAmountToTrigger) {
            revert Uniflow__MinAmountNotFulfilled();
        }

        if (tokenConfig.receiver != receiver) {
            revert Uniflow__InvalidReceiverConfigured();
        }

        // Perform approval and transfer
        IERC20(token).safeTransferFrom(owner(), address(this), amount);
        IERC20(token).forceApprove(approvalContract, amount);
        emit AcrossBridgeInitiated(token, receiver, amount);
        (bool depositSuccess,) = depositContract.call(depositData);
        if (!depositSuccess) {
            revert Uniflow__AcrossBridgeDepositFailed();
        }
    }

    function _performChainlinkCCIPBridgeOp(bytes calldata report) internal {
        (address receiver, address token, uint256 amount, uint64 destinationChainSelector) =
            abi.decode(report, (address, address, uint256, uint64));
        if (!s_allowlistedChains[destinationChainSelector]) {
            revert Uniflow__CCIPChainNotAllowlisted(destinationChainSelector);
        }

        _onlyConfiguredToken(token);

        TokenConfig memory tokenConfig = s_tokenConfig[token];

        if (amount < tokenConfig.minAmountToTrigger) {
            revert Uniflow__MinAmountNotFulfilled();
        }

        if (tokenConfig.receiver != receiver) {
            revert Uniflow__InvalidReceiverConfigured();
        }

        Client.EVM2AnyMessage memory ccipMessage = _buildCCIPMessage(receiver, token, amount);
        uint256 ccipFees = s_ccipRouter.getFee(destinationChainSelector, ccipMessage);

        uint256 requiredLinkBalance = ccipFees;

        if (token == address(s_linkToken)) {
            requiredLinkBalance += amount;
        }

        IERC20(token).safeTransferFrom(owner(), address(this), amount);

        if (s_linkToken.balanceOf(address(this)) < requiredLinkBalance) {
            revert Uniflow__InsufficientLinkTokenForCCIPBridge();
        }

        s_linkToken.forceApprove(address(s_ccipRouter), requiredLinkBalance);
        if (token != address(s_linkToken)) {
            IERC20(token).forceApprove(address(s_ccipRouter), amount);
        }

        bytes32 messageId = s_ccipRouter.ccipSend(destinationChainSelector, ccipMessage);
        emit CCIPBridgeInitiated(messageId, token, receiver, amount);
    }

    function _buildCCIPMessage(address receiver, address token, uint256 amount)
        internal
        view
        returns (Client.EVM2AnyMessage memory)
    {
        Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
        tokenAmounts[0] = Client.EVMTokenAmount({token: token, amount: amount});

        return Client.EVM2AnyMessage({
            receiver: abi.encode(receiver),
            data: "",
            tokenAmounts: tokenAmounts,
            feeToken: address(s_linkToken),
            extraArgs: Client._argsToBytes(Client.GenericExtraArgsV2({gasLimit: 0, allowOutOfOrderExecution: true}))
        });
    }

    function _processReport(bytes calldata report) internal override {
        if (report.length == 0) {
            revert Uniflow__EmptyReport();
        }

        bytes1 op = report[0];
        if (op == 0x01) {
            _performAcrossBridgeOp(report[1:]);
        } else if (op == 0x02) {
            _performChainlinkCCIPBridgeOp(report[1:]);
        }
    }

    function _onlyConfiguredToken(address token) internal view {
        if (s_tokenConfig[token].receiver == address(0)) {
            revert Uniflow__TokenNotConfigured();
        }
    }
}
