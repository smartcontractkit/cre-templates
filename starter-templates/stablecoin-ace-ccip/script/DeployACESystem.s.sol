// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {PolicyEngine} from "@chainlink/policy-management/core/PolicyEngine.sol";
import {Policy} from "@chainlink/policy-management/core/Policy.sol";
import {AddressBlacklistPolicy} from "../contracts/policies/AddressBlacklistPolicy.sol";

/**
 * @title DeployACESystem
 * @notice Deploys the complete ACE infrastructure for the stablecoin demo.
 * @dev This script deploys:
 *      1. PolicyEngine (via proxy) - Central policy orchestrator
 *      2. AddressBlacklistPolicy (via proxy) - Blacklist enforcement
 * 
 * Prerequisites:
 * - Set PRIVATE_KEY in .env
 * - Have testnet ETH for gas
 * 
 * Usage:
 * forge script script/DeployACESystem.s.sol:DeployACESystem \
 *   --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast
 * 
 * After deployment:
 * - Save all deployed addresses to .env
 * - PolicyEngine address needed for UnifiedExtractor deployment
 * - Policy address needed for blacklist management
 */
contract DeployACESystem is Script {
    function run() external {
        uint256 deployerPK = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPK);

        vm.startBroadcast(deployerPK);

        console.log("=== Deploying ACE Infrastructure ===");
        console.log("Deployer:", deployer);
        console.log("");

        // ========================================
        // 1. Deploy PolicyEngine via Proxy
        // ========================================
        console.log("1. Deploying PolicyEngine...");
        
        PolicyEngine policyEngineImpl = new PolicyEngine();
        console.log("  PolicyEngine implementation:", address(policyEngineImpl));
        
        bytes memory policyEngineData = abi.encodeWithSelector(
            PolicyEngine.initialize.selector,
            true,     // defaultAllow = true (allow by default unless policy rejects)
            deployer  // owner
        );
        
        ERC1967Proxy policyEngineProxy = new ERC1967Proxy(
            address(policyEngineImpl),
            policyEngineData
        );
        PolicyEngine policyEngine = PolicyEngine(address(policyEngineProxy));
        console.log("  PolicyEngine proxy:", address(policyEngine));
        console.log("");

        // ========================================
        // 2. Deploy AddressBlacklistPolicy via Proxy
        // ========================================
        console.log("2. Deploying AddressBlacklistPolicy...");
        
        AddressBlacklistPolicy blacklistPolicyImpl = new AddressBlacklistPolicy();
        console.log("  BlacklistPolicy implementation:", address(blacklistPolicyImpl));
        
        // No initial blacklist (empty configuration)
        bytes memory blacklistPolicyConfig = abi.encode(new address[](0));
        
        bytes memory blacklistPolicyData = abi.encodeWithSelector(
            Policy.initialize.selector,
            address(policyEngine),  // policyEngine
            deployer,               // owner
            blacklistPolicyConfig   // config
        );
        
        ERC1967Proxy blacklistPolicyProxy = new ERC1967Proxy(
            address(blacklistPolicyImpl),
            blacklistPolicyData
        );
        AddressBlacklistPolicy blacklistPolicy = AddressBlacklistPolicy(address(blacklistPolicyProxy));
        console.log("  BlacklistPolicy proxy:", address(blacklistPolicy));
        console.log("");

        vm.stopBroadcast();

        // ========================================
        // Summary
        // ========================================
        console.log("=== Deployment Complete ===");
        console.log("");
        console.log("Copy these addresses to your .env file:");
        console.log("");
        console.log("# ACE Infrastructure");
        console.log("POLICY_ENGINE=", address(policyEngine));
        console.log("BLACKLIST_POLICY=", address(blacklistPolicy));
        console.log("");
        console.log("Next steps:");
        console.log("1. Update .env with the above addresses");
        console.log("2. Deploy UnifiedExtractor (script/DeployUnifiedExtractor.s.sol)");
        console.log("3. Deploy VolumePolicy (script/DeployVolumePolicy.s.sol)");
        console.log("4. Deploy ACE consumers (script/DeployACEConsumers.s.sol)");
    }
}

