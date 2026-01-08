// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.26;

import {Script} from "forge-std/Script.sol";
import {console} from "forge-std/console.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {MintingConsumerWithACE} from "../contracts/MintingConsumerWithACE.sol";
import {CCIPTransferConsumerWithACE} from "../contracts/CCIPTransferConsumerWithACE.sol";
import {StablecoinERC20} from "../contracts/StablecoinERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title DeployACEConsumers
 * @notice Deploys ACE-protected consumer contracts for stablecoin operations.
 * @dev This script deploys:
 *      1. MintingConsumerWithACE (via proxy) - For mint/redeem with ACE
 *      2. CCIPTransferConsumerWithACE (via proxy) - For cross-chain transfers with ACE (Sepolia only)
 *      3. Grants MINTER_ROLE and BURN_ROLE to consumers
 *      4. Funds CCIP consumer with LINK
 * 
 * Prerequisites:
 * - StablecoinERC20 deployed on the target chain
 * - PolicyEngine deployed (from DeployACESystem.s.sol)
 * - Set environment variables in .env:
 *   - PRIVATE_KEY
 *   - STABLECOIN_SEPOLIA or STABLECOIN_FUJI
 *   - POLICY_ENGINE (from previous deployment)
 *   - SEPOLIA_ROUTER and LINK_SEPOLIA (for Sepolia only)
 *   - FUJI_ROUTER and LINK_FUJI (for Fuji only)
 * 
 * Usage (Sepolia):
 * forge script script/DeployACEConsumers.s.sol:DeployACEConsumers \
 *   --rpc-url $SEPOLIA_RPC --private-key $PRIVATE_KEY --broadcast \
 *   --sig "run(bool)" true
 * 
 * Usage (Fuji - no CCIP consumer):
 * forge script script/DeployACEConsumers.s.sol:DeployACEConsumers \
 *   --rpc-url $FUJI_RPC --private-key $PRIVATE_KEY --broadcast \
 *   --sig "run(bool)" false
 */
contract DeployACEConsumers is Script {
    /**
     * @notice Deploy ACE consumers.
     * @param deployCCIP Whether to deploy CCIPTransferConsumerWithACE (true for Sepolia, false for Fuji).
     */
    function run(bool deployCCIP) external {
        uint256 deployerPK = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPK);

        // Load environment variables
        address stablecoin;
        address policyEngine = vm.envAddress("POLICY_ENGINE");
        address router;
        address linkToken;

        // Determine chain and load appropriate addresses
        if (deployCCIP) {
            console.log("=== Deploying on Sepolia (with CCIP Consumer) ===");
            stablecoin = vm.envAddress("STABLECOIN_SEPOLIA");
            router = vm.envAddress("SEPOLIA_ROUTER");
            linkToken = vm.envAddress("LINK_SEPOLIA");
        } else {
            console.log("=== Deploying on Fuji (Minting Consumer Only) ===");
            stablecoin = vm.envAddress("STABLECOIN_FUJI");
            // Router and LINK not needed for Fuji (no CCIP consumer)
        }

        console.log("Deployer:", deployer);
        console.log("Stablecoin:", stablecoin);
        console.log("PolicyEngine:", policyEngine);
        console.log("");

        vm.startBroadcast(deployerPK);

        // ========================================
        // 1. Deploy MintingConsumerWithACE via Proxy
        // ========================================
        console.log("1. Deploying MintingConsumerWithACE...");
        
        MintingConsumerWithACE mintingConsumerImpl = new MintingConsumerWithACE();
        console.log("  MintingConsumer implementation:", address(mintingConsumerImpl));
        
        bytes memory mintingConsumerData = abi.encodeWithSelector(
            MintingConsumerWithACE.initialize.selector,
            deployer,                    // initialOwner
            stablecoin,                  // stablecoin
            policyEngine,                // policyEngine
            address(0),                  // expectedAuthor (0x0 for testing)
            bytes10("dummy")             // expectedWorkflowName (dummy for testing)
        );
        
        ERC1967Proxy mintingConsumerProxy = new ERC1967Proxy(
            address(mintingConsumerImpl),
            mintingConsumerData
        );
        MintingConsumerWithACE mintingConsumer = MintingConsumerWithACE(address(mintingConsumerProxy));
        console.log("  MintingConsumer proxy:", address(mintingConsumer));
        console.log("");

        // ========================================
        // 2. Grant Roles to MintingConsumer
        // ========================================
        console.log("2. Granting roles to MintingConsumer...");
        
        StablecoinERC20(stablecoin).grantMintRole(address(mintingConsumer));
        console.log("  Granted MINTER_ROLE to MintingConsumer");
        
        StablecoinERC20(stablecoin).grantBurnRole(address(mintingConsumer));
        console.log("  Granted BURN_ROLE to MintingConsumer");
        console.log("");

        // ========================================
        // 3. Deploy CCIPTransferConsumerWithACE (Sepolia only)
        // ========================================
        address ccipConsumer = address(0);
        
        if (deployCCIP) {
            console.log("3. Deploying CCIPTransferConsumerWithACE...");
            
            CCIPTransferConsumerWithACE ccipConsumerImpl = new CCIPTransferConsumerWithACE();
            console.log("  CCIPConsumer implementation:", address(ccipConsumerImpl));
            
            bytes memory ccipConsumerData = abi.encodeWithSelector(
                CCIPTransferConsumerWithACE.initialize.selector,
                deployer,                    // initialOwner
                stablecoin,                  // stablecoin
                router,                      // router
                linkToken,                   // linkToken
                policyEngine,                // policyEngine
                address(0),                  // expectedAuthor
                bytes10("dummy")             // expectedWorkflowName
            );
            
            ERC1967Proxy ccipConsumerProxy = new ERC1967Proxy(
                address(ccipConsumerImpl),
                ccipConsumerData
            );
            ccipConsumer = address(ccipConsumerProxy);
            console.log("  CCIPConsumer proxy:", ccipConsumer);
            console.log("");

            // ========================================
            // 4. Grant Roles to CCIPConsumer
            // ========================================
            console.log("4. Granting roles to CCIPConsumer...");
            
            // CCIP consumer doesn't need MINTER_ROLE (only transfers existing tokens)
            // But it does need approval from senders
            console.log("  Note: CCIP consumer requires users to approve it for token transfers");
            console.log("");
        }

        vm.stopBroadcast();

        // ========================================
        // Summary
        // ========================================
        console.log("=== Deployment Complete ===");
        console.log("");
        console.log("Deployed Addresses:");
        console.log("MintingConsumerWithACE:", address(mintingConsumer));
        if (deployCCIP) {
            console.log("CCIPTransferConsumerWithACE:", ccipConsumer);
        }
        console.log("");
        console.log("Next steps:");
        console.log("1. Update workflow config.json with these addresses");
        console.log("2. Attach extractors to PolicyEngine:");
        console.log("   - policyEngine.setExtractor(mintingConsumer._processReport.selector, mintingExtractor)");
        if (deployCCIP) {
            console.log("   - policyEngine.setExtractor(ccipConsumer._processReport.selector, ccipExtractor)");
        }
        console.log("3. Attach policies to functions:");
        console.log("   - policyEngine.addPolicy(mintingConsumer, selector, blacklistPolicy, ['beneficiary'])");
        if (deployCCIP) {
            console.log("   - policyEngine.addPolicy(ccipConsumer, selector, blacklistPolicy, ['recipient'])");
            console.log("4. Fund CCIP consumer with LINK:");
            console.log("   - cast send $LINK_SEPOLIA 'transfer(address,uint256)' ", ccipConsumer, " 5000000000000000000");
        }
    }
}

