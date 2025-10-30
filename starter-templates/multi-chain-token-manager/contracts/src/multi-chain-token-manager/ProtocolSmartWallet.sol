/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {MockPool} from "./MockPool.sol";

import {Client} from "@chainlink/contracts-ccip/libraries/Client.sol";
import {IReceiver} from "@chainlink/contracts/src/v0.8/keystone/interfaces/IReceiver.sol";
import {KeystoneFeedDefaultMetadataLib} from "@chainlink/contracts/src/v0.8/keystone/lib/KeystoneFeedDefaultMetadataLib.sol";
import {CCIPReceiver} from "@chainlink/contracts-ccip/applications/CCIPReceiver.sol";
import {LinkTokenInterface} from "@chainlink/contracts/src/v0.8/shared/interfaces/LinkTokenInterface.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/interfaces/IRouterClient.sol";
import {IERC20} from
    "@chainlink/contracts/src/v0.8/vendor/forge-std/src/interfaces/IERC20.sol";
import {SafeERC20} from
    "@openzeppelin/contracts@5.0.2/token/ERC20/utils/SafeERC20.sol";
import {IERC165} from "@openzeppelin/contracts@5.0.2/utils/introspection/IERC165.sol";

import {OwnerIsCreator} from "@chainlink/contracts/src/v0.8/shared/access/OwnerIsCreator.sol";

