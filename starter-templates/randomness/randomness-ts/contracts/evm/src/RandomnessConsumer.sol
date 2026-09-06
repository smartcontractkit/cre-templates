// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/**
 * @title RandomnessConsumer
 * @notice VRF-style request/fulfill randomness backed by a CRE workflow.
 *
 *         Flow (mirrors Chainlink VRF's requestRandomWords/fulfillRandomWords):
 *           1. A caller invokes requestRandomness(), which emits RandomnessRequested.
 *           2. A CRE workflow watches that event via a LogTrigger, generates a
 *              consensus-safe random value, and writes a signed report back here.
 *           3. ReceiverTemplate verifies the report came from the CRE Forwarder,
 *              then _processReport() records the random word for that requestId.
 *
 *         The randomness is produced inside a CRE workflow (agreed by the DON)
 *         instead of by the VRF coordinator. Customize the reveal logic in
 *         _processReport() for your own use case (loot boxes, raffles, sampling...).
 */
contract RandomnessConsumer is ReceiverTemplate {
    /// @notice Monotonic counter used to mint request IDs.
    uint256 public nextRequestId;

    /// @notice requestId => is awaiting fulfillment.
    mapping(uint256 => bool) public pendingRequests;

    /// @notice requestId => fulfilled random word (0 until fulfilled).
    mapping(uint256 => uint256) public randomWords;

    event RandomnessRequested(uint256 indexed requestId, address indexed requester);
    event RandomnessFulfilled(uint256 indexed requestId, uint256 randomWord);

    constructor(address forwarder) ReceiverTemplate(forwarder) {}

    /// @notice Request a random number. Emits an event the CRE workflow reacts to.
    /// @dev Commit any user choices tied to this draw BEFORE calling, and reveal
    ///      only after fulfillment — never in the same transaction.
    /// @return requestId Identifier used to correlate the later fulfillment.
    function requestRandomness() external returns (uint256 requestId) {
        requestId = ++nextRequestId;
        pendingRequests[requestId] = true;
        emit RandomnessRequested(requestId, msg.sender);
    }

    /// @notice Called by the CRE Forwarder via ReceiverTemplate after metadata validation.
    /// @param report ABI-encoded (uint256 requestId, uint256 randomWord)
    function _processReport(bytes calldata report) internal override {
        (uint256 requestId, uint256 randomWord) = abi.decode(report, (uint256, uint256));

        // Only fulfill a known, still-pending request. Mark it fulfilled BEFORE any
        // downstream logic to eliminate re-entrancy and replay/reroll windows.
        require(pendingRequests[requestId], "unknown or already-fulfilled request");
        delete pendingRequests[requestId];

        randomWords[requestId] = randomWord;
        emit RandomnessFulfilled(requestId, randomWord);
    }
}
