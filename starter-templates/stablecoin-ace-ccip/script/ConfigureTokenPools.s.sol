// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console} from "forge-std/Script.sol";

// TokenPool interface
interface ITokenPool {
        struct ChainUpdate {
            uint64 remoteChainSelector;
            bytes[] remotePoolAddresses;
            bytes remoteTokenAddress;
            RateLimitConfig outboundRateLimiterConfig;
            RateLimitConfig inboundRateLimiterConfig;
        }
        
        struct RateLimitConfig {
            bool isEnabled;
            uint128 capacity;
            uint128 rate;
        }
        
        function applyChainUpdates(
            uint64[] memory chainsToRemove,
            ChainUpdate[] memory chainsToAdd
        ) external;
        
        function isSupportedChain(uint64 remoteChainSelector) external view returns (bool);
}

/**
 * @title ConfigureTokenPools
 * @notice Configures cross-chain routes for TokenPools (Sepolia ↔ Fuji)
 * @dev Replaces tokenpool-configure-workflow (more reliable, no timeouts)
 */
contract ConfigureTokenPools is Script {
    // Chain selectors
    uint64 constant SEPOLIA_CHAIN_SELECTOR = 16015286601757825753;
    uint64 constant FUJI_CHAIN_SELECTOR = 14767482510784806043;
    
    function run(bool isFuji) external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        address localPool;
        address remotePool;
        address remoteToken;
        uint64 remoteChainSelector;
        string memory direction;
        string memory rpcUrl;
        
        if (isFuji) {
            // Fuji → Sepolia
            localPool = vm.envAddress("POOL_FUJI");
            remotePool = vm.envAddress("POOL_SEPOLIA");
            remoteToken = vm.envAddress("STABLECOIN_SEPOLIA");
            remoteChainSelector = SEPOLIA_CHAIN_SELECTOR;
            direction = "Fuji -> Sepolia";
            rpcUrl = vm.envString("FUJI_RPC");
        } else {
            // Sepolia → Fuji
            localPool = vm.envAddress("POOL_SEPOLIA");
            remotePool = vm.envAddress("POOL_FUJI");
            remoteToken = vm.envAddress("STABLECOIN_FUJI");
            remoteChainSelector = FUJI_CHAIN_SELECTOR;
            direction = "Sepolia -> Fuji";
            rpcUrl = vm.envString("SEPOLIA_RPC");
        }
        
        console.log("\n===========================================");
        console.log("Configuring TokenPool Route:", direction);
        console.log("===========================================");
        console.log("Local Pool:", localPool);
        console.log("Remote Pool:", remotePool);
        console.log("Remote Token:", remoteToken);
        console.log("Remote Chain Selector:", remoteChainSelector);
        
        // Prepare ChainUpdate
        ITokenPool.ChainUpdate[] memory updates = new ITokenPool.ChainUpdate[](1);
        
        bytes[] memory remotePoolAddresses = new bytes[](1);
        remotePoolAddresses[0] = abi.encode(remotePool);
        
        updates[0] = ITokenPool.ChainUpdate({
            remoteChainSelector: remoteChainSelector,
            remotePoolAddresses: remotePoolAddresses,
            remoteTokenAddress: abi.encode(remoteToken),
            outboundRateLimiterConfig: ITokenPool.RateLimitConfig({
                isEnabled: false,
                capacity: 0,
                rate: 0
            }),
            inboundRateLimiterConfig: ITokenPool.RateLimitConfig({
                isEnabled: false,
                capacity: 0,
                rate: 0
            })
        });
        
        console.log("\nApplying chain updates...");
        
        vm.broadcast(deployerPrivateKey);
        ITokenPool(localPool).applyChainUpdates(new uint64[](0), updates);
        
        console.log("Success!");
        
        // Verify configuration
        bool isSupported = ITokenPool(localPool).isSupportedChain(remoteChainSelector);
        require(isSupported, "Route configuration failed!");
        
        console.log("\n===========================================");
        console.log("Route Configuration Complete!");
        console.log("===========================================");
        console.log("Route:", direction);
        console.log("Verified: Remote chain is now supported");
        console.log("\nTokenPools are ready for CCIP transfers!");
    }
}