contract ProtocolSmartWallet is CCIPReceiver, IReceiver, OwnerIsCreator {
    using KeystoneFeedDefaultMetadataLib for bytes;
    using SafeERC20 for IERC20;

    struct Message {
        uint256 amount;
    }

    struct RebalanceParams {
        address asset;
        uint256 amount;
        uint64 destinationChainSelector;
        address destinationProtocolSmartWallet;
    }

    event KeystoneForwarderSet(address indexed keystoneForwarder);
    event KeystoneForwarderRemoved(address indexed keystoneForwarder);

    event WorkflowOwnerSet(address indexed workflowOwner);
    event WorkflowOwnerRemoved(address indexed workflowOwner);

    event SenderForSourceChainSet(uint64 indexed _sourceChainSelector, address indexed _sender);
    event SenderForSourceChainRemoved(uint64 indexed _sourceChainSelector, address indexed _sender);

    event MessageSent(
        bytes32 indexed messageId,
        uint64 indexed destinationChainSelector,
        address receiver,
        Client.EVMTokenAmount tokenAmount,
        uint256 fees
    );

    event MessageReceived(
        bytes32 indexed messageId,
        uint64 indexed sourceChainSelector,
        address sender,
        Client.EVMTokenAmount tokenAmount
    );

    // Include RebalanceParams so `cre generate-bindings evm` picks up this 
    // type in abigen.
    event ReportReceived(address indexed workflowOwner, bytes10 indexed workflowName, RebalanceParams params);

    error MustBeKeystoneForwarder();
    error InvalidKeystoneForwarder();

    error InvalidWorkflowOwner();
    error UnauthorizedWorkflowOwner(address workflowOwner);

    error ZeroAddress(uint256 index);
    error InsufficientTokenAmount();
    error InsufficientFeeTokenAmount();
    error MismatchedTokenAmount();

    error InvalidSourceChain(); // Used when the source chain is 0
    error InvalidSenderAddress(); // Used when the sender address is 0
    error NoSenderOnSourceChain(uint64 sourceChainSelector); // Used when there is no sender for a given source chain
    error WrongSenderForSourceChain(uint64 sourceChainSelector); // Used when the sender contract is not the correct one

    MockPool public pool;
    address link;
    mapping(address keystoneForwarder => bool allowed) public allowedKeystoneForwarders;
    mapping(address workflowOwner => bool allowed) public allowedWorkflowOwners;
    // Mapping to keep track of the sender contract per source chain.
    mapping(uint64 chainSelector => address sender) public allowedCcipSenders;

    modifier validateSourceChain(uint64 _sourceChainSelector) {
        _validateSourceChain(_sourceChainSelector);
        _;
    }
 
    function _validateSourceChain(uint64 _sourceChainSelector) pure internal {
        if (_sourceChainSelector == 0) revert InvalidSourceChain();
    }

    constructor(address[] memory _keystoneForwarders, address[] memory _allowedWorkflowOwners, address _pool, address _router, address _link) CCIPReceiver(_router)  {
        pool = MockPool(_pool);
        link = _link;
        for (uint256 i = 0; i < _keystoneForwarders.length; i++) {
            allowedKeystoneForwarders[_keystoneForwarders[i]] = true;
        }
        for (uint256 i = 0; i < _allowedWorkflowOwners.length; i++) {
            allowedWorkflowOwners[_allowedWorkflowOwners[i]] = true;
        }
    }

    function setSenderForSourceChain(
        uint64 _sourceChainSelector,
        address _sender
    ) external onlyOwner validateSourceChain(_sourceChainSelector) {
        if (_sender == address(0)) revert InvalidSenderAddress();
        allowedCcipSenders[_sourceChainSelector] = _sender;
        emit SenderForSourceChainSet(_sourceChainSelector, _sender);
    }

    function removeSenderForSourceChain(
        uint64 _sourceChainSelector
    ) external onlyOwner validateSourceChain(_sourceChainSelector) {
        if (allowedCcipSenders[_sourceChainSelector] == address(0)) {
            revert NoSenderOnSourceChain(_sourceChainSelector);
        }
        address sender = allowedCcipSenders[_sourceChainSelector];
        delete allowedCcipSenders[_sourceChainSelector];
        emit SenderForSourceChainRemoved(_sourceChainSelector, sender);
    }

    function setKeystoneForwarder(address _keystoneForwarder) external onlyOwner {
        if (_keystoneForwarder == address(0)) revert InvalidKeystoneForwarder();
        allowedKeystoneForwarders[_keystoneForwarder] = true;
        emit KeystoneForwarderSet(_keystoneForwarder);
    }

    function removeKeystoneForwarder(address _keystoneForwarder) external onlyOwner {
        if (allowedKeystoneForwarders[_keystoneForwarder] == false) {
            revert InvalidKeystoneForwarder();
        }
        delete allowedKeystoneForwarders[_keystoneForwarder];
        emit KeystoneForwarderRemoved(_keystoneForwarder);
    }

    function setWorkflowOwner(address _workflowOwner) external onlyOwner {
        if (_workflowOwner == address(0)) revert InvalidWorkflowOwner();
        allowedWorkflowOwners[_workflowOwner] = true;
        emit WorkflowOwnerSet(_workflowOwner);
    }

    function removeWorkflowOwner(address _workflowOwner) external onlyOwner {
        if (allowedWorkflowOwners[_workflowOwner] == false) {
            revert InvalidWorkflowOwner();
        }
        delete allowedWorkflowOwners[_workflowOwner];
        emit WorkflowOwnerRemoved(_workflowOwner);
    }

    function setPool(address _pool) external onlyOwner {
        pool = MockPool(_pool);
    }

    function depositToPool(address asset, uint256 amount) onlyOwner external {
        _depositToPool(asset, amount);
    }

    function _depositToPool(address asset, uint256 amount) internal {
        IERC20(asset).approve(address(pool), amount);
        pool.supply(asset, amount, address(this), 0);
    }

    function withdrawFromPool(address asset, uint256 amount, address to) onlyOwner external {
        pool.withdraw(asset, amount, to);
    }

    function getPoolAddress() external view returns (address) {
        return address(pool);
    }

    // IReceiver interface implementation.
    // Called by CRE KeystoneForwarder contract.
    function onReport(bytes calldata metadata, bytes calldata report) external override {
        if (allowedKeystoneForwarders[msg.sender] != true) {
            revert MustBeKeystoneForwarder();
        }

        (bytes10 workflowName, address workflowOwner,) = metadata._extractMetadataInfo();
        // TODO: allow multiple owners to support simulate
        if (workflowOwner != owner()) {
            revert UnauthorizedWorkflowOwner(workflowOwner);
        }

        RebalanceParams memory rebalanceParams = abi.decode(report, (RebalanceParams));

        withdrawFromPoolAndDepositCrossChain(
           rebalanceParams
        );

        emit ReportReceived(workflowOwner, workflowName, rebalanceParams);
    }

    function withdrawFromPoolAndDepositCrossChain(
        RebalanceParams memory params
    ) public {
        uint256 amountWithdrawn = pool.withdraw(params.asset, params.amount, address(this));
        Message memory message = Message({amount: amountWithdrawn});
        _sendMessage(params.destinationChainSelector, params.destinationProtocolSmartWallet, message, params.asset, amountWithdrawn);
    }

    function _sendMessage(
        uint64 destinationChainSelector,
        address receiver,
        Message memory crossChainMessage,
        address token,
        uint256 amount
    ) internal returns (bytes32 messageId) {

        if (receiver == address(0)) revert ZeroAddress(0);

        Client.EVMTokenAmount[] memory tokenAmounts = new Client.EVMTokenAmount[](1);
        Client.EVMTokenAmount memory tokenAmount = Client.EVMTokenAmount({token: token, amount: amount});
        tokenAmounts[0] = tokenAmount;

        Client.EVM2AnyMessage memory evm2AnyMessage = Client.EVM2AnyMessage({
            receiver: abi.encode(receiver),
            data: abi.encode(crossChainMessage),
            tokenAmounts: tokenAmounts,
            extraArgs: Client._argsToBytes(
                Client.EVMExtraArgsV1({gasLimit: 300_000 })
            ),
            feeToken: link
        });

        if (IERC20(token).balanceOf(address(this)) < amount) {
            revert InsufficientTokenAmount();
        }

        IRouterClient router = IRouterClient(this.getRouter());

        // Approve the Router to spend tokens on contract's behalf.
        IERC20(token).approve(address(router), amount);

        // Get the fee required to send the message.
        uint256 fees = router.getFee(destinationChainSelector, evm2AnyMessage);

        // Reverts if insufficient funds to cover the CCIP fees.
        if (LinkTokenInterface(link).balanceOf(address(this)) < fees) {
            revert InsufficientFeeTokenAmount();
        }

        // Approve the Router to spend the required fees
        LinkTokenInterface(link).approve(this.getRouter(), fees);
        messageId = router.ccipSend(destinationChainSelector, evm2AnyMessage);

        emit MessageSent(messageId, destinationChainSelector, receiver, tokenAmount, fees);

        return messageId;
    }

    // CCIPReceiver interface implementation.
    // Called by CCIP Router contract.
    function _ccipReceive(Client.Any2EVMMessage memory any2EvmMessage) internal override {

        uint64 sourceChainSelector = any2EvmMessage.sourceChainSelector;
        address sender = abi.decode(any2EvmMessage.sender, (address));

        // validate the sender contract
        if (sender != allowedCcipSenders[any2EvmMessage.sourceChainSelector]) {
            revert WrongSenderForSourceChain(any2EvmMessage.sourceChainSelector);
        }

        bytes32 messageId = any2EvmMessage.messageId;
        Client.EVMTokenAmount[] memory tokenAmounts = any2EvmMessage.destTokenAmounts;
        address token = tokenAmounts[0].token;
        uint256 amount = tokenAmounts[0].amount;

        // TODO: consider dropping this check and simply deposit the amount
        // that is received.
        Message memory message = abi.decode(any2EvmMessage.data, (Message));
        if (amount != message.amount) {
            revert MismatchedTokenAmount();
        }

        _depositToPool(token, amount);

        emit MessageReceived(messageId, sourceChainSelector, sender, tokenAmounts[0]);
    }

    function supportsInterface(
        bytes4 interfaceId
    ) public pure virtual override(CCIPReceiver,IERC165) returns (bool) {
        return super.supportsInterface(interfaceId) || interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
    }
}
