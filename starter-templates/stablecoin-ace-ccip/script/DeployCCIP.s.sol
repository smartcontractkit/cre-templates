// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

/**
 * @notice All-in-one CCIP deployment script (Workshop Style)
 * @dev Deploys and configures TokenPools on both Sepolia and Fuji in a single script.
 *      Uses multi-fork mode to switch between chains automatically.
 */
import {Script, console} from "forge-std/Script.sol";
import {BurnMintTokenPool} from "@chainlink/contracts-ccip/ccip/pools/BurnMintTokenPool.sol";
import {TokenAdminRegistry} from "@chainlink/contracts-ccip/ccip/tokenAdminRegistry/TokenAdminRegistry.sol";
import {RegistryModuleOwnerCustom} from "@chainlink/contracts-ccip/ccip/tokenAdminRegistry/RegistryModuleOwnerCustom.sol";
import {RateLimiter} from "@chainlink/contracts-ccip/ccip/libraries/RateLimiter.sol";
import {TokenPool} from "@chainlink/contracts-ccip/ccip/pools/TokenPool.sol";
import {IBurnMintERC20} from "@chainlink/contracts/src/v0.8/shared/token/ERC20/IBurnMintERC20.sol";

/**
 * @title DeployCCIP
 * @notice All-in-one CCIP setup script for CRE stablecoin demo
 * @dev Deploys and configures TokenPools on Sepolia and Fuji with automatic fork switching
 * 
 * What this script does:
 * 1. Deploys BurnMintTokenPool on Sepolia
 * 2. Grants mint/burn roles to Sepolia pool
 * 3. Registers Sepolia pool with TokenAdminRegistry
 * 4. Deploys BurnMintTokenPool on Fuji
 * 5. Grants mint/burn roles to Fuji pool
 * 6. Registers Fuji pool with TokenAdminRegistry
 * 7. Configures bidirectional routes (Sepolia ↔ Fuji)
 * 
 * Usage:
 * ETHERSCAN_API_KEY=dummy forge script script/DeployCCIP.s.sol:DeployCCIP \
 *   --broadcast --slow
 */
