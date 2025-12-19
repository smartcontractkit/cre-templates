// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {IPolicyEngine} from "@chainlink/policy-management/interfaces/IPolicyEngine.sol";
import {UnifiedExtractor} from "../contracts/extractors/UnifiedExtractor.sol";

/**
 * @title ConfigureACEWithConstants
 * @notice Configures ACE using proper keccak256 parameter names (ACE convention)
 * @dev This script:
 *      1. Verifies UnifiedExtractor is attached (already set by DeployUnifiedExtractor.s.sol)
 *      2. Attaches BlacklistPolicy to MintingConsumer (beneficiary parameter)
 *      3. Attaches VolumePolicy to CCIPConsumer (amount parameter)
 */
contract ConfigureACEWithConstants is Script {
    function run() external {
        uint256 deployerPK = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPK);
        
        address policyEngine = vm.envAddress("POLICY_ENGINE");
        address blacklistPolicy = vm.envAddress("BLACKLIST_POLICY");
        address volumePolicy = vm.envAddress("VOLUME_POLICY");
        address unifiedExtractor = vm.envAddress("UNIFIED_EXTRACTOR");
        address mintingConsumer = vm.envAddress("MINTING_CONSUMER_ACE");
        address ccipConsumer = vm.envAddress("CCIP_CONSUMER_ACE");
        
        // CRITICAL: Calculate selector dynamically to avoid hardcoding errors
        // onReport(bytes,bytes) = 0x805f2132
        bytes4 onReportSelector = bytes4(keccak256("onReport(bytes,bytes)"));
        
        console.log("=== Configuring ACE with ACE Convention (keccak256) ===");
        console.log("");
        console.log("PolicyEngine:", policyEngine);
        console.log("UnifiedExtractor:", unifiedExtractor);
        console.log("BlacklistPolicy:", blacklistPolicy);
        console.log("VolumePolicy:", volumePolicy);
        console.log("MintingConsumer:", mintingConsumer);
        console.log("CCIPConsumer:", ccipConsumer);
        console.log("");
        
        vm.startBroadcast(deployerPK);
        
        // 1. Verify UnifiedExtractor is attached (already done by DeployUnifiedExtractor.s.sol)
        console.log("1. Verifying UnifiedExtractor...");
        address currentExtractor = IPolicyEngine(policyEngine).getExtractor(onReportSelector);
        require(currentExtractor == unifiedExtractor, "UnifiedExtractor not set correctly");
        console.log("   UnifiedExtractor verified:", currentExtractor);
        
        // 2. Attach BlacklistPolicy to MintingConsumer (checks 'beneficiary' parameter)
        console.log("2. Attaching BlacklistPolicy to MintingConsumer...");
        
        // Use the CONSTANT from the extractor (ACE convention!)
        UnifiedExtractor extractor = UnifiedExtractor(unifiedExtractor);
        bytes32[] memory mintParams = new bytes32[](1);
        mintParams[0] = extractor.PARAM_BENEFICIARY(); // keccak256("beneficiary")
        
        console.log("   Parameter name (keccak256('beneficiary')):");
        console.logBytes32(mintParams[0]);
        
        // Remove existing policy if any
        try IPolicyEngine(policyEngine).removePolicy(mintingConsumer, onReportSelector, blacklistPolicy) {
            console.log("   Removed existing policy");
        } catch {
            console.log("   No existing policy");
        }
        
        // Add policy with proper parameter
        IPolicyEngine(policyEngine).addPolicy(
            mintingConsumer,
            onReportSelector,
            blacklistPolicy,
            mintParams
        );
        console.log("   BlacklistPolicy attached!");
        
        // 3. Attach VolumePolicy to CCIPConsumer (checks 'amount' parameter)
        console.log("3. Attaching VolumePolicy to CCIPConsumer...");
        
        bytes32[] memory ccipParams = new bytes32[](1);
        ccipParams[0] = extractor.PARAM_AMOUNT(); // keccak256("amount")
        
        console.log("   Parameter name (keccak256('amount')):");
        console.logBytes32(ccipParams[0]);
        
        // Remove existing policy if any
        try IPolicyEngine(policyEngine).removePolicy(ccipConsumer, onReportSelector, volumePolicy) {
            console.log("   Removed existing policy");
        } catch {
            console.log("   No existing policy");
        }
        
        IPolicyEngine(policyEngine).addPolicy(
            ccipConsumer,
            onReportSelector,
            volumePolicy,
            ccipParams
        );
        console.log("   VolumePolicy attached!");
        
        vm.stopBroadcast();
        
        console.log("");
        console.log("=== Configuration Complete (ACE Convention) ===");
        console.log("");
        console.log("Verification:");
        console.log("1. Extractor attached:", IPolicyEngine(policyEngine).getExtractor(onReportSelector));
        console.log("2. Policies on MintingConsumer: Use getPolicies() to verify");
        console.log("3. Policies on CCIPConsumer: Use getPolicies() to verify");
        console.log("");
        console.log("Next steps:");
        console.log("1. Blacklist test addresses if needed");
        console.log("2. Fund CCIP consumers with LINK");
        console.log("3. Update workflow config.json");
        console.log("4. Test mint + CCIP with PoR + ACE");
    }
}
