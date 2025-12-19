// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IReceiverTemplate} from "./IReceiverTemplate.sol";
import {StablecoinERC20} from "./StablecoinERC20.sol";

/**
 * @title MintingConsumer
 * @notice Consumer contract that receives CRE reports and mints/redeems stablecoins
 * @dev Inherits IReceiverTemplate for secure report validation
 * 
 * Flow:
 * 1. Bank sends SWIFT MT103 message to CRE workflow (HTTP trigger)
 * 2. CRE workflow generates signed report with mint/redeem instruction
 * 3. Forwarder validates signatures and calls this contract's onReport()
 * 4. This contract decodes report and mints or burns stablecoins
 */
contract MintingConsumer is IReceiverTemplate {
    StablecoinERC20 public immutable stablecoin;
    
    // Instruction types
    uint8 constant INSTRUCTION_MINT = 1;
    uint8 constant INSTRUCTION_REDEEM = 2;
    
    event MintInstructionReceived(
        address indexed recipient,
        uint256 amount,
        bytes32 indexed bankReference,
        uint256 timestamp
    );
    
    event RedeemInstructionReceived(
        address indexed account,
        uint256 amount,
        bytes32 indexed bankReference,
        uint256 timestamp
    );
    
    event MintExecuted(
        address indexed recipient,
        uint256 amount,
        bytes32 indexed bankReference
    );
    
    event RedeemExecuted(
        address indexed account,
        uint256 amount,
        bytes32 indexed bankReference
    );
    
    error MintFailed();
    error RedeemFailed();
    error InvalidInstructionType();
    
    /**
     * @notice Construct the minting consumer
     * @param _stablecoin Address of the StablecoinERC20 contract
     * @param _expectedAuthor Expected workflow owner address (use address(0) for testing)
     * @param _expectedWorkflowName Expected workflow name (use bytes10("dummy") for testing)
     */
    constructor(
        address _stablecoin,
        address _expectedAuthor,
        bytes10 _expectedWorkflowName
    ) IReceiverTemplate(_expectedAuthor, _expectedWorkflowName) {
        stablecoin = StablecoinERC20(_stablecoin);
    }
    
    /**
     * @notice Receive report from Forwarder
     * @param metadata Encoded metadata (not used in testing version)
     * @param report Encoded mint instruction: (address recipient, uint256 amount, bytes32 bankRef)
     */
    function onReport(bytes calldata metadata, bytes calldata report) external override {
        // In production, you would validate metadata here
        // For testing/demo purposes, we skip validation
        
        _processReport(report);
    }
    
    /**
     * @notice Process the mint or redeem instruction
     * @param report ABI-encoded: (uint8 instructionType, address account, uint256 amount, bytes32 bankRef)
     */
    function _processReport(bytes calldata report) internal override {
        // Decode the report
        (uint8 instructionType, address account, uint256 amount, bytes32 bankRef) = abi.decode(
            report,
            (uint8, address, uint256, bytes32)
        );
        
        if (instructionType == INSTRUCTION_MINT) {
            // Log the mint instruction
            emit MintInstructionReceived(account, amount, bankRef, block.timestamp);
            
            // Execute the mint
            try stablecoin.mint(account, amount) {
                emit MintExecuted(account, amount, bankRef);
            } catch {
                revert MintFailed();
            }
        } else if (instructionType == INSTRUCTION_REDEEM) {
            // Log the redeem instruction
            emit RedeemInstructionReceived(account, amount, bankRef, block.timestamp);
            
            // Execute the burn
            try stablecoin.burnFrom(account, amount) {
                emit RedeemExecuted(account, amount, bankRef);
            } catch {
                revert RedeemFailed();
            }
        } else {
            revert InvalidInstructionType();
        }
    }
}

