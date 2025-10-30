// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console} from "forge-std/Script.sol";
import {MockPool} from "../../src/multi-chain-token-manager/MockPool.sol";

contract Deploy is Script {
    function run() external {
        vm.startBroadcast();

        uint128 apr = 10e25; // 10% APR

        address poolAddress;
        address assetAddress;

        if (block.chainid == 11155111) {
            // ethereum sepolia
            poolAddress = 0x1Feac342C432770cD804C1EB8EC8ef514BAAf1E7; // Replace with address of your deployed MockPool contract
            assetAddress = 0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05; // BNM token address
        } else if (block.chainid == 84532) {
            // base sepolia
            poolAddress = 0xd79773a680885012863f1Ab2446a96fb830A825e; // Replace with address of your deployed MockPool contract
            assetAddress = 0x88A2d74F47a237a62e7A51cdDa67270CE381555e; // BNM token address
        }

        if (poolAddress == address(0) || assetAddress == address(0)) {
            revert("Please configure deployment for this chainid");
        }

        MockPool pool = MockPool(poolAddress);
        pool.setCurrentLiquidityRate(
            assetAddress, // BNM token address
            apr
        );

        uint256 aprPercent = (uint256(apr) * 100) / 1e27;
        console.log("Set asset %s APR on MockPool to %s%%", address(assetAddress), aprPercent);

        vm.stopBroadcast();
    }
}
