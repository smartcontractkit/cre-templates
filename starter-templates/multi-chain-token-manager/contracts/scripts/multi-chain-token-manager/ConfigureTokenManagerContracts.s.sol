// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {Script, console} from "forge-std/Script.sol";
import {ProtocolSmartWallet} from "../../src/multi-chain-token-manager/ProtocolSmartWallet.sol";

contract Deploy is Script {
    function run() external {
        vm.startBroadcast();

        address protocolSmartWallet = getProtocolSmartWallet(block.chainid);
        uint64 sourceChainSelector;
        address sourceChainProtocolSmartWallet;

        if (block.chainid == 11155111) {
            // ethereum sepolia; source will be base sepolia
            uint256 sourceChainID = 84532;
            sourceChainSelector = getChainSelector(sourceChainID);
            sourceChainProtocolSmartWallet = getProtocolSmartWallet(sourceChainID);
        } else if (block.chainid == 84532) {
            // base sepolia; source will be ethereum sepolia
            uint256 sourceChainID = 11155111;
            sourceChainSelector = getChainSelector(sourceChainID);
            sourceChainProtocolSmartWallet = getProtocolSmartWallet(sourceChainID);
        }

        if (protocolSmartWallet == address(0) || sourceChainProtocolSmartWallet == address(0) || sourceChainSelector == 0) {
            revert("Please configure deployment for this chainid");
        }

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

        vm.stopBroadcast();
    }

    function getProtocolSmartWallet(uint256 chainid) internal pure returns (address) {
        if (chainid == 11155111) {
            return 0x8D51aAE7F52B74D26D4a38803c0328dfD5C87af0;
        } else if (chainid == 84532) {
            return 0xe698BA634C9266962e19f3D21667e8036202CA33;
        }
        revert("chainid not configured");
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
