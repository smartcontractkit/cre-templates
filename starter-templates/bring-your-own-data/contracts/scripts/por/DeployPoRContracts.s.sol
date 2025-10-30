// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console} from "forge-std/Script.sol";
import {BundleAggregatorProxy} from "@chainlink/contracts/src/v0.8/data-feeds/BundleAggregatorProxy.sol";
import {DataFeedsCache} from "@chainlink/contracts/src/v0.8/data-feeds/DataFeedsCache.sol";
import {Workflow} from "../../src/util/Workflow.sol";

contract Deploy is Script {
    function run() external {
        vm.startBroadcast();

        string memory workflowName = "cre-por";
        address keystoneForwarder = Workflow.KeystoneForwarder(block.chainid);

        DataFeedsCache dfc = new DataFeedsCache();
        console.log("DataFeedsCache deployed at:", address(dfc));

        dfc.setFeedAdmin(msg.sender, true);

        uint256 n = 1;
        bytes16[] memory dataIds = new bytes16[](n);
        dataIds[0] = hex"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa";

        string[] memory descriptions = new string[](n);
        descriptions[0] = "PoR feed";

        uint8[][] memory decimalsMatrix = new uint8[][](n);
        decimalsMatrix[0] = new uint8[](n);
        decimalsMatrix[0][0] = 18; // 18 decimals

        DataFeedsCache.WorkflowMetadata[] memory workflowMetadata;

        if (vm.envOr("ENABLE_SIMULATE_WORKFLOW", false) == true) {
            address mockKeystoneForwarder = Workflow.MockKeystoneForwarder(block.chainid);

            // Enables writing data feeds cache with `cre workflow simulate`,
            // which uses a local signer and the mock keystone forwarder 
            // on-chain.

            bytes10 onChainSimWorkflowName = Workflow.ToOnChainSimWorkflowName(workflowName);
            console.log(
                string.concat("on-chain sim workflow name (bytes10): 0x", Workflow.ToHexString(abi.encodePacked(onChainSimWorkflowName)))
            );

            workflowMetadata = new DataFeedsCache.WorkflowMetadata[](2);
            workflowMetadata[1] = DataFeedsCache.WorkflowMetadata({
                allowedSender: mockKeystoneForwarder,
                allowedWorkflowOwner: 0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa, // simulator owner
                allowedWorkflowName: onChainSimWorkflowName
            });
        } else {
            workflowMetadata = new DataFeedsCache.WorkflowMetadata[](1);
        }

        bytes10 onChainWorkflowName = Workflow.ToOnChainWorkflowName(workflowName);
        console.log(
            string.concat("on-chain workflow name (bytes10): 0x", Workflow.ToHexString(abi.encodePacked(onChainWorkflowName)))
        );

        workflowMetadata[0] = DataFeedsCache.WorkflowMetadata({
            allowedSender: keystoneForwarder,
            allowedWorkflowOwner: msg.sender,
            allowedWorkflowName: onChainWorkflowName // <= 10 bytes
        });

        dfc.setBundleFeedConfigs(
            dataIds,
            descriptions,
            decimalsMatrix,
            workflowMetadata
        );

        BundleAggregatorProxy bap = new BundleAggregatorProxy(
            address(dfc),
            msg.sender
        );
        console.log("BundleAggregatorProxy deployed at:", address(bap));

        address[] memory proxies = new address[](n);
        proxies[0] = address(bap);
        dfc.updateDataIdMappingsForProxies(proxies, dataIds);
        console.log("Data ID mappings for proxies set in DataFeedsCache");

        vm.stopBroadcast();
    }
}
