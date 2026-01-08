// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.26;

import {IPolicyEngine} from "@chainlink/policy-management/interfaces/IPolicyEngine.sol";
import {Policy} from "@chainlink/policy-management/core/Policy.sol";

/**
 * @title AddressBlacklistPolicy
 * @notice A policy that rejects transactions if the target address is blacklisted.
 * @dev This policy checks if an address (extracted from transaction parameters) is on a blacklist.
 * 
 * Use case: Prevent stablecoin minting or transfers to sanctioned/blocked addresses.
 * 
 * Integration:
 * 1. Deploy via ERC1967Proxy
 * 2. Attach to protected function (e.g., _processReport)
 * 3. Specify parameter name (e.g., "beneficiary" or "recipient")
 * 4. Policy will check that address against blacklist
 */
contract AddressBlacklistPolicy is Policy {
    /// @custom:storage-location erc7201:policy-management.AddressBlacklistPolicy
    struct AddressBlacklistPolicyStorage {
        /// @notice Mapping of blacklisted addresses. If true, address is blocked.
        mapping(address account => bool isBlacklisted) blacklist;
    }

    // keccak256(abi.encode(uint256(keccak256("policy-management.AddressBlacklistPolicy")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant AddressBlacklistPolicyStorageLocation =
        0x2d4a6f88cf11456224a5d041b8f08e5bbcd4be4911056bec4149096758242d00;

    function _getAddressBlacklistPolicyStorage() private pure returns (AddressBlacklistPolicyStorage storage $) {
        assembly {
            $.slot := AddressBlacklistPolicyStorageLocation
        }
    }

    event AddressBlacklisted(address indexed account);
    event AddressRemovedFromBlacklist(address indexed account);

    /**
     * @notice Configures the policy with initial blacklisted addresses (optional).
     * @dev This is called automatically during initialization via proxy.
     * @param parameters ABI-encoded array of addresses to blacklist initially.
     *                   Pass empty bytes for no initial blacklist.
     */
    function configure(bytes calldata parameters) internal override onlyInitializing {
        // Allow empty initialization (no addresses blacklisted initially)
        if (parameters.length == 0) {
            return;
        }
        
        // Decode initial blacklist
        address[] memory initialBlacklist = abi.decode(parameters, (address[]));
        AddressBlacklistPolicyStorage storage $ = _getAddressBlacklistPolicyStorage();
        
        for (uint256 i = 0; i < initialBlacklist.length; i++) {
            $.blacklist[initialBlacklist[i]] = true;
            emit AddressBlacklisted(initialBlacklist[i]);
        }
    }

    /**
     * @notice Adds an address to the blacklist.
     * @dev Only the policy owner can blacklist addresses.
     * @param account The address to blacklist.
     */
    function addToBlacklist(address account) public onlyOwner {
        AddressBlacklistPolicyStorage storage $ = _getAddressBlacklistPolicyStorage();
        require(!$.blacklist[account], "AddressBlacklistPolicy: already blacklisted");
        $.blacklist[account] = true;
        emit AddressBlacklisted(account);
    }

    /**
     * @notice Removes an address from the blacklist.
     * @dev Only the policy owner can remove addresses from blacklist.
     * @param account The address to remove from blacklist.
     */
    function removeFromBlacklist(address account) public onlyOwner {
        AddressBlacklistPolicyStorage storage $ = _getAddressBlacklistPolicyStorage();
        require($.blacklist[account], "AddressBlacklistPolicy: not blacklisted");
        $.blacklist[account] = false;
        emit AddressRemovedFromBlacklist(account);
    }

    /**
     * @notice Checks if an address is blacklisted.
     * @param account The address to check.
     * @return isBlacklisted True if the address is blacklisted, false otherwise.
     */
    function isBlacklisted(address account) public view returns (bool) {
        AddressBlacklistPolicyStorage storage $ = _getAddressBlacklistPolicyStorage();
        return $.blacklist[account];
    }

    /**
     * @notice Function called by the PolicyEngine to check if execution is allowed.
     * @dev This policy expects exactly 1 parameter: the address to check.
     *      The parameter name is configured when attaching the policy (e.g., "beneficiary").
     * 
     * @param parameters The extracted parameters from the transaction.
     *                   Expected: [address] - The address to check against blacklist.
     * 
     * @return result PolicyResult.Continue if address is not blacklisted.
     *                Reverts with PolicyRejected if address is blacklisted.
     */
    function run(
        address, /* caller */
        address, /* subject */
        bytes4, /* selector */
        bytes[] calldata parameters,
        bytes calldata /* context */
    )
        public
        view
        override
        returns (IPolicyEngine.PolicyResult)
    {
        // Validate parameters
        require(parameters.length == 1, "AddressBlacklistPolicy: expected 1 parameter");
        
        // Decode the address to check
        address addressToCheck = abi.decode(parameters[0], (address));
        
        // Check blacklist
        AddressBlacklistPolicyStorage storage $ = _getAddressBlacklistPolicyStorage();
        if ($.blacklist[addressToCheck]) {
            revert IPolicyEngine.PolicyRejected("address is blacklisted");
        }
        
        // Allow execution to continue
        return IPolicyEngine.PolicyResult.Continue;
    }
}

