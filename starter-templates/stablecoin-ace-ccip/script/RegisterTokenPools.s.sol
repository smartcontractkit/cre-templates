// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Script, console} from "forge-std/Script.sol";

// Interfaces
interface IRegistryModuleOwnerCustom {
    function registerAdminViaOwner(address token) external;
}


interface ITokenAdminRegistry {
    function acceptAdminRole(address token) external;
    function setPool(address token, address pool) external;
    function getPool(address token) external view returns (address);
}

/**
 * @title RegisterTokenPools
 * @notice Registers TokenPools with TokenAdminRegistry on Sepolia and Fuji
 * @dev Replaces tokenpool-register-workflow (more reliable, no timeouts)
 */
contract RegisterTokenPools is Script {
    
    // Sepolia CCIP infrastructure
    address constant SEPOLIA_REGISTRY_MODULE = 0x62e731218d0D47305aba2BE3751E7EE9E5520790;
    address constant SEPOLIA_TOKEN_ADMIN_REGISTRY = 0x95F29FEE11c5C55d26cCcf1DB6772DE953B37B82;
    
    // Fuji CCIP infrastructure
    address constant FUJI_REGISTRY_MODULE = 0x34e503A31C6C085818703aB73b874BFF84AADF3C;
    address constant FUJI_TOKEN_ADMIN_REGISTRY = 0xbf8B49aDACaA58c98aB344F024C591E5E068E048;
    
    function run(bool isFuji) external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        
        address token;
        address pool;
        address registryModule;
        address tokenAdminRegistry;
        string memory chainName;
        string memory rpcUrl;
        
        if (isFuji) {
            token = vm.envAddress("STABLECOIN_FUJI");
            pool = vm.envAddress("POOL_FUJI");
            registryModule = FUJI_REGISTRY_MODULE;
            tokenAdminRegistry = FUJI_TOKEN_ADMIN_REGISTRY;
            chainName = "Fuji";
            rpcUrl = vm.envString("FUJI_RPC");
        } else {
            token = vm.envAddress("STABLECOIN_SEPOLIA");
            pool = vm.envAddress("POOL_SEPOLIA");
            registryModule = SEPOLIA_REGISTRY_MODULE;
            tokenAdminRegistry = SEPOLIA_TOKEN_ADMIN_REGISTRY;
            chainName = "Sepolia";
            rpcUrl = vm.envString("SEPOLIA_RPC");
        }
        
        console.log("\n===========================================");
        console.log("Registering TokenPool on", chainName);
        console.log("===========================================");
        console.log("Token:", token);
        console.log("Pool:", pool);
        console.log("\n[1/3] Registering admin via owner...");
        
        vm.broadcast(deployerPrivateKey);
        IRegistryModuleOwnerCustom(registryModule).registerAdminViaOwner(token);
        
        console.log("Success!");
        console.log("\n[2/3] Accepting admin role...");
        
        vm.broadcast(deployerPrivateKey);
        ITokenAdminRegistry(tokenAdminRegistry).acceptAdminRole(token);
        
        console.log("Success!");
        console.log("\n[3/3] Setting pool in TokenAdminRegistry...");
        
        vm.broadcast(deployerPrivateKey);
        ITokenAdminRegistry(tokenAdminRegistry).setPool(token, pool);
        
        console.log("Success!");
        
        // Verify registration
        address registeredPool = ITokenAdminRegistry(tokenAdminRegistry).getPool(token);
        require(registeredPool == pool, "Pool registration failed!");
        
        console.log("\n===========================================");
        console.log("TokenPool Registration Complete!");
        console.log("===========================================");
        console.log("Registered pool:", registeredPool);
        console.log("\nNext: Run ConfigureTokenPools.s.sol to set up cross-chain routes");
    }
}

