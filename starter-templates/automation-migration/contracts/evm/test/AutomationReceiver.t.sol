// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {AutomationReceiver} from "../src/AutomationReceiver.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

interface Vm {
    function prank(address msgSender) external;
    function expectRevert(bytes4 revertData) external;
    function expectRevert(bytes calldata revertData) external;
}

/// @dev Minimal Automation-style target. `performUpkeep` records the last performData and
///      can be toggled to revert, to exercise the execution-failure path.
contract MockUpkeep {
    bool public shouldRevert;
    uint256 public performCount;
    bytes public lastPerformData;

    function setShouldRevert(bool value) external {
        shouldRevert = value;
    }

    function performUpkeep(bytes calldata performData) external {
        if (shouldRevert) {
            revert("upkeep failed");
        }
        performCount++;
        lastPerformData = performData;
    }
}

contract AutomationReceiverTest {
    Vm private constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    address private constant FORWARDER = address(uint160(1));
    address private constant ATTACKER = address(uint160(3));
    bytes4 private constant PERFORM_SELECTOR = bytes4(keccak256("performUpkeep(bytes)"));

    AutomationReceiver private receiver;
    MockUpkeep private target;

    constructor() {
        receiver = new AutomationReceiver(FORWARDER);
        target = new MockUpkeep();
    }

    // ─── helpers ────────────────────────────────────────────────
    function _performCall(bytes memory performData) private pure returns (bytes memory) {
        return abi.encodeWithSignature("performUpkeep(bytes)", performData);
    }

    function _report(address tgt, bytes memory callData) private pure returns (bytes memory) {
        return abi.encode(tgt, callData);
    }

    function _deliver(bytes memory report) private {
        vm.prank(FORWARDER);
        receiver.onReport("", report);
    }

    // ─── inbound auth (delegated to ReceiverTemplate) ───────────
    function testOnlyForwarderCanDeliver() external {
        bytes memory report = _report(address(target), _performCall(hex"01"));

        vm.expectRevert(abi.encodeWithSelector(_invalidSenderSelector(), ATTACKER, FORWARDER));
        vm.prank(ATTACKER);
        receiver.onReport("", report);
    }

    // ─── outbound allowlist ─────────────────────────────────────
    function testUnauthorizedTargetSelectorReverts() external {
        // Not allowlisted → must revert loudly.
        bytes memory report = _report(address(target), _performCall(hex"01"));

        vm.expectRevert(
            abi.encodeWithSelector(AutomationReceiver.CallNotAllowed.selector, address(target), PERFORM_SELECTOR)
        );
        _deliver(report);
    }

    function testAllowedCallExecutes() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        bytes memory performData = hex"deadbeef";
        _deliver(_report(address(target), _performCall(performData)));

        _assertEq(target.performCount(), 1);
        _assertEq(keccak256(target.lastPerformData()), keccak256(performData));
    }

    function testRevokedCallReverts() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, false);

        vm.expectRevert(
            abi.encodeWithSelector(AutomationReceiver.CallNotAllowed.selector, address(target), PERFORM_SELECTOR)
        );
        _deliver(_report(address(target), _performCall(hex"01")));
    }

    // ─── execution failure is swallowed (Automation parity) ─────
    function testAllowedButFailingCallDoesNotRevert() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        target.setShouldRevert(true);

        // Must NOT revert: the report is consumed and CallFailed is emitted.
        _deliver(_report(address(target), _performCall(hex"01")));

        // The reverting upkeep changed no state.
        _assertEq(target.performCount(), 0);
    }

    // ─── malformed reports ──────────────────────────────────────
    function testZeroTargetReverts() external {
        vm.expectRevert(AutomationReceiver.InvalidTargetAddress.selector);
        _deliver(_report(address(0), _performCall(hex"01")));
    }

    function testMissingSelectorReverts() external {
        // 3 bytes of calldata → no full selector.
        vm.expectRevert(AutomationReceiver.MissingSelector.selector);
        _deliver(_report(address(target), hex"010203"));
    }

    // ─── allowlist administration ───────────────────────────────
    function testSetCallAllowedIsOwnerOnly() external {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, ATTACKER));
        vm.prank(ATTACKER);
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
    }

    function testSetCallAllowedRejectsZeroTarget() external {
        vm.expectRevert(AutomationReceiver.InvalidTargetAddress.selector);
        receiver.setCallAllowed(address(0), PERFORM_SELECTOR, true);
    }

    function testIsCallAllowedReflectsState() external {
        _assertFalse(receiver.isCallAllowed(address(target), PERFORM_SELECTOR));
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        _assertTrue(receiver.isCallAllowed(address(target), PERFORM_SELECTOR));
    }

    // ─── tiny assertion helpers (no forge-std dependency) ───────
    function _invalidSenderSelector() private pure returns (bytes4) {
        // ReceiverTemplate.InvalidSender(address,address)
        return bytes4(keccak256("InvalidSender(address,address)"));
    }

    function _assertEq(uint256 actual, uint256 expected) private pure {
        if (actual != expected) revert("uint mismatch");
    }

    function _assertEq(bytes32 actual, bytes32 expected) private pure {
        if (actual != expected) revert("bytes32 mismatch");
    }

    function _assertTrue(bool value) private pure {
        if (!value) revert("expected true");
    }

    function _assertFalse(bool value) private pure {
        if (value) revert("expected false");
    }
}
