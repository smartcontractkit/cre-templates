// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console} from "forge-std/Script.sol";
import {MockPool} from "../../src/multi-chain-token-manager/MockPool.sol";
import {ProtocolSmartWallet} from "../../src/multi-chain-token-manager/ProtocolSmartWallet.sol";
import {LinkTokenInterface} from "@chainlink/contracts/src/v0.8/shared/interfaces/LinkTokenInterface.sol";
import {IERC20} from
    "@openzeppelin/contracts@5.0.2/token/ERC20/IERC20.sol";
import {Workflow} from "../../src/util/Workflow.sol";

contract Deploy is Script {
    function run() external {
        // Deploy on Sepolia
        uint256 sepoliaFork = vm.createSelectFork("sepolia");
        uint256 sepoliaChainID = block.chainid;
        vm.startBroadcast();
        address sepoliaPSW = deployAndDeposit();
        vm.stopBroadcast();
 
        // Deploy and configure on Base Sepolia
        vm.createSelectFork("base-sepolia");
        uint256 baseSepoliaChainID = block.chainid;
        vm.startBroadcast();
        address baseSepoliaPSW = deployAndDeposit();
        configureProtocolSmartWallet(baseSepoliaPSW, sepoliaPSW, getChainSelector(sepoliaChainID));
        vm.stopBroadcast();

        // Finally, configure on Sepolia now that the Base Sepolia
        // ProtocolSmartWallet contract address is known.
        vm.selectFork(sepoliaFork);
        vm.startBroadcast();
        configureProtocolSmartWallet(sepoliaPSW, baseSepoliaPSW, getChainSelector(baseSepoliaChainID));
        vm.stopBroadcast();
    }

    // Returns address of the deployed ProtocolSmartWallet.
    function deployAndDeposit() internal returns (address) {
        console.log("Deploying to chainid:", block.chainid);

        address[] memory keystoneForwarders;
        address[] memory workflowOwners;
        address ccipRouter;
        address linkToken;
        address bnmToken;
        uint128 apr;

        if (vm.envOr("ENABLE_WORKFLOW_SIMULATION", false) == true) {
            // Enables `cre workflow simulate` which uses the mock keystone 
            // forwarder on-chain.
            keystoneForwarders = new address[](2);
            address mockKeystoneForwarder = Workflow.mockKeystoneForwarder(block.chainid);
            keystoneForwarders[1] = mockKeystoneForwarder;
            workflowOwners = new address[](2);
            workflowOwners[1] = 0xaAaAaAaaAaAaAaaAaAAAAAAAAaaaAaAaAaaAaaAa; // simulator owner
        } else {
            keystoneForwarders = new address[](1);
            workflowOwners = new address[](1);
        }
        keystoneForwarders[0] = Workflow.keystoneForwarder(block.chainid);
        workflowOwners[0] = msg.sender;

        if (block.chainid == 11155111) {
            // ethereum sepolia
            ccipRouter = 0x0BF3dE8c5D3e8A2B34D2BEeB17ABfCeBaf363A59; 
            linkToken = 0x779877A7B0D9E8603169DdbD7836e478b4624789;
            bnmToken =  0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05;
            apr = 5e25; // 5% APR
        } else if (block.chainid == 84532) {
            // base sepolia
            ccipRouter = 0xD3b06cEbF099CE7DA4AcCf578aaebFDBd6e88a93;
            linkToken = 0xE4aB69C077896252FAFBD49EFD26B5D171A32410;
            bnmToken = 0x88A2d74F47a237a62e7A51cdDa67270CE381555e;
            apr = 6e25; // 6% APR
        }

        if (ccipRouter == address(0) || linkToken == address(0) || bnmToken == address(0) || keystoneForwarders.length == 0) {
            revert("Please configure deployment for this chainid");
        }

        MockPool pool = new MockPool();
        console.log("MockPool deployed at:", address(pool));

        ProtocolSmartWallet psw = new ProtocolSmartWallet(
            keystoneForwarders,
            workflowOwners,
            address(pool),
            ccipRouter,
            linkToken
        );
        console.log("Protocol Smart Wallet deployed at:", address(psw));

        pool.setCurrentLiquidityRate(bnmToken, apr); // Xe25 = X% APR
        console.log("Set MockPool BNM APR to:", apr / 1e25, "%");

        // Transfer LINK to PSW for CCIP send
        uint256 linkBal = LinkTokenInterface(linkToken).balanceOf(msg.sender);
        if (linkBal <= 1e18) {
            revert("Deployer requires at least 1 LINK, please use faucet");
        }
        console.log("Deployer LINK balance:", linkBal);
        uint256 amountOfLinkToFund = linkBal;
        if (amountOfLinkToFund > 1e18) { 
            amountOfLinkToFund = 1e18; // 1 LINK
        }

        console.log("Funding PSW with", amountOfLinkToFund, "LINK");
        LinkTokenInterface(linkToken).transfer(address(psw), amountOfLinkToFund);

        // Deposit BNM into the pool from the PSW
        uint256 bnmBal = IERC20(bnmToken).balanceOf(msg.sender);
        if (bnmBal == 0) {
            revert("Deployer has no BNM");
            
        }
        console.log("Deployer BNM balance:", bnmBal);
        uint256 amountOfBNMToDeposit = bnmBal;
        if (amountOfBNMToDeposit > 1e18) {
            amountOfBNMToDeposit = 1e18; // 1 BNM
        }

        IERC20(bnmToken).transfer(address(psw), amountOfBNMToDeposit);
        psw.depositToPool(bnmToken, amountOfBNMToDeposit);
        console.log(
            "Deposited",
            amountOfBNMToDeposit,
            "BNM into the Pool from PSW"
        );

        return address(psw);
    }

    function configureProtocolSmartWallet(address protocolSmartWallet, address sourceChainProtocolSmartWallet, uint64 sourceChainSelector) internal {

        console.log("Configuring ProtocolSmartWallet on chainid:", block.chainid);

        ProtocolSmartWallet psw = ProtocolSmartWallet(protocolSmartWallet);
        psw.setSenderForSourceChain(
            sourceChainSelector,
            sourceChainProtocolSmartWallet
        );

        console.log(
            "Set CCIP sender on chain %s for source chain selector %s to address %s",
            block.chainid,
            sourceChainSelector,
            address(sourceChainProtocolSmartWallet)
        );
    }

    function getChainSelector(uint256 chainid) internal pure returns (uint64) {
        if (chainid == 11155111) {
            return 16015286601757825753;
        } else if (chainid == 84532) {
            return 10344971235874465080;
        }
        revert("chainid not configured");
    }
}
