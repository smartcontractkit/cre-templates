// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import "forge-std/Test.sol";
import "@chainlink/policy-management/interfaces/IPolicyEngine.sol";
import {VolumePolicy} from "@chainlink/policy-management/policies/VolumePolicy.sol";
import {OnReportCCIPExtractor} from "../contracts/extractors/OnReportCCIPExtractor.sol";

/**
 * @title VolumePolicyIntegrationTest
 * @notice Integration test simulating the full PolicyEngine flow
 */
contract VolumePolicyIntegrationTest is Test {
    VolumePolicy volumePolicy;
    OnReportCCIPExtractor extractor;
    
    function setUp() public {
        // Use deployed contracts
        volumePolicy = VolumePolicy(vm.envAddress("VOLUME_POLICY"));
        extractor = OnReportCCIPExtractor(vm.envAddress("CCIP_EXTRACTOR"));
    }
    
    function testExtractorOutput() public view {
        console.log("============================================================");
        console.log("TEST 1: EXTRACTOR OUTPUT");
        console.log("============================================================");
        console.log("");
        
        // Simulate what the consumer receives
        uint64 destChain = 14767482510784806043; // Fuji
        address sender = address(0x4fed0A5B65eac383D36E65733786386709B86be8);
        address beneficiary = address(0x742D35CC6634c0532925A3b844BC9E7595F0BEb0);
        uint256 amount = 500e18; // 500 creUSD
        bytes32 bankRef = bytes32(0);
        
        // Encode the report (what CRE sends)
        bytes memory report = abi.encode(destChain, sender, beneficiary, amount, bankRef);
        
        // Encode onReport parameters
        bytes memory metadata = "";
        bytes memory data = abi.encode(metadata, report);
        
        // Create payload (what PolicyEngine receives)
        IPolicyEngine.Payload memory payload = IPolicyEngine.Payload({
            selector: bytes4(keccak256("onReport(bytes,bytes)")),
            sender: sender,
            data: data,
            context: ""
        });
        
        // Call extractor
        IPolicyEngine.Parameter[] memory params = extractor.extract(payload);
        
        console.log("Extractor returned", params.length, "parameters:");
        for (uint i = 0; i < params.length; i++) {
            console.log("  Parameter", i, ":");
            console.log("    Name:");
            console.logBytes32(params[i].name);
            
            if (params[i].name == keccak256("amount")) {
                uint256 extractedAmount = abi.decode(params[i].value, (uint256));
                console.log("    Value (amount):", extractedAmount);
            }
        }
        
        console.log("");
        console.log("Expected 'amount' parameter:");
        console.logBytes32(keccak256("amount"));
    }
    
    function testPolicyWithExtractorOutput() public view {
        console.log("============================================================");
        console.log("TEST 2: POLICY WITH EXTRACTOR OUTPUT");
        console.log("============================================================");
        console.log("");
        
        // Simulate extractor output
        uint64 destChain = 14767482510784806043;
        address sender = address(0x4fed0A5B65eac383D36E65733786386709B86be8);
        address beneficiary = address(0x742D35CC6634c0532925A3b844BC9E7595F0BEb0);
        uint256 amount = 500e18;
        bytes32 bankRef = bytes32(0);
        
        bytes memory report = abi.encode(destChain, sender, beneficiary, amount, bankRef);
        bytes memory metadata = "";
        bytes memory data = abi.encode(metadata, report);
        
        IPolicyEngine.Payload memory payload = IPolicyEngine.Payload({
            selector: bytes4(keccak256("onReport(bytes,bytes)")),
            sender: sender,
            data: data,
            context: ""
        });
        
        IPolicyEngine.Parameter[] memory allParams = extractor.extract(payload);
        
        // Simulate what PolicyEngine does: filter parameters
        bytes32 requestedParam = keccak256("amount");
        bytes[] memory filteredParams = new bytes[](1);
        bool found = false;
        
        for (uint i = 0; i < allParams.length; i++) {
            if (allParams[i].name == requestedParam) {
                filteredParams[0] = allParams[i].value;
                found = true;
                break;
            }
        }
        
        console.log("Parameter filtering:");
        console.log("  Requested:", vm.toString(requestedParam));
        console.log("  Found:", found);
        
        if (!found) {
            console.log("  ERROR: Parameter not found in extractor output!");
            return;
        }
        
        // Call policy with filtered parameters
        console.log("");
        console.log("Calling VolumePolicy with filtered parameters...");
        
        try volumePolicy.run(sender, address(0), bytes4(0), filteredParams, "") returns (
            IPolicyEngine.PolicyResult result
        ) {
            console.log("  Result: PolicyResult.", uint8(result));
            console.log("  Status: ALLOWED (expected for 500 creUSD)");
        } catch Error(string memory reason) {
            console.log("  Reverted with reason:", reason);
            console.log("  Status: BLOCKED (unexpected!)");
        } catch {
            console.log("  Reverted without reason");
            console.log("  Status: BLOCKED (unexpected!)");
        }
    }
}

