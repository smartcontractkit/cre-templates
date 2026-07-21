// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title MonitoredToken
 * @notice A simple contract that emits LargeTransfer events when transfers
 *         exceed a threshold. CRE watches these events via LogTrigger.
 *
 *         In production, this would be your actual token or protocol contract.
 */
contract MonitoredToken {
    event LargeTransfer(
        address indexed from,
        address indexed to,
        uint256 amount,
        uint256 timestamp
    );

    uint256 public largeTransferThreshold = 1000 ether;

    /// @notice Simulate a large transfer that triggers CRE event monitoring
    function simulateLargeTransfer(address to, uint256 amount) external {
        if (amount >= largeTransferThreshold) {
            emit LargeTransfer(msg.sender, to, amount, block.timestamp);
        }
    }
}
