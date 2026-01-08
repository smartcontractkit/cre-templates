// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title IReceiverTemplate
 * @notice Abstract contract for receiving DON-signed reports from CRE workflows
 * @dev Consumer contracts should inherit from this to implement secure report validation
 */
abstract contract IReceiverTemplate {
    address public immutable EXPECTED_AUTHOR;
    bytes10 public immutable EXPECTED_WORKFLOW_NAME;

    error InvalidAuthor(address received, address expected);
    error InvalidWorkflowName(bytes10 received, bytes10 expected);

    constructor(address expectedAuthor, bytes10 expectedWorkflowName) {
        EXPECTED_AUTHOR = expectedAuthor;
        EXPECTED_WORKFLOW_NAME = expectedWorkflowName;
    }

    /**
     * @notice Receive and process a signed report from the Forwarder
     * @param metadata Encoded metadata containing workflow info
     * @param report Encoded report data from the workflow
     */
    function onReport(bytes calldata metadata, bytes calldata report) external virtual;

    /**
     * @notice Internal function to be implemented by consumers
     * @param report The decoded report data to process
     */
    function _processReport(bytes calldata report) internal virtual;

    /**
     * @dev ERC165 interface support (required by some implementations)
     */
    function supportsInterface(bytes4 interfaceId) public pure virtual returns (bool) {
        return interfaceId == this.onReport.selector;
    }
}

