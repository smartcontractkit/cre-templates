// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import "@chainlink/policy-management/interfaces/IPolicyEngine.sol";

/**
 * @title PolicyEngineRunTest
 * @notice Test PolicyEngine.run() with actual onReport payload
 */
contract PolicyEngineRunTest is Test {
    IPolicyEngine policyEngine;
    address ccipConsumer;
    
    function setUp() public {
        policyEngine = IPolicyEngine(vm.envAddress("POLICY_ENGINE"));
        ccipConsumer = vm.envAddress("CCIP_CONSUMER_ACE");
    }
    
    function testPolicyEngineRun500() public {
        console.log("============================================================");
        console.log("TEST: PolicyEngine.run() with 500 creUSD");
        console.log("============================================================");
        console.log("");
        console.log("PolicyEngine:", address(policyEngine));
        console.log("CCIP Consumer:", ccipConsumer);
        console.log("");
        
        // Simulate onReport call with 500 creUSD
        uint64 destChain = 14767482510784806043;
        address sender = address(0x4fed0A5B65eac383D36E65733786386709B86be8);
        address beneficiary = address(0x742D35CC6634c0532925A3b844BC9E7595F0BEb0);
        uint256 amount = 500e18;
        bytes32 bankRef = bytes32(0);
        
        bytes memory report = abi.encode(destChain, sender, beneficiary, amount, bankRef);
        bytes memory metadata = "";
        
        // Create the data that would be passed to onReport (both parameters)
        bytes memory onReportData = abi.encode(metadata, report);
        
        // Create payload (what PolicyProtected.runPolicy sends to PolicyEngine)
        IPolicyEngine.Payload memory payload = IPolicyEngine.Payload({
            selector: bytes4(keccak256("onReport(bytes,bytes)")),
            sender: sender,
            data: onReportData,
            context: ""
        });
        
        console.log("Calling PolicyEngine.run()...");
        console.log("  Selector:", vm.toString(payload.selector));
        console.log("  Sender:", payload.sender);
        console.log("  Amount: 500 creUSD");
        console.log("");
        
        // Set msg.sender to match what the consumer would send
        vm.prank(ccipConsumer);
        
        try policyEngine.run(payload) {
            console.log("Result: SUCCESS (policy allowed)");
            console.log("Expected: SUCCESS for 500 creUSD");
            console.log("Status: PASS");
        } catch Error(string memory reason) {
            console.log("Result: REVERTED");
            console.log("Reason:", reason);
            console.log("Expected: SUCCESS");
            console.log("Status: FAIL");
        } catch (bytes memory lowLevelData) {
            console.log("Result: REVERTED (low-level)");
            console.log("Data length:", lowLevelData.length);
            if (lowLevelData.length > 0) {
                console.logBytes(lowLevelData);
            }
            console.log("Expected: SUCCESS");
            console.log("Status: FAIL");
        }
    }
    
    function testPolicyEngineRun50() public {
        console.log("============================================================");
        console.log("TEST: PolicyEngine.run() with 50 creUSD");
        console.log("============================================================");
        console.log("");
        
        // Simulate onReport call with 50 creUSD (should be blocked)
        uint64 destChain = 14767482510784806043;
        address sender = address(0x4fed0A5B65eac383D36E65733786386709B86be8);
        address beneficiary = address(0x742D35CC6634c0532925A3b844BC9E7595F0BEb0);
        uint256 amount = 50e18;
        bytes32 bankRef = bytes32(0);
        
        bytes memory report = abi.encode(destChain, sender, beneficiary, amount, bankRef);
        bytes memory metadata = "";
        bytes memory onReportData = abi.encode(metadata, report);
        
        IPolicyEngine.Payload memory payload = IPolicyEngine.Payload({
            selector: bytes4(keccak256("onReport(bytes,bytes)")),
            sender: sender,
            data: onReportData,
            context: ""
        });
        
        console.log("Calling PolicyEngine.run() with 50 creUSD...");
        
        vm.prank(ccipConsumer);
        
        try policyEngine.run(payload) {
            console.log("Result: SUCCESS (unexpected!)");
            console.log("Expected: REVERTED");
            console.log("Status: FAIL");
        } catch Error(string memory reason) {
            console.log("Result: REVERTED");
            console.log("Reason:", reason);
            console.log("Expected: REVERTED for 50 creUSD");
            console.log("Status: PASS");
        } catch {
            console.log("Result: REVERTED (expected)");
            console.log("Status: PASS");
        }
    }
}

