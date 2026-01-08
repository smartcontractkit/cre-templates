// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import "@chainlink/policy-management/policies/VolumePolicy.sol";

/**
 * @title DeployVolumePolicy
 * @notice Deploys VolumePolicy for CCIP transfer amount validation
 * @dev Deploys via proxy pattern with configurable min/max amounts
 * 
 * Environment Variables Required:
 * - POLICY_ENGINE: Address of the deployed PolicyEngine
 * - PRIVATE_KEY: Deployer's private key
 * - SEPOLIA_RPC_URL: Sepolia RPC endpoint
 * 
 * Usage:
 * forge script script/DeployVolumePolicy.s.sol:DeployVolumePolicy \
 *   --rpc-url $SEPOLIA_RPC_URL \
 *   --private-key $PRIVATE_KEY \
 *   --broadcast \
 *   --verify
 */
contract DeployVolumePolicy is Script {
    function run() external {
        // Load environment variables
        address policyEngine = vm.envAddress("POLICY_ENGINE");
        address deployer = vm.addr(vm.envUint("PRIVATE_KEY"));
        
        console.log("============================================================");
        console.log("DEPLOYING VOLUMEPOLICY FOR CCIP");
        console.log("============================================================");
        console.log("");
        console.log("Deployer:", deployer);
        console.log("PolicyEngine:", policyEngine);
        console.log("");
        
        vm.startBroadcast();
        
        // 1. Deploy VolumePolicy implementation
        console.log("Step 1: Deploying VolumePolicy implementation...");
        VolumePolicy volumePolicyImpl = new VolumePolicy();
        console.log("VolumePolicy implementation:", address(volumePolicyImpl));
        
        // 2. Configure min/max amounts
        // Min: 100 creUSD (100e18 wei)
        // Max: 10,000 creUSD (10000e18 wei)
        uint256 minAmount = 100e18;
        uint256 maxAmount = 10000e18;
        
        console.log("");
        console.log("Configuration:");
        console.log("  Min amount: 100 creUSD (100e18 wei)");
        console.log("  Max amount: 10,000 creUSD (10000e18 wei)");
        
        // 3. Encode initialize call
        bytes memory initData = abi.encodeWithSignature(
            "initialize(address,address,bytes)",
            policyEngine,
            deployer,
            abi.encode(minAmount, maxAmount)
        );
        
        // 4. Deploy proxy
        console.log("");
        console.log("Step 2: Deploying VolumePolicy proxy...");
        ERC1967Proxy volumePolicyProxy = new ERC1967Proxy(
            address(volumePolicyImpl),
            initData
        );
        console.log("VolumePolicy proxy:", address(volumePolicyProxy));
        
        // 5. Verify configuration
        VolumePolicy volumePolicy = VolumePolicy(address(volumePolicyProxy));
        uint256 deployedMin = volumePolicy.getMin();
        uint256 deployedMax = volumePolicy.getMax();
        
        console.log("");
        console.log("Step 3: Verifying configuration...");
        console.log("  Deployed min:", deployedMin);
        console.log("  Deployed max:", deployedMax);
        
        require(deployedMin == minAmount, "Min amount mismatch");
        require(deployedMax == maxAmount, "Max amount mismatch");
        console.log("  Verification: PASSED");
        
        vm.stopBroadcast();
        
        console.log("");
        console.log("============================================================");
        console.log("DEPLOYMENT COMPLETE");
        console.log("============================================================");
        console.log("");
        console.log("Summary:");
        console.log("  VolumePolicy Implementation:", address(volumePolicyImpl));
        console.log("  VolumePolicy Proxy:", address(volumePolicyProxy));
        console.log("  Min Amount:", minAmount, "wei (100 creUSD)");
        console.log("  Max Amount:", maxAmount, "wei (10,000 creUSD)");
        console.log("");
        console.log("Add to .env:");
        console.log("  VOLUME_POLICY=", address(volumePolicyProxy));
        console.log("");
    }
}

