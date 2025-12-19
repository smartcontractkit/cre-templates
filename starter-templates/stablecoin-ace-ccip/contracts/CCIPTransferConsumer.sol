// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiverTemplate} from "./IReceiverTemplate.sol";
import {IRouterClient} from "@chainlink/contracts-ccip/ccip/interfaces/IRouterClient.sol";
import {Client} from "@chainlink/contracts-ccip/ccip/libraries/Client.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title CCIPTransferConsumer
 * @notice Consumer contract that receives CRE reports and executes cross-chain CCIP transfers
 * @dev Inherits IReceiverTemplate for secure report validation
 * 
 * Flow:
 * 1. Bank sends transfer instruction to CRE workflow (HTTP trigger)
 * 2. CRE workflow generates signed report with transfer data
 * 3. Forwarder validates signatures and calls this contract's onReport()
 * 4. This contract:
 *    - Pulls tokens from sender
 *    - Approves Router to spend tokens (Router transfers to Pool)
 *    - Approves LINK to Router
 *    - Calls Router.ccipSend()
 *    - Emits transfer events
 */
contract CCIPTransferConsumer is IReceiverTemplate {
    IRouterClient public immutable router;
    IERC20 public immutable linkToken;
    IERC20 public immutable stablecoin;
    
    event TransferInitiated(
        bytes32 indexed messageId,
        uint64 indexed destinationChainSelector,
        address indexed recipient,
        uint256 amount,
        bytes32 bankReference,
        uint256 fees
    );
    
    event TransferInstructionReceived(
        uint64 indexed destinationChainSelector,
        address indexed recipient,
        uint256 amount,
        bytes32 indexed bankReference,
        uint256 timestamp
    );
    
    error TransferFailed();
    error InsufficientLinkBalance(uint256 required, uint256 available);
    
    /**
     * @notice Construct the CCIP transfer consumer
     * @param _stablecoin Address of the StablecoinERC20 contract
     * @param _router Address of the CCIP Router contract
     * @param _linkToken Address of the LINK token contract
     * @param _expectedAuthor Expected workflow owner address (use address(0) for testing)
     * @param _expectedWorkflowName Expected workflow name (use bytes10("dummy") for testing)
     */
    constructor(
        address _stablecoin,
        address _router,
        address _linkToken,
        address _expectedAuthor,
        bytes10 _expectedWorkflowName
    ) IReceiverTemplate(_expectedAuthor, _expectedWorkflowName) {
        stablecoin = IERC20(_stablecoin);
        router = IRouterClient(_router);
        linkToken = IERC20(_linkToken);
    }
    
    /**
     * @notice Receive report from Forwarder
     * @param metadata Encoded metadata (not used in testing version)
     * @param report Encoded transfer instruction: (uint64 destChainSelector, address recipient, uint256 amount, bytes32 bankRef)
     */
    function onReport(bytes calldata metadata, bytes calldata report) external override {
        // In production, validate metadata here
        _processReport(report);
    }
    
    /**
     * @notice Process the cross-chain transfer instruction
     * @param report ABI-encoded: (uint64 destChainSelector, address sender, address recipient, uint256 amount, bytes32 bankRef)
     */
    function _processReport(bytes calldata report) internal override {
        // 1. Decode the transfer instruction
        (uint64 destinationChainSelector, address sender, address recipient, uint256 amount, bytes32 bankRef) = abi.decode(
            report,
            (uint64, address, address, uint256, bytes32)
        );
        
        // Validate decoded values
        require(sender != address(0), "Invalid sender");
        require(recipient != address(0), "Invalid recipient");
        require(amount > 0, "Invalid amount");
        
        // 2. Log the instruction
        emit TransferInstructionReceived(
            destinationChainSelector,
            recipient,
            amount,
            bankRef,
            block.timestamp
        );
        
        // 3. Build CCIP message
        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(recipient),
            data: "",
            tokenAmounts: new Client.EVMTokenAmount[](1),
            feeToken: address(linkToken),
            extraArgs: Client._argsToBytes(
                Client.EVMExtraArgsV1({gasLimit: 0}) // 0 for default
            )
        });
        
        message.tokenAmounts[0] = Client.EVMTokenAmount({
            token: address(stablecoin),
            amount: amount
        });
        
        // 4. Calculate fee
        uint256 fees = router.getFee(destinationChainSelector, message);
        
        // 5. Check LINK balance
        uint256 linkBalance = linkToken.balanceOf(address(this));
        if (linkBalance < fees) {
            revert InsufficientLinkBalance(fees, linkBalance);
        }
        
        // 6. Transfer tokens from sender to this contract (sender must have approved consumer)
        require(
            stablecoin.transferFrom(sender, address(this), amount),
            "Token transfer failed"
        );
        
        // 7. Approve Router to spend tokens (Router will transfer to Pool, then Pool burns)
        stablecoin.approve(address(router), amount);
        
        // 8. Approve Router to spend LINK for fees
        linkToken.approve(address(router), fees);
        
        // 10. Execute cross-chain transfer
        bytes32 messageId = router.ccipSend(destinationChainSelector, message);
        
        emit TransferInitiated(
            messageId,
            destinationChainSelector,
            recipient,
            amount,
            bankRef,
            fees
        );
    }
    
    /**
     * @notice Fund this contract with LINK for CCIP fees
     * @dev Owner or anyone can send LINK to this contract
     */
    receive() external payable {}
}






















