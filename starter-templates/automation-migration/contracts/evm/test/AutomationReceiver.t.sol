// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {AutomationReceiver} from "../src/AutomationReceiver.sol";
import {ReceiverTemplate} from "../src/ReceiverTemplate.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";

interface Vm {
    function prank(address msgSender) external;
    function expectRevert(bytes4 revertData) external;
    function expectRevert(bytes calldata revertData) external;
    /// @dev Marks `target` cold, as if it had never been accessed in the current transaction.
    ///      Used to simulate production delivery, where setCallAllowed ran in an earlier
    ///      transaction and the target account is therefore cold at delivery time — unlike a
    ///      same-transaction Foundry test, which would otherwise leave it warm.
    function cool(address target) external;
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

/// @dev Worst case for GAS_OVERHEAD: consumes every unit of forwarded gas and reverts (OOG).
///      After the failing CALL the receiver must still have enough gas left to emit the
///      CallFailed LOG3 event. Used to prove GAS_OVERHEAD covers the post-guard path even
///      when the target burns 100% of consumerGasLimit.
contract MockGasBurner {
    function performUpkeep(bytes calldata) external pure {
        // Infinite loop: burns all forwarded gas, then the frame reverts with out-of-gas.
        while (true) {}
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
    uint256 private constant GAS_OVERHEAD = 7000;

    AutomationReceiver private receiver;
    MockUpkeep private target;
    MockGasHog private gasHog;
    MockGasRecorder private gasRecorder;
    MockGasBurner private gasBurner;

    constructor() {
        receiver = new AutomationReceiver(FORWARDER);
        // Use option (a): workflowId alone satisfies the identity guard. The author is also
        // set so ReceiverTemplate additionally validates it against the delivered metadata.
        receiver.setExpectedWorkflowId(WORKFLOW_ID);
        receiver.setExpectedAuthor(WORKFLOW_OWNER);
        target = new MockUpkeep();
        gasHog = new MockGasHog();
        gasRecorder = new MockGasRecorder();
        gasBurner = new MockGasBurner();
    }

    // ─── helpers ────────────────────────────────────────────────
    function _performCall(bytes memory performData) private pure returns (bytes memory) {
        return abi.encodeWithSignature("performUpkeep(bytes)", performData);
    }

    /// @dev Builds a report with block number 0 (ignored unless the monotonicity check is enabled
    ///      for the pair). `callData` is the full function call (selector + arguments).
    function _report(address tgt, bytes memory callData) private pure returns (bytes memory) {
        return abi.encode(tgt, uint256(0), callData);
    }

    /// @dev Builds a report with an explicit block number for the monotonicity check.
    function _reportAtBlock(address tgt, uint256 blockNumber, bytes memory callData)
        private
        pure
        returns (bytes memory)
    {
        return abi.encode(tgt, blockNumber, callData);
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

    /// @dev Delivers a report with a bounded gas budget and classifies the outcome:
    ///      0 = onReport returned (success, including the swallowed CallFailed path),
    ///      1 = reverted with InsufficientGas (the gas guard fired),
    ///      2 = any other revert (e.g. the receiver itself ran out of gas post-guard).
    ///      Cools `target` first so every delivery pays the cold-CALL cost that a real,
    ///      separate-transaction delivery would — otherwise the account would stay warm
    ///      across repeated calls within this single test transaction.
    function _deliverOutcome(address target, uint256 gasAmount, bytes memory report) private returns (uint8) {
        vm.cool(target);
        vm.prank(FORWARDER);
        try receiver.onReport{gas: gasAmount}(_metadata(WORKFLOW_ID, WORKFLOW_OWNER), report) {
            return 0;
        } catch (bytes memory data) {
            bytes4 sel;
            if (data.length >= 4) {
                assembly {
                    sel := mload(add(data, 32))
                }
            }
            return sel == AutomationReceiver.InsufficientGas.selector ? 1 : 2;
        }
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

    function testSetCallAllowedRevocationSkipsCodeCheck() external {
        // A codeless target (e.g. a self-destructed contract) must still be revocable.
        address eoa = address(uint160(99));
        receiver.setCallAllowed(eoa, PERFORM_SELECTOR, false);
        _assertFalse(receiver.isCallAllowed(eoa, PERFORM_SELECTOR));
    }

    function testIsCallAllowedReflectsState() external {
        _assertFalse(receiver.isCallAllowed(address(target), PERFORM_SELECTOR));
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        _assertTrue(receiver.isCallAllowed(address(target), PERFORM_SELECTOR));
    }

    // ─── emergency pause (OpenZeppelin Pausable) ────────────────
    function testPauseIsOwnerOnly() external {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, ATTACKER));
        vm.prank(ATTACKER);
        receiver.pause(true);
    }

    function testUnpauseIsOwnerOnly() external {
        receiver.pause(true);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, ATTACKER));
        vm.prank(ATTACKER);
        receiver.unpause();
    }

