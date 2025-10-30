
/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

library Workflow {

    function keystoneForwarder(uint256 chainid) internal pure returns (address) {
        if (chainid == 11155111) {
            // ethereum sepolia
            return 0xA73aFb610aC79A729A393D679029fCf9618854AD;
        } else if (chainid == 84532) {
            // base sepolia
            return 0xf6914c461AC0943Cfa908C64DCdF84355257c9D5;
        } else {
            revert("keystone forwarder address not set for this chain");
        }
    }

    function mockKeystoneForwarder(uint256 chainid) internal pure returns (address) {
        if (chainid == 11155111) {
            // ethereum sepolia
            return 0x15fC6ae953E024d975e77382eEeC56A9101f9F88;
        } else if (chainid == 84532) {
            // base sepolia
            return 0x82300bd7c3958625581cc2F77bC6464dcEcDF3e5;
        } else {
            revert("mock keystone forwarder address not set for this chain");
        }
    }

    // Helper: convert raw workflow name to on-chain bytes10 representation.
    function toOnChainWorkflowName(string memory rawName) internal pure returns (bytes10) {
        return hashTruncateName(rawName);
    }

    // Helper: convert raw workflow name to on-chain bytes10 representation when
    // using `cre workflow simulate`.
    function toOnChainSimWorkflowName(string memory rawName) internal pure returns (bytes10) {
        // Convert to bytes
        bytes memory raw = bytes(rawName);

        // Pad or truncate to exactly 10 bytes
        bytes memory padded = new bytes(10);
        for (uint256 i = 0; i < raw.length && i < 10; i++) {
            padded[i] = raw[i];
        }

        // Encode to hex string
        string memory encodedWorkflowName = toHexString(padded);

        return hashTruncateName(encodedWorkflowName);
    }

    function hashTruncateName(string memory name) internal pure returns (bytes10) {
        // Hash the hex-encoded string using SHA-256
        bytes32 hash = sha256(bytes(name));

        // Convert hash to hex string and truncate to 10 bytes (20 hex chars)
        string memory hashHex = toHexString(abi.encodePacked(hash));
        bytes memory hashBytes = bytes(hashHex);

        bytes memory truncated = new bytes(10);
        for (uint256 i = 0; i < 10; i++) {
            truncated[i] = hashBytes[i];
        }

        bytes10 truncated10;
        assembly {
            truncated10 := mload(add(truncated, 32))
        }

        return truncated10;
    }

    // Helper: bytes -> hex string
    function toHexString(bytes memory data) internal pure returns (string memory) {
        bytes memory hexChars = "0123456789abcdef";
        bytes memory str = new bytes(data.length * 2);
        for (uint256 i = 0; i < data.length; i++) {
            str[2*i] = hexChars[uint8(data[i] >> 4)];
            str[2*i+1] = hexChars[uint8(data[i] & 0x0f)];
        }
        return string(str);
    }
}
