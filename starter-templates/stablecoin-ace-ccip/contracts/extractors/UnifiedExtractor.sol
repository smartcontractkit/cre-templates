// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.26;

import {IExtractor} from "@chainlink/policy-management/interfaces/IExtractor.sol";
import {IPolicyEngine} from "@chainlink/policy-management/interfaces/IPolicyEngine.sol";

/**
 * @title UnifiedExtractor
 * @notice Unified extractor for both minting and CCIP operations via onReport(bytes,bytes)
 * @dev This extractor handles TWO different report formats based on instructionType:
 * 
 * MINT Report (instructionType = 1):
 *   (uint8 instructionType, address beneficiary, uint256 amount, bytes32 bankRef)
 * 
 * CCIP Report (instructionType = 2):
 *   (uint64 destChain, address sender, address beneficiary, uint256 amount, bytes32 bankRef)
 * 
 * The extractor detects the format and extracts all relevant parameters with common names.
 * Each consumer's policies will filter and use only the parameters they need.
 * 
 * Common Parameters (extracted for both):
 *   - beneficiary: Who receives tokens
 *   - amount: How much to mint/transfer
 *   - bankRef: Off-chain reference
 * 
 * CCIP-only Parameters:
 *   - destinationChain: Target chain selector
 *   - sender: Who initiated the transfer
 * 
 * Mint-only Parameters:
 *   - instructionType: 1=MINT, 2=BURN (legacy from Phase 1)
 * 
 * Usage:
 *   policyEngine.setExtractor(onReport.selector, unifiedExtractor);
 *   
 *   // Then configure policies per consumer:
 *   policyEngine.addPolicy(mintingConsumer, selector, blacklistPolicy, ["beneficiary"]);
 *   policyEngine.addPolicy(ccipConsumer, selector, volumePolicy, ["amount"]);
 */
contract UnifiedExtractor is IExtractor {
    /// @notice Parameter keys (ACE convention: keccak256)
    bytes32 public constant PARAM_INSTRUCTION_TYPE = keccak256("instructionType");
    bytes32 public constant PARAM_DESTINATION_CHAIN = keccak256("destinationChain");
    bytes32 public constant PARAM_SENDER = keccak256("sender");
    bytes32 public constant PARAM_BENEFICIARY = keccak256("beneficiary");
    bytes32 public constant PARAM_AMOUNT = keccak256("amount");
    bytes32 public constant PARAM_BANK_REF = keccak256("bankRef");
    
    /**
     * @notice Instruction types
     */
    uint8 constant INSTRUCTION_MINT = 1;
    uint8 constant INSTRUCTION_BURN = 2;
    
    /**
     * @notice Extracts parameters from onReport(bytes metadata, bytes report) calls.
     * @dev Handles both minting and CCIP report formats.
     * 
     * @param payload The transaction payload from PolicyEngine.
     * @return parameters Array of extracted parameters with standardized names.
     */
    function extract(IPolicyEngine.Payload memory payload) 
        external 
        pure 
        returns (IPolicyEngine.Parameter[] memory) 
    {
        // Decode onReport parameters: (bytes metadata, bytes report)
        (, bytes memory report) = abi.decode(payload.data, (bytes, bytes));
        
        // Determine report type by decoding the first field
        // Mint reports: (uint8 instructionType, ...) - instructionType is 1 or 2
        // CCIP reports: (uint64 destinationChainSelector, ...) - chain selector > 255
        
        // Decode first field as uint256 (covers both uint8 and uint64)
        uint256 firstField = abi.decode(report, (uint256));
        
        // If firstField > 255, it's a CCIP report (chain selector)
        // If firstField <= 255, it's a Mint report (instruction type)
        if (firstField > 255) {
            return extractCCIPReport(report);
        } else {
            return extractMintReport(report);
        }
    }
    
    /**
     * @notice Extracts parameters from a minting report.
     * @dev Format: (uint8 instructionType, address beneficiary, uint256 amount, bytes32 bankRef)
     */
    function extractMintReport(bytes memory report) 
        internal 
        pure 
        returns (IPolicyEngine.Parameter[] memory) 
    {
        (uint8 instructionType, address beneficiary, uint256 amount, bytes32 bankRef) 
            = abi.decode(report, (uint8, address, uint256, bytes32));
        
        // Return 4 parameters for minting operations
        IPolicyEngine.Parameter[] memory parameters = new IPolicyEngine.Parameter[](4);
        
        parameters[0] = IPolicyEngine.Parameter({
            name: PARAM_INSTRUCTION_TYPE,
            value: abi.encode(instructionType)
        });
        
        parameters[1] = IPolicyEngine.Parameter({
            name: PARAM_BENEFICIARY,
            value: abi.encode(beneficiary)
        });
        
        parameters[2] = IPolicyEngine.Parameter({
            name: PARAM_AMOUNT,
            value: abi.encode(amount)
        });
        
        parameters[3] = IPolicyEngine.Parameter({
            name: PARAM_BANK_REF,
            value: abi.encode(bankRef)
        });
        
        return parameters;
    }
    
    /**
     * @notice Extracts parameters from a CCIP report.
     * @dev Format: (uint64 destChain, address sender, address beneficiary, uint256 amount, bytes32 bankRef)
     */
    function extractCCIPReport(bytes memory report) 
        internal 
        pure 
        returns (IPolicyEngine.Parameter[] memory) 
    {
        (uint64 destinationChainSelector, address sender, address beneficiary, uint256 amount, bytes32 bankRef) 
            = abi.decode(report, (uint64, address, address, uint256, bytes32));
        
        // Return 5 parameters for CCIP operations
        IPolicyEngine.Parameter[] memory parameters = new IPolicyEngine.Parameter[](5);
        
        parameters[0] = IPolicyEngine.Parameter({
            name: PARAM_DESTINATION_CHAIN,
            value: abi.encode(destinationChainSelector)
        });
        
        parameters[1] = IPolicyEngine.Parameter({
            name: PARAM_SENDER,
            value: abi.encode(sender)
        });
        
        parameters[2] = IPolicyEngine.Parameter({
            name: PARAM_BENEFICIARY,
            value: abi.encode(beneficiary)
        });
        
        parameters[3] = IPolicyEngine.Parameter({
            name: PARAM_AMOUNT,
            value: abi.encode(amount)
        });
        
        parameters[4] = IPolicyEngine.Parameter({
            name: PARAM_BANK_REF,
            value: abi.encode(bankRef)
        });
        
        return parameters;
    }
}