    /// @dev Retryable mode (pause(true)): an otherwise-valid delivery must revert with EnforcedPause
    ///      and the target must not be called (the report stays unconsumed and is retryable).
    function testPausedRetryableDeliveryReverts() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.pause(true);

        vm.expectRevert(Pausable.EnforcedPause.selector);
        _deliver(_report(address(target), _performCall(hex"01")));

        _assertEq(target.performCount(), 0);
    }

    /// @dev Non-retryable mode (pause(false)): a delivery must NOT revert — the report is consumed
    ///      (dropped) and the target is not called.
    function testPausedNonRetryableDeliveryConsumesReport() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.pause(false);

        // Must not revert: the report is swallowed while paused in drop mode.
        _deliver(_report(address(target), _performCall(hex"01")));

        _assertEq(target.performCount(), 0);
    }

    /// @dev After unpause the same delivery succeeds — demonstrates a retryable-mode paused report
    ///      is retryable rather than permanently consumed.
    function testUnpauseResumesDelivery() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        receiver.pause(true);
        vm.expectRevert(Pausable.EnforcedPause.selector);
        _deliver(_report(address(target), _performCall(hex"01")));

        receiver.unpause();
        _deliver(_report(address(target), _performCall(hex"01")));

        _assertEq(target.performCount(), 1);
    }

    function testPausedReflectsState() external {
        _assertFalse(receiver.paused());
        receiver.pause(true);
        _assertTrue(receiver.paused());
        receiver.unpause();
        _assertFalse(receiver.paused());
    }

    /// @dev The retryable mode chosen at pause time is reflected by retryableWhilePaused().
    function testRetryableWhilePausedReflectsMode() external {
        receiver.pause(true);
        _assertTrue(receiver.retryableWhilePaused());
        receiver.unpause();

        receiver.pause(false);
        _assertFalse(receiver.retryableWhilePaused());
    }

    function testDoublePauseReverts() external {
        receiver.pause(true);
        vm.expectRevert(Pausable.EnforcedPause.selector);
        receiver.pause(true);
    }

    function testUnpauseWhenNotPausedReverts() external {
        vm.expectRevert(Pausable.ExpectedPause.selector);
        receiver.unpause();
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
    ///      For example, at limit = 1,000,000 without the fix:
    ///        available ≈ limit + GAS_OVERHEAD → 63/64 × 1,007,000 ≈ 991,109 < 1,000,000.
    ///      With the fix (required includes limit/63), the CALL always has enough headroom.
    function testEIP150TermEnsuresFullGasForwardedAtHighLimit() external {
        receiver.setCallAllowed(address(gasRecorder), PERFORM_SELECTOR, true);
        uint256 limit = 1_000_000; // above 63 × GAS_OVERHEAD = 441,000
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
    ///      (cold CALL dispatch ~2,600 + returnData handling + LOG3 emission ~2,200).
    ///      A no-op consumer is used as the worst case: the target records gasleft() via
    ///      a single cold SSTORE (≈ 22,100 gas) and returns — no business logic.
    ///      `vm.cool` marks the target cold before delivery, matching production where
    ///      setCallAllowed ran in an earlier transaction (see GAS_OVERHEAD's docstring);
    ///      without it, Foundry would leave the target warm from `setCallAllowed`'s
    ///      `EXTCODESIZE` read above, understating the real cold-CALL cost.
    ///      If GAS_OVERHEAD were severely underestimated (e.g. 1,000) the function would
    ///      OOG during LOG3 emission and this test would fail.
    ///      The 60,000 pre-check allowance covers ReceiverTemplate metadata validation,
    ///      identity checks in _processReport (direct reads of the inherited permission
    ///      fields), report decoding, and all cold/warm slot accesses before the gasleft() guard.
    function testGasOverheadCorrectlyCoversNoOpConsumer() external {
        receiver.setCallAllowed(address(gasRecorder), PERFORM_SELECTOR, true);
        // consumerGasLimit must cover the cold SSTORE inside performUpkeep (~22,100),
        // leaving room for function prologue and return.
        uint256 limit = 30_000;
        receiver.setConsumerGasLimit(address(gasRecorder), PERFORM_SELECTOR, limit);

        bytes memory report = _report(address(gasRecorder), _performCall(hex""));
        vm.cool(address(gasRecorder));
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

    /// @dev Worst-case proof that GAS_OVERHEAD is sufficient. The consumer burns 100% of the
    ///      forwarded gas and reverts (OOG), which is the most expensive post-guard path: the
    ///      receiver must still emit CallFailed (LOG3) out of the GAS_OVERHEAD headroom.
    ///
    ///      A deliberately small consumerGasLimit is the true worst case for a fixed overhead: the
    ///      `consumerGasLimit / 63` term (which doubles as post-call headroom, see the constant's
    ///      docstring) shrinks toward zero, so GAS_OVERHEAD alone must cover the cold CALL dispatch
    ///      plus the LOG3 emission. `_deliverOutcome` cools the target before every delivery so the
    ///      CALL actually pays the cold-account surcharge, matching production where setCallAllowed
    ///      ran in an earlier transaction — without cooling, the target would stay warm across the
    ///      binary search's repeated calls within this one test transaction and understate the real
    ///      worst case (this is what let the previous, too-small GAS_OVERHEAD value pass here).
    ///
    ///      onReport's success is monotonic in the gas budget (more gas never turns a success into a
    ///      failure). We binary-search the smallest budget that succeeds, then assert that ONE unit
    ///      below it the call reverted with InsufficientGas — i.e. the guard was still firing. If
    ///      GAS_OVERHEAD were undersized, there would instead be an out-of-gas "dead zone" just above
    ///      the guard threshold, so the budget one unit below the first success would revert for a
    ///      different reason (outcome 2) and this assertion would fail.
    function testGasOverheadCoversWorstCaseGasBurningConsumer() external {
        receiver.setCallAllowed(address(gasBurner), PERFORM_SELECTOR, true);
        uint256 limit = 1_000;
        receiver.setConsumerGasLimit(address(gasBurner), PERFORM_SELECTOR, limit);
        bytes memory report = _report(address(gasBurner), _performCall(hex""));

        uint256 lo = 0; // no gas: reverts (outcome != 0)
        uint256 hi = 500_000; // ample: onReport runs to completion (outcome 0)
        _assertEq(_deliverOutcome(address(gasBurner), hi, report), 0);

        // Binary-search the smallest budget at which onReport succeeds.
        while (hi - lo > 1) {
            uint256 mid = (lo + hi) / 2;
            if (_deliverOutcome(address(gasBurner), mid, report) == 0) {
                hi = mid;
            } else {
                lo = mid;
            }
        }

        // `lo` is the largest budget that does NOT succeed. It must revert with InsufficientGas
        // (the guard firing), proving there is no out-of-gas dead zone above the guard threshold:
        // the moment the guard is satisfied, GAS_OVERHEAD is enough to finish (CALL + LOG3).
        _assertEq(_deliverOutcome(address(gasBurner), lo, report), 1);
    }

    // ─── block-number monotonicity (staleness protection) ───────
    /// @dev Disabled by default: an older report still executes (no ordering enforced).
    function testBlockNumberCheckDisabledAllowsOutOfOrder() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        _deliver(_reportAtBlock(address(target), 100, _performCall(hex"01")));
        // Older block — still executes because the check is off.
        _deliver(_reportAtBlock(address(target), 50, _performCall(hex"01")));

        _assertEq(target.performCount(), 2);
    }

    /// @dev With an explicit floor: below-floor is skipped, at-floor and above execute, and a
    ///      report older than the last accepted one is skipped. Stale reports never revert.
    function testBlockNumberCheckWithExplicitFloor() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, true, 100);

        // Below floor → skipped (consumed, no revert).
        _deliver(_reportAtBlock(address(target), 99, _performCall(hex"01")));
        _assertEq(target.performCount(), 0);

        // At floor (>=) → accepted.
        _deliver(_reportAtBlock(address(target), 100, _performCall(hex"01")));
        _assertEq(target.performCount(), 1);

        // Higher → accepted, advances the stored block.
        _deliver(_reportAtBlock(address(target), 101, _performCall(hex"01")));
        _assertEq(target.performCount(), 2);

        // Older than last accepted (101) → skipped.
        _deliver(_reportAtBlock(address(target), 100, _performCall(hex"01")));
        _assertEq(target.performCount(), 2);
    }

    /// @dev CLA-parity: a report at the same block as the last accepted one is still accepted (>=).
    function testBlockNumberCheckEqualBlockAccepted() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, true, 100);

        _deliver(_reportAtBlock(address(target), 100, _performCall(hex"01")));
        _deliver(_reportAtBlock(address(target), 100, _performCall(hex"01")));

        _assertEq(target.performCount(), 2);
    }

    /// @dev initialBlockNumber == 0 snapshots the current block.number as the floor.
    function testBlockNumberCheckSnapshotUsesCurrentBlock() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, true, 0);

        (bool enabled, uint256 floor) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertTrue(enabled);
        _assertEq(floor, block.number);

        // Below the snapshot → skipped.
        if (block.number > 0) {
            _deliver(_reportAtBlock(address(target), block.number - 1, _performCall(hex"01")));
            _assertEq(target.performCount(), 0);
        }

        // At the snapshot → accepted.
        _deliver(_reportAtBlock(address(target), block.number, _performCall(hex"01")));
        _assertEq(target.performCount(), 1);
    }

    /// @dev Getter reflects state, the stored block advances on accepted reports, and disabling clears it.
    function testGetBlockNumberCheckReflectsState() external {
        (bool enabled, uint256 last) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertFalse(enabled);
        _assertEq(last, 0);

        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, true, 500);
        (enabled, last) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertTrue(enabled);
        _assertEq(last, 500);

        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        _deliver(_reportAtBlock(address(target), 600, _performCall(hex"01")));
        (, last) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertEq(last, 600);

        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, false, 0);
        (enabled, last) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertFalse(enabled);
        _assertEq(last, 0);
    }

    /// @dev Enabling the check for one pair does not affect another pair.
    function testBlockNumberCheckIsPerPair() external {
        receiver.setBlockNumberCheck(address(gasHog), PERFORM_SELECTOR, true, 300);

        (bool enabledHog, uint256 lastHog) = receiver.getBlockNumberCheck(address(gasHog), PERFORM_SELECTOR);
        _assertTrue(enabledHog);
        _assertEq(lastHog, 300);

        (bool enabledTarget, uint256 lastTarget) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertFalse(enabledTarget);
        _assertEq(lastTarget, 0);
    }

    function testSetBlockNumberCheckIsOwnerOnly() external {
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, ATTACKER));
        vm.prank(ATTACKER);
        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, true, 0);
    }

    function testSetBlockNumberCheckRejectsZeroTarget() external {
        vm.expectRevert(AutomationReceiver.InvalidTargetAddress.selector);
        receiver.setBlockNumberCheck(address(0), PERFORM_SELECTOR, true, 0);
    }

    /// @dev An InsufficientGas revert must roll back the block-number advance so a retry at the
    ///      same block remains valid (see _processReport staleness comment).
    function testInsufficientGasDoesNotAdvanceBlockNumber() external {
        receiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);
        receiver.setBlockNumberCheck(address(target), PERFORM_SELECTOR, true, 100);
        uint256 limit = 200_000;
        receiver.setConsumerGasLimit(address(target), PERFORM_SELECTOR, limit);

        bytes memory report = _reportAtBlock(address(target), 100, _performCall(hex"01"));

        bool reverted;
        vm.prank(FORWARDER);
        try receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD - 1}(
            _metadata(WORKFLOW_ID, WORKFLOW_OWNER), report
        ) {
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

        (, uint256 lastBlock) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertEq(lastBlock, 100);
        _assertEq(target.performCount(), 0);

        vm.prank(FORWARDER);
        receiver.onReport{gas: limit + limit / 63 + GAS_OVERHEAD + 100_000}(
            _metadata(WORKFLOW_ID, WORKFLOW_OWNER), report
        );
        _assertEq(target.performCount(), 1);
        (, lastBlock) = receiver.getBlockNumberCheck(address(target), PERFORM_SELECTOR);
        _assertEq(lastBlock, 100);
    }

    /// @dev All three identity fields may be configured together for the strongest binding.
    function testCombinedIdentityOptionsAccepted() external {
        AutomationReceiver freshReceiver = new AutomationReceiver(FORWARDER);
        freshReceiver.setExpectedWorkflowId(WORKFLOW_ID);
        freshReceiver.setExpectedAuthor(WORKFLOW_OWNER);
        freshReceiver.setExpectedWorkflowName("my-workflow");
        freshReceiver.setCallAllowed(address(target), PERFORM_SELECTOR, true);

        bytes10 wfName = freshReceiver.getExpectedWorkflowName();
        vm.prank(FORWARDER);
        freshReceiver.onReport(
            abi.encodePacked(WORKFLOW_ID, wfName, WORKFLOW_OWNER),
            _report(address(target), _performCall(hex"01"))
        );

        _assertEq(target.performCount(), 1);
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
