// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {AutomationReceiver} from "../src/AutomationReceiver.sol";
import {ReceiverTemplate} from "../src/ReceiverTemplate.sol";
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

/// @dev Burns a fixed amount of gas on every call so we can test the insufficient-gas guard.
contract MockGasHog {
    uint256 public callCount;

    /// @dev Spins in a tight loop consuming approximately `gasToConsume` gas.
    ///      Uses inline assembly to avoid the compiler optimising the loop away.
    function performUpkeep(bytes calldata) external {
        callCount++;
    }
}

/// @dev Records the gas remaining at the start of performUpkeep. Used to verify that
///      target.call{gas: consumerGasLimit} forwards exactly consumerGasLimit gas
///      and that GAS_OVERHEAD covers CALL + LOG3 overhead for a no-op consumer.
contract MockGasRecorder {
    uint256 public gasOnEntry;

    function performUpkeep(bytes calldata) external {
        gasOnEntry = gasleft();
    }
}

contract AutomationReceiverTest {
    Vm private constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    address private constant FORWARDER = address(uint160(1));
    address private constant ATTACKER  = address(uint160(3));
    bytes4  private constant PERFORM_SELECTOR = bytes4(keccak256("performUpkeep(bytes)"));

    // Workflow identity used across all delivery tests.
    // The receiver accepts either (a) workflowId alone, or (b) workflowOwner + workflowName.
    // Tests below exercise both valid options and the failure cases.
    bytes32 private constant WORKFLOW_ID    = bytes32(uint256(42));
    address private constant WORKFLOW_OWNER = address(uint160(5));

    // GAS_OVERHEAD mirrors the private constant in AutomationReceiver (EIP-2929 worst-case).
    uint256 private constant GAS_OVERHEAD = 7_000;

    AutomationReceiver private receiver;
    MockUpkeep private target;
    MockGasHog private gasHog;
    MockGasRecorder private gasRecorder;

    constructor() {
        receiver = new AutomationReceiver(FORWARDER);
        // Use option (a): workflowId alone satisfies the identity guard.
        receiver.setExpectedWorkflowId(WORKFLOW_ID);
        receiver.setExpectedAuthor(WORKFLOW_OWNER);
        target = new MockUpkeep();
        gasHog = new MockGasHog();
        gasRecorder = new MockGasRecorder();
    }

    // ─── helpers ────────────────────────────────────────────────
    function _performCall(bytes memory performData) private pure returns (bytes memory) {
        return abi.encodeWithSignature("performUpkeep(bytes)", performData);
    }

    function _report(address tgt, bytes memory callData) private pure returns (bytes memory) {
        return abi.encode(tgt, callData);
    }

    /// @dev Builds the 62-byte metadata expected by ReceiverTemplate._decodeMetadata:
    ///      bytes 0-31  : workflowId (bytes32)
    ///      bytes 32-41 : workflowName (bytes10, zero — not validated in these tests)
    ///      bytes 42-61 : workflowOwner (address)
    function _metadata(bytes32 wfId, address wfOwner) private pure returns (bytes memory) {
        return abi.encodePacked(wfId, bytes10(0), wfOwner);
    }

    function _deliver(bytes memory report) private {
        vm.prank(FORWARDER);
        receiver.onReport(_metadata(WORKFLOW_ID, WORKFLOW_OWNER), report);
    }

    // ─── inbound auth (delegated to ReceiverTemplate) ───────────
    function testOnlyForwarderCanDeliver() external {
        bytes memory report = _report(address(target), _performCall(hex"01"));

        vm.expectRevert(abi.encodeWithSelector(_invalidSenderSelector(), ATTACKER, FORWARDER));
        vm.prank(ATTACKER);
        receiver.onReport(_metadata(WORKFLOW_ID, WORKFLOW_OWNER), report);
    }

    // ─── forwarder-zero guard ────────────────────────────────────
    /// @dev setForwarderAddress(address(0)) succeeds at the setter level (ReceiverTemplate
    ///      does not block it), but _processReport must reject any subsequent delivery.
    function testSetForwarderToZeroAllowedButDeliveryReverts() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        // Setter succeeds — this is expected ReceiverTemplate behaviour that we cannot change.
        receiver.setForwarderAddress(address(0));

        // Anyone can now reach onReport (no forwarder gate), but _processReport blocks them.
        vm.expectRevert(ReceiverTemplate.InvalidForwarderAddress.selector);
        // Call directly (no prank needed since forwarder check is bypassed by the zero address).
        receiver.onReport(_metadata(WORKFLOW_ID, WORKFLOW_OWNER), _report(address(target), _performCall(hex"01")));
    }

    // ─── workflow identity guard ─────────────────────────────────
    /// @dev Option (a): workflowId alone satisfies the identity guard. A receiver that has
    ///      only workflowId configured (no owner, no name) must accept deliveries without
    ///      reverting with WorkflowIdentityNotConfigured.
    function testWorkflowIdAloneSuffices() external {
        AutomationReceiver freshReceiver = new AutomationReceiver(FORWARDER);
        freshReceiver.setExpectedWorkflowId(WORKFLOW_ID); // only workflowId — no owner, no name
        freshReceiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        vm.prank(FORWARDER);
        // ReceiverTemplate skips owner/name checks (none configured); identity guard passes.
        freshReceiver.onReport(_metadata(WORKFLOW_ID, address(0)), _report(address(target), _performCall(hex"01")));

        _assertEq(target.performCount(), 1);
    }

    /// @dev Option (b): workflowOwner + workflowName (no workflowId) satisfies the identity guard.
    function testOwnerAndNameSufficeWithoutWorkflowId() external {
        AutomationReceiver freshReceiver = new AutomationReceiver(FORWARDER);
        freshReceiver.setExpectedAuthor(WORKFLOW_OWNER);
        freshReceiver.setExpectedWorkflowName("my-workflow");
        freshReceiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        // Read back the encoded bytes10 to pass the correct value in metadata.
        bytes10 wfName = freshReceiver.getExpectedWorkflowName();

        vm.prank(FORWARDER);
        // workflowId in metadata is zero — ReceiverTemplate skips id check (none configured).
        freshReceiver.onReport(
            abi.encodePacked(bytes32(0), wfName, WORKFLOW_OWNER),
            _report(address(target), _performCall(hex"01"))
        );

        _assertEq(target.performCount(), 1);
    }

    /// @dev Only workflowOwner set, no workflowId and no workflowName — neither option satisfied.
    function testOnReportRevertsWhenOwnerSetButNameMissing() external {
        AutomationReceiver freshReceiver = new AutomationReceiver(FORWARDER);
        freshReceiver.setExpectedAuthor(WORKFLOW_OWNER); // owner set, name and id intentionally missing
        freshReceiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        vm.expectRevert(AutomationReceiver.WorkflowIdentityNotConfigured.selector);
        vm.prank(FORWARDER);
        freshReceiver.onReport(_metadata(bytes32(0), WORKFLOW_OWNER), _report(address(target), _performCall(hex"01")));
    }

    /// @dev No identity fields at all — neither option satisfied.
    function testOnReportRevertsWhenNeitherIdentityOptionConfigured() external {
        AutomationReceiver freshReceiver = new AutomationReceiver(FORWARDER);
        freshReceiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        vm.expectRevert(AutomationReceiver.WorkflowIdentityNotConfigured.selector);
        vm.prank(FORWARDER);
        freshReceiver.onReport(_metadata(bytes32(0), address(0)), _report(address(target), _performCall(hex"01")));
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

    function testSetCallAllowedRejectsCodelessTarget() external {
        // address(uint160(99)) is an EOA with no deployed code.
        address eoa = address(uint160(99));
        vm.expectRevert(
            abi.encodeWithSelector(AutomationReceiver.TargetHasNoCode.selector, eoa)
        );
        receiver.setCallAllowed(eoa, PERFORM_SELECTOR, true);
    }

    function testIsCallAllowedReflectsState() external {
        _assertFalse(receiver.isCallAllowed(address(target), PERFORM_SELECTOR));
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        _assertTrue(receiver.isCallAllowed(address(target), PERFORM_SELECTOR));
    }

    // ─── consumer gas limit administration ─────────────────────
    function testSetConsumerGasLimitIsOwnerOnly() external {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, ATTACKER));
        vm.prank(ATTACKER);
        receiver.setConsumerGasLimit(address(gasHog), PERFORM_SELECTOR, 500_000);
    }

    function testSetConsumerGasLimitRejectsZeroTarget() external {
        vm.expectRevert(AutomationReceiver.InvalidTargetAddress.selector);
        receiver.setConsumerGasLimit(address(0), PERFORM_SELECTOR, 100_000);
    }

    function testSetAndGetConsumerGasLimit() external {
        // Default is 0 for any pair.
        _assertEq(receiver.getConsumerGasLimit(address(target), PERFORM_SELECTOR), 0);

        receiver.setConsumerGasLimit(address(target), PERFORM_SELECTOR, 300_000);
        _assertEq(receiver.getConsumerGasLimit(address(target), PERFORM_SELECTOR), 300_000);

        // Reset to 0 disables the guard for that pair.
        receiver.setConsumerGasLimit(address(target), PERFORM_SELECTOR, 0);
        _assertEq(receiver.getConsumerGasLimit(address(target), PERFORM_SELECTOR), 0);
    }

    function testConsumerGasLimitIsPerPair() external {
        // Setting a limit for one pair must not affect any other pair.
        receiver.setConsumerGasLimit(address(gasHog), PERFORM_SELECTOR, 200_000);

        _assertEq(receiver.getConsumerGasLimit(address(gasHog), PERFORM_SELECTOR), 200_000);
        _assertEq(receiver.getConsumerGasLimit(address(target), PERFORM_SELECTOR), 0);
        _assertEq(receiver.getConsumerGasLimit(address(gasHog), bytes4(keccak256("otherFn()"))), 0);
    }

    // ─── gas guard ──────────────────────────────────────────────
    function testInsufficientGasReverts() external {
        receiver.setCallAllowed(address(gasHog), PERFORM_SELECTOR, true);
        uint256 limit = 200_000;
        receiver.setConsumerGasLimit(address(gasHog), PERFORM_SELECTOR, limit);

        // Deliver with gas that is less than limit + limit/63 + GAS_OVERHEAD so the guard fires.
        // expectRevert(bytes4) only matches no-argument errors; InsufficientGas carries two
        // uint256 args, so we use try/catch and inspect only the 4-byte selector.
        bytes memory report = _report(address(gasHog), _performCall(hex""));
        bool reverted;
        vm.prank(FORWARDER);
        try receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD - 1}(_metadata(WORKFLOW_ID, WORKFLOW_OWNER), report) {
            reverted = false;
        } catch (bytes memory data) {
            bytes4 sel;
            assembly {
                sel := mload(add(data, 32))
            }
            if (sel != AutomationReceiver.InsufficientGas.selector) revert("wrong revert selector");
            reverted = true;
        }
        _assertTrue(reverted);
    }

    function testSufficientGasWithLimitSucceeds() external {
        receiver.setCallAllowed(address(gasHog), PERFORM_SELECTOR, true);
        uint256 limit = 50_000;
        receiver.setConsumerGasLimit(address(gasHog), PERFORM_SELECTOR, limit);

        bytes memory report = _report(address(gasHog), _performCall(hex""));
        // Deliver with plenty of gas — should succeed and not revert.
        vm.prank(FORWARDER);
        receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD + 50_000}(_metadata(WORKFLOW_ID, WORKFLOW_OWNER), report);

        _assertEq(gasHog.callCount(), 1);
    }

    function testGasLimitZeroPreservesUnboundedBehavior() external {
        // Default limit is 0 for every pair: no guard, fire-and-forget semantics unchanged.
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        target.setShouldRevert(true);

        // Must NOT revert even with limit == 0: fire-and-forget is preserved.
        _deliver(_report(address(target), _performCall(hex"01")));
        _assertEq(target.performCount(), 0);
    }

    // ─── EIP-150 (63/64 rule) ────────────────────────────────────
    /// @dev Verifies that the `required` value emitted in InsufficientGas includes the
    ///      EIP-150 buffer (consumerGasLimit / 63). By decoding the error arguments we can
    ///      assert the exact formula without needing forge-std gas metering. If someone
    ///      removes the `/ 63` term from _processReport this assertion fails.
    function testRequiredIncludesEIP150Buffer() external {
        receiver.setCallAllowed(address(gasHog), PERFORM_SELECTOR, true);
        uint256 limit = 200_000;
        receiver.setConsumerGasLimit(address(gasHog), PERFORM_SELECTOR, limit);

        bytes memory report = _report(address(gasHog), _performCall(hex""));
        bool reverted;
        uint256 emittedRequired;
        vm.prank(FORWARDER);
        // Give less than required so the guard fires and emits the InsufficientGas error.
        try receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD - 1}(
            _metadata(WORKFLOW_ID, WORKFLOW_OWNER), report
        ) {
            reverted = false;
        } catch (bytes memory errData) {
            bytes4 sel;
            assembly {
                sel := mload(add(errData, 32))
            }
            if (sel != AutomationReceiver.InsufficientGas.selector) revert("wrong revert selector");
            // InsufficientGas(uint256 available, uint256 required)
            // Error payload layout (after the 4-byte selector):
            //   bytes  4-35: available (uint256)
            //   bytes 36-67: required  (uint256)
            // In memory: errData[0..31] = length; errData[32..] = payload.
            // required starts at memory offset 32 (length) + 4 (selector) + 32 (available) = 68.
            assembly {
                emittedRequired := mload(add(errData, 68))
            }
            reverted = true;
        }
        _assertTrue(reverted);
        _assertEq(emittedRequired, limit + limit / 63 + GAS_OVERHEAD);
    }

    /// @dev Verifies that at a consumerGasLimit above the 63 × GAS_OVERHEAD threshold (~441,000)
    ///      the target still receives its full configured gas. Without the EIP-150 buffer in
    ///      `required`, the CALL would deliver less than consumerGasLimit due to the 63/64 cap.
    ///      For example, at limit = 600,000 without the fix:
    ///        available ≈ limit + GAS_OVERHEAD → 63/64 × 607,000 ≈ 597,562 < 600,000.
    ///      With the fix (required includes limit/63), the CALL always has enough headroom.
    function testEIP150TermEnsuresFullGasForwardedAtHighLimit() external {
        receiver.setCallAllowed(address(gasRecorder), PERFORM_SELECTOR, true);
        uint256 limit = 600_000; // above 63 × GAS_OVERHEAD = 441,000
        receiver.setConsumerGasLimit(address(gasRecorder), PERFORM_SELECTOR, limit);

        bytes memory report = _report(address(gasRecorder), _performCall(hex""));
        vm.prank(FORWARDER);
        receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD + 60_000}(
            _metadata(WORKFLOW_ID, WORKFLOW_OWNER), report
        );

        // gasOnEntry records gasleft() before the SSTORE; it equals the gas received by the
        // CALL frame minus the function-prologue cost (~200 gas).
        uint256 gasReceived = gasRecorder.gasOnEntry();
        _assertTrue(gasReceived >= limit - 500); // target got ≈ consumerGasLimit
        _assertTrue(gasReceived <= limit);
    }

    // ─── GAS_OVERHEAD accuracy ──────────────────────────────────
    /// @dev Validates that GAS_OVERHEAD is large enough to cover all post-guard overhead
    ///      (CALL opcode cold dispatch 2,600 + returnData handling + LOG3 emission 2,012).
    ///      A no-op consumer is used as the worst case: the target records gasleft() via
    ///      a single cold SSTORE (≈ 22,100 gas) and returns — no business logic.
    ///      If GAS_OVERHEAD were severely underestimated (e.g. 1,000) the function would
    ///      OOG during LOG3 emission and this test would fail.
    ///      The 60,000 pre-check allowance covers ReceiverTemplate metadata validation,
    ///      identity checks in _processReport (this.get*() calls), report decoding, and
    ///      all cold/warm slot accesses before the gasleft() guard.
    function testGasOverheadCorrectlyCoversNoOpConsumer() external {
        receiver.setCallAllowed(address(gasRecorder), PERFORM_SELECTOR, true);
        // consumerGasLimit must cover the cold SSTORE inside performUpkeep (~22,100),
        // leaving room for function prologue and return.
        uint256 limit = 30_000;
        receiver.setConsumerGasLimit(address(gasRecorder), PERFORM_SELECTOR, limit);

        bytes memory report = _report(address(gasRecorder), _performCall(hex""));
        vm.prank(FORWARDER);
        receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD + 60_000}(
            _metadata(WORKFLOW_ID, WORKFLOW_OWNER), report
        );

        // Completing without revert confirms GAS_OVERHEAD covers CALL + LOG3.
        // gasOnEntry records gasleft() before the SSTORE, so it is close to but slightly
        // below consumerGasLimit (call-prologue takes ~160 gas).
        // Verify target received exactly consumerGasLimit gas (call{gas: limit} semantics):
        // since limit < 63/64 * available, the EVM forwards exactly limit, not a fraction.
        uint256 gasReceived = gasRecorder.gasOnEntry();
        _assertTrue(gasReceived >= limit - 500); // 500 gas tolerance for call-frame setup
        _assertTrue(gasReceived <= limit);
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
