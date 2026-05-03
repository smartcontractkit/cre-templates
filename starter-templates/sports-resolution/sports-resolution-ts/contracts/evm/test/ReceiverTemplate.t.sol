// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "../src/ReceiverTemplate.sol";

interface Vm {
    function prank(address msgSender) external;
    function expectRevert(bytes4 revertData) external;
    function expectRevert(bytes calldata revertData) external;
}

contract MockReceiver is ReceiverTemplate {
    uint256 public reportCount;
    bytes32 public lastReportHash;

    constructor(address forwarder) ReceiverTemplate(forwarder) {}

    function _processReport(bytes calldata report) internal override {
        reportCount++;
        lastReportHash = keccak256(report);
    }
}

contract ReceiverTemplateTest {
    Vm private constant vm = Vm(address(uint160(uint256(keccak256("hevm cheat code")))));

    address private constant FORWARDER = address(uint160(1));
    address private constant AUTHOR = address(uint160(2));
    address private constant ATTACKER = address(uint160(3));
    bytes32 private constant WORKFLOW_ID = bytes32(uint256(0x1234));

    function testConstructorRejectsZeroForwarder() external {
        vm.expectRevert(ReceiverTemplate.InvalidForwarderAddress.selector);
        new MockReceiver(address(0));
    }

    function testOnlyConfiguredForwarderCanSubmitReports() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);

        vm.expectRevert(abi.encodeWithSelector(ReceiverTemplate.InvalidSender.selector, ATTACKER, FORWARDER));
        vm.prank(ATTACKER);
        receiver.onReport("", "report");
    }

    function testOwnerCannotDisableForwarderValidation() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);

        vm.expectRevert(ReceiverTemplate.InvalidForwarderAddress.selector);
        receiver.setForwarderAddress(address(0));
    }

    function testWorkflowNameRequiresAuthorBeforeConfiguration() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);

        vm.expectRevert(ReceiverTemplate.WorkflowNameRequiresAuthorValidation.selector);
        receiver.setExpectedWorkflowName("game-resolution");
    }

    function testCannotClearAuthorWhileWorkflowNameValidationIsEnabled() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);
        receiver.setExpectedAuthor(AUTHOR);
        receiver.setExpectedWorkflowName("game-resolution");

        vm.expectRevert(ReceiverTemplate.WorkflowNameRequiresAuthorValidation.selector);
        receiver.setExpectedAuthor(address(0));
    }

    function testRejectsMalformedMetadataLength() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);
        receiver.setExpectedAuthor(AUTHOR);
        bytes memory malformedMetadata = new bytes(61);

        vm.expectRevert(abi.encodeWithSelector(ReceiverTemplate.InvalidMetadataLength.selector, 61, 62));
        vm.prank(FORWARDER);
        receiver.onReport(malformedMetadata, "report");
    }

    function testAcceptsReportWhenForwarderAndWorkflowMetadataMatch() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);
        receiver.setExpectedAuthor(AUTHOR);
        receiver.setExpectedWorkflowId(WORKFLOW_ID);
        receiver.setExpectedWorkflowName("game-resolution");

        bytes memory metadata = abi.encodePacked(WORKFLOW_ID, _workflowName("game-resolution"), AUTHOR);
        bytes memory report = abi.encode("verified report");

        vm.prank(FORWARDER);
        receiver.onReport(metadata, report);

        _assertEq(receiver.reportCount(), 1);
        _assertEq(receiver.lastReportHash(), keccak256(report));
    }

    function testRejectsWrongWorkflowOwner() external {
        MockReceiver receiver = new MockReceiver(FORWARDER);
        receiver.setExpectedAuthor(AUTHOR);

        bytes memory metadata = abi.encodePacked(WORKFLOW_ID, _workflowName("game-resolution"), ATTACKER);

        vm.expectRevert(abi.encodeWithSelector(ReceiverTemplate.InvalidAuthor.selector, ATTACKER, AUTHOR));
        vm.prank(FORWARDER);
        receiver.onReport(metadata, "report");
    }

    function _workflowName(string memory name) private pure returns (bytes10) {
        bytes32 hash = sha256(bytes(name));
        bytes16 hexChars = "0123456789abcdef";
        bytes memory hexString = new bytes(64);

        for (uint256 i = 0; i < 32; i++) {
            hexString[i * 2] = hexChars[uint8(hash[i] >> 4)];
            hexString[i * 2 + 1] = hexChars[uint8(hash[i] & 0x0f)];
        }

        bytes memory first10 = new bytes(10);
        for (uint256 i = 0; i < 10; i++) {
            first10[i] = hexString[i];
        }

        return bytes10(first10);
    }

    function _assertEq(uint256 actual, uint256 expected) private pure {
        if (actual != expected) {
            revert("uint mismatch");
        }
    }

    function _assertEq(bytes32 actual, bytes32 expected) private pure {
        if (actual != expected) {
            revert("bytes32 mismatch");
        }
    }
}