contract DeployCCIP is Script {
    
    // ========================================
    // Network Configurations
    // ========================================
    
    // Sepolia
    uint64 constant SEPOLIA_CHAIN_SELECTOR = 16015286601757825753;
    address constant SEPOLIA_RMN_PROXY = 0xba3f6251de62dED61Ff98590cB2fDf6871FbB991;
    address constant SEPOLIA_ROUTER = 0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59;
    address constant SEPOLIA_REGISTRY_MODULE = 0x62e731218d0D47305aba2BE3751E7EE9E5520790;
    address constant SEPOLIA_TOKEN_ADMIN_REGISTRY = 0x95F29FEE11c5C55d26cCcf1DB6772DE953B37B82;
    
    // Fuji
    uint64 constant FUJI_CHAIN_SELECTOR = 14767482510784806043;
    address constant FUJI_RMN_PROXY = 0xAc8CFc3762a979628334a0E4C1026244498E821b;
    address constant FUJI_ROUTER = 0xF694E193200268f9a4868e4Aa017A0118C9a8177;
    address constant FUJI_REGISTRY_MODULE = 0x97300785aF1edE1343DB6d90706A35CF14aA3d81;
    address constant FUJI_TOKEN_ADMIN_REGISTRY = 0xA92053a4a3922084d992fD2835bdBa4caC6877e6;
    
    // State
    address public sepoliaPool;
    address public fujiPool;
    
    function run() external {
        // Load from environment
        address sepoliaToken = vm.envAddress("STABLECOIN_SEPOLIA");
        address fujiToken = vm.envAddress("STABLECOIN_FUJI");
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        string memory sepoliaRpc = vm.envString("SEPOLIA_RPC");
        string memory fujiRpc = vm.envString("FUJI_RPC");
        
        console.log("============================================================");
        console.log("DEPLOYING CCIP INFRASTRUCTURE (SEPOLIA <> FUJI)");
        console.log("============================================================");
        console.log("");
        console.log("Tokens:");
        console.log("  Sepolia:", sepoliaToken);
        console.log("  Fuji:", fujiToken);
        console.log("");
        
        // ========================================
        // SEPOLIA DEPLOYMENT
        // ========================================
        console.log("-----------------------------------------------------------");
        console.log("PHASE 1: SEPOLIA SETUP");
        console.log("-----------------------------------------------------------");
        uint256 sepoliaFork = vm.createSelectFork(sepoliaRpc);
        vm.selectFork(sepoliaFork);
        
        vm.startBroadcast(deployerKey);
        sepoliaPool = deployTokenPool(sepoliaToken, SEPOLIA_RMN_PROXY, SEPOLIA_ROUTER, "Sepolia");
        grantRoles(sepoliaToken, sepoliaPool);
        registerAdmin(SEPOLIA_REGISTRY_MODULE, sepoliaToken);
        acceptAdmin(SEPOLIA_TOKEN_ADMIN_REGISTRY, sepoliaToken);
        setPool(SEPOLIA_TOKEN_ADMIN_REGISTRY, sepoliaToken, sepoliaPool);
        vm.stopBroadcast();
        
        console.log("Sepolia TokenPool:", sepoliaPool);
        console.log("");
        
        // ========================================
        // FUJI DEPLOYMENT
        // ========================================
        console.log("-----------------------------------------------------------");
        console.log("PHASE 2: FUJI SETUP");
        console.log("-----------------------------------------------------------");
        uint256 fujiFork = vm.createSelectFork(fujiRpc);
        vm.selectFork(fujiFork);
        
        vm.startBroadcast(deployerKey);
        fujiPool = deployTokenPool(fujiToken, FUJI_RMN_PROXY, FUJI_ROUTER, "Fuji");
        grantRoles(fujiToken, fujiPool);
        registerAdmin(FUJI_REGISTRY_MODULE, fujiToken);
        acceptAdmin(FUJI_TOKEN_ADMIN_REGISTRY, fujiToken);
        setPool(FUJI_TOKEN_ADMIN_REGISTRY, fujiToken, fujiPool);
        vm.stopBroadcast();
        
        console.log("Fuji TokenPool:", fujiPool);
        console.log("");
        
        // ========================================
        // CONFIGURE CROSS-CHAIN ROUTES
        // ========================================
        console.log("-----------------------------------------------------------");
        console.log("PHASE 3: CONFIGURE BIDIRECTIONAL ROUTES");
        console.log("-----------------------------------------------------------");
        
        // Sepolia → Fuji route
        vm.selectFork(sepoliaFork);
        vm.startBroadcast(deployerKey);
        configureRoute(sepoliaPool, FUJI_CHAIN_SELECTOR, fujiPool, fujiToken, "Sepolia -> Fuji");
        vm.stopBroadcast();
        
        // Fuji → Sepolia route
        vm.selectFork(fujiFork);
        vm.startBroadcast(deployerKey);
        configureRoute(fujiPool, SEPOLIA_CHAIN_SELECTOR, sepoliaPool, sepoliaToken, "Fuji -> Sepolia");
        vm.stopBroadcast();
        
        // ========================================
        // FINAL SUMMARY
        // ========================================
        console.log("");
        console.log("============================================================");
        console.log("CCIP DEPLOYMENT COMPLETE");
        console.log("============================================================");
        console.log("");
        console.log("Sepolia:");
        console.log("  Token:", sepoliaToken);
        console.log("  Pool:", sepoliaPool);
        console.log("");
        console.log("Fuji:");
        console.log("  Token:", fujiToken);
        console.log("  Pool:", fujiPool);
        console.log("");
        console.log("Add to .env:");
        console.log("  POOL_SEPOLIA=", sepoliaPool);
        console.log("  POOL_FUJI=", fujiPool);
        console.log("");
        console.log("Routes configured: Sepolia <-> Fuji (bidirectional)");
        console.log("TokenPools are ready for CCIP transfers!");
        console.log("");
    }
    
    function deployTokenPool(
        address token,
        address rmnProxy,
        address router,
        string memory chainName
    ) internal returns (address) {
        console.log("[", chainName, "] Deploying BurnMintTokenPool...");
        
        BurnMintTokenPool pool = new BurnMintTokenPool(
            IBurnMintERC20(token),
            18, // decimals
            new address[](0), // allowlist (empty)
            rmnProxy,
            router
        );
        
        console.log("[", chainName, "] Pool deployed:", address(pool));
        return address(pool);
    }
    
    function grantRoles(address token, address pool) internal {
        console.log("  Granting mint/burn roles to pool...");
        
        // Call grantMintAndBurnRoles (our StablecoinERC20 has this convenience function)
        (bool success, ) = token.call(
            abi.encodeWithSignature("grantMintAndBurnRoles(address)", pool)
        );
        require(success, "Failed to grant roles");
        
        console.log("  Roles granted");
    }
    
    function registerAdmin(address registryModule, address token) internal {
        console.log("  Registering admin via RegistryModuleOwnerCustom...");
        RegistryModuleOwnerCustom(registryModule).registerAdminViaOwner(token);
        console.log("  Admin registered");
    }
    
    function acceptAdmin(address registry, address token) internal {
        console.log("  Accepting admin role...");
        TokenAdminRegistry(registry).acceptAdminRole(token);
        console.log("  Admin role accepted");
    }
    
    function setPool(address registry, address token, address pool) internal {
        console.log("  Setting pool in TokenAdminRegistry...");
        TokenAdminRegistry(registry).setPool(token, pool);
        console.log("  Pool set");
    }
    
    function configureRoute(
        address localPool,
        uint64 remoteChainSelector,
        address remotePool,
        address remoteToken,
        string memory routeName
    ) internal {
        console.log("[Route]", routeName);
        console.log("  Remote chain selector:", remoteChainSelector);
        console.log("  Remote pool:", remotePool);
        
        TokenPool.ChainUpdate[] memory updates = new TokenPool.ChainUpdate[](1);
        bytes[] memory remotePoolAddresses = new bytes[](1);
        remotePoolAddresses[0] = abi.encode(remotePool);
        
        updates[0] = TokenPool.ChainUpdate({
            remoteChainSelector: remoteChainSelector,
            remotePoolAddresses: remotePoolAddresses,
            remoteTokenAddress: abi.encode(remoteToken),
            outboundRateLimiterConfig: RateLimiter.Config({
                isEnabled: false,
                capacity: 0,
                rate: 0
            }),
            inboundRateLimiterConfig: RateLimiter.Config({
                isEnabled: false,
                capacity: 0,
                rate: 0
            })
        });
        
        TokenPool(localPool).applyChainUpdates(new uint64[](0), updates);
        console.log("  Route configured");
    }
}
