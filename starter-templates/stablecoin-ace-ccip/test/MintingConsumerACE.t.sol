// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.26;

import {Test} from "forge-std/Test.sol";
import {console} from "forge-std/console.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {PolicyEngine} from "@chainlink/policy-management/core/PolicyEngine.sol";
import {AddressBlacklistPolicy} from "../contracts/policies/AddressBlacklistPolicy.sol";
import {OnReportMintingExtractor} from "../contracts/extractors/OnReportMintingExtractor.sol";
import {MintingConsumerWithACE} from "../contracts/MintingConsumerWithACE.sol";
import {StablecoinERC20} from "../contracts/StablecoinERC20.sol";
import {IPolicyEngine} from "@chainlink/policy-management/interfaces/IPolicyEngine.sol";

/**
 * @title MintingConsumerACETest
 * @notice Test ACE blocking exactly like ACE's own tests
 */
contract MintingConsumerACETest is Test {
    PolicyEngine public policyEngine;
    AddressBlacklistPolicy public blacklistPolicy;
    OnReportMintingExtractor public extractor;
    MintingConsumerWithACE public consumer;
    StablecoinERC20 public stablecoin;
    
    address public deployer;
    address public blacklistedAddr = 0x3333333333333333333333333333333333333333;
    address public normalAddr = 0x4444444444444444444444444444444444444444;
    
    function setUp() public {
        deployer = address(this);
        
        console.log("=== ACE Test Setup (Replicating ACE Pattern) ===");
        
        // 1. Deploy PolicyEngine
        PolicyEngine policyEngineImpl = new PolicyEngine();
        bytes memory policyEngineData = abi.encodeWithSelector(
            PolicyEngine.initialize.selector,
            true,  // defaultAllow = true (same as ACE tests)
            deployer
        );
        ERC1967Proxy policyEngineProxy = new ERC1967Proxy(address(policyEngineImpl), policyEngineData);
        policyEngine = PolicyEngine(address(policyEngineProxy));
        console.log("PolicyEngine:", address(policyEngine));
        
        // 2. Deploy BlacklistPolicy
        AddressBlacklistPolicy blacklistPolicyImpl = new AddressBlacklistPolicy();
        bytes memory blacklistPolicyData = abi.encodeWithSelector(
            blacklistPolicyImpl.initialize.selector,
            address(policyEngine),
            deployer,
            abi.encode(new address[](0))  // no initial blacklist
        );
        ERC1967Proxy blacklistPolicyProxy = new ERC1967Proxy(address(blacklistPolicyImpl), blacklistPolicyData);
        blacklistPolicy = AddressBlacklistPolicy(address(blacklistPolicyProxy));
        console.log("BlacklistPolicy:", address(blacklistPolicy));
        
        // 3. Deploy Stablecoin
        stablecoin = new StablecoinERC20("Test USD", "TUSD");
        console.log("Stablecoin:", address(stablecoin));
        
        // 4. Deploy MintingConsumer
        MintingConsumerWithACE consumerImpl = new MintingConsumerWithACE();
        bytes memory consumerData = abi.encodeWithSelector(
            consumerImpl.initialize.selector,
            deployer,
            address(stablecoin),
            address(policyEngine),
            address(0),  // no expected author
            bytes10("dummy")
        );
        ERC1967Proxy consumerProxy = new ERC1967Proxy(address(consumerImpl), consumerData);
        consumer = MintingConsumerWithACE(address(consumerProxy));
        console.log("Consumer:", address(consumer));
        
        // 5. Grant consumer mint/burn roles
        stablecoin.grantMintRole(address(consumer));
        stablecoin.grantBurnRole(address(consumer));
        
        // 6. Deploy Extractor
        extractor = new OnReportMintingExtractor();
        console.log("Extractor:", address(extractor));
        
        // 7. Configure ACE (exactly like ACE tests)
        bytes4 onReportSelector = bytes4(keccak256("onReport(bytes,bytes)"));
        policyEngine.setExtractor(onReportSelector, address(extractor));
        
        bytes32[] memory policyParams = new bytes32[](1);
        policyParams[0] = extractor.PARAM_BENEFICIARY();
        policyEngine.addPolicy(
            address(consumer),
            onReportSelector,
            address(blacklistPolicy),
            policyParams
        );
        
        console.log("");
        console.log("Setup complete!");
    }
    
    function test_mint_notBlacklisted_succeeds() public {
        console.log("=== Test 1: Mint to normal address (should succeed) ===");
        
        bytes memory report = abi.encode(
            uint8(1),      // MINT
            normalAddr,    // beneficiary
            uint256(100e18),  // amount
            bytes32("TEST")   // bankRef
        );
        
        consumer.onReport(bytes(""), report);
        
        uint256 balance = stablecoin.balanceOf(normalAddr);
        console.log("Balance:", balance);
        assertEq(balance, 100e18, "Mint should succeed");
    }
    
    function test_mint_blacklisted_reverts() public {
        console.log("=== Test 2: Mint to blacklisted address (should REVERT) ===");
        
        // Add to blacklist
        blacklistPolicy.addToBlacklist(blacklistedAddr);
        assertTrue(blacklistPolicy.isBlacklisted(blacklistedAddr), "Should be blacklisted");
        console.log("Blacklisted:", blacklistedAddr);
        
        // Try to mint
        bytes memory report = abi.encode(
            uint8(1),      // MINT
            blacklistedAddr,  // beneficiary (BLACKLISTED!)
            uint256(100e18),  // amount
            bytes32("TEST")   // bankRef
        );
        
        // EXPECT REVERT (exactly like ACE test line 106-109)
        vm.expectRevert(
            abi.encodeWithSelector(
                IPolicyEngine.PolicyRunRejected.selector,
                bytes4(keccak256("onReport(bytes,bytes)")),
                address(blacklistPolicy),
                "address is blacklisted"
            )
        );
        consumer.onReport(bytes(""), report);
        
        console.log("SUCCESS: Transaction reverted as expected!");
    }
}

