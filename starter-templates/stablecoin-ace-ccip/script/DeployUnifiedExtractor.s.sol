// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "forge-std/Script.sol";
import "../contracts/extractors/UnifiedExtractor.sol";

/**
 * @title DeployUnifiedExtractor
 * @notice Deploys UnifiedExtractor and updates PolicyEngine
 * @dev Replaces both OnReportMintingExtractor and OnReportCCIPExtractor
 * 
 * Environment Variables Required:
 * - POLICY_ENGINE: Address of the deployed PolicyEngine
 * - PRIVATE_KEY: Deployer's private key
 * - SEPOLIA_RPC_URL: Sepolia RPC endpoint
 * 
 * Usage:
 * ETHERSCAN_API_KEY=dummy forge script script/DeployUnifiedExtractor.s.sol:DeployUnifiedExtractor \
 *   --rpc-url $SEPOLIA_RPC_URL \
 *   --private-key $PRIVATE_KEY \
 *   --broadcast
 */
contract DeployUnifiedExtractor is Script {
    function run() external {
        // Load environment variables
        address policyEngine = vm.envAddress("POLICY_ENGINE");
        address deployer = vm.addr(vm.envUint("PRIVATE_KEY"));
        
        console.log("============================================================");
        console.log("DEPLOYING UNIFIED EXTRACTOR");
        console.log("============================================================");
        console.log("");
        console.log("Deployer:", deployer);
        console.log("PolicyEngine:", policyEngine);
        console.log("");
        
        vm.startBroadcast();
        
        // 1. Deploy UnifiedExtractor
        console.log("Step 1: Deploying UnifiedExtractor...");
        UnifiedExtractor unifiedExtractor = new UnifiedExtractor();
        console.log("UnifiedExtractor deployed:", address(unifiedExtractor));
        
        // 2. Update PolicyEngine to use UnifiedExtractor for onReport selector
        console.log("");
        console.log("Step 2: Updating PolicyEngine...");
        bytes4 onReportSelector = bytes4(keccak256("onReport(bytes,bytes)"));
        console.log("Selector:", vm.toString(onReportSelector));
        
        IPolicyEngine(policyEngine).setExtractor(onReportSelector, address(unifiedExtractor));
        console.log("PolicyEngine updated successfully");
        
        // 3. Verify configuration
        console.log("");
        console.log("Step 3: Verifying configuration...");
        address currentExtractor = IPolicyEngine(policyEngine).getExtractor(onReportSelector);
        console.log("Current extractor for onReport:", currentExtractor);
        console.log("Expected:", address(unifiedExtractor));
        
        require(currentExtractor == address(unifiedExtractor), "Extractor verification failed");
        console.log("Verification: PASSED");
        
        vm.stopBroadcast();
        
        console.log("");
        console.log("============================================================");
        console.log("DEPLOYMENT COMPLETE");
        console.log("============================================================");
        console.log("");
        console.log("Summary:");
        console.log("  UnifiedExtractor:", address(unifiedExtractor));
        console.log("  Function selector:", vm.toString(onReportSelector));
        console.log("");
        console.log("Capabilities:");
        console.log("  Mint reports: instructionType, beneficiary, amount, bankRef");
        console.log("  CCIP reports: destinationChain, sender, beneficiary, amount, bankRef");
        console.log("");
        console.log("Add to .env:");
        console.log("  UNIFIED_EXTRACTOR=", address(unifiedExtractor));
        console.log("");
        console.log("Next steps:");
        console.log("  1. Test minting with blacklist policy");
        console.log("  2. Test CCIP with volume policy");
        console.log("");
    }
}

