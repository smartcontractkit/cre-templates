
/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
// SPDX-License-Identifier: MIT
pragma solidity 0.8.26;

import {ConfirmedOwner} from "@chainlink/contracts/src/v0.8/shared/access/ConfirmedOwner.sol";
import {ITypeAndVersion} from "@chainlink/contracts/src/v0.8/shared/interfaces/ITypeAndVersion.sol";
import {ICommonAggregator} from "@chainlink/contracts/src/v0.8/data-feeds/interfaces/ICommonAggregator.sol";
import {IDecimalAggregator} from "@chainlink/contracts/src/v0.8/data-feeds/interfaces/IDecimalAggregator.sol";
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

interface IAggregator is IDecimalAggregator, ICommonAggregator {}

/// @title A trusted proxy for updating where current answers are read from
/// @notice This contract provides a consistent address for the
/// CurrentAnswerInterface but delegates where it reads from to the owner, who is
/// trusted to update it.
contract DecimalAggregatorProxy is AggregatorV3Interface, ITypeAndVersion, ConfirmedOwner {
    string public constant override typeAndVersion =
        "DecimalAggregatorProxy 1.0.0";

    IAggregator private s_currentAggregator;
    IAggregator private s_proposedAggregator;

    event AggregatorProposed(address indexed current, address indexed proposed);
    event AggregatorConfirmed(address indexed previous, address indexed latest);

    error AggregatorNotProposed(address aggregator);

    constructor(
        address aggregatorAddress,
        address owner
    ) ConfirmedOwner(owner) {
        s_currentAggregator = IAggregator(aggregatorAddress);
    }

    function latestAnswer() external view returns (int256) {
        return s_currentAggregator.latestAnswer();
    }

    function latestRound() external view returns (uint256) {
        return s_currentAggregator.latestRound();
    }

    function latestTimestamp() external view returns (uint256 timestamp) {
        return s_currentAggregator.latestTimestamp();
    }

    function getAnswer(uint256 roundId) external view returns (int256) {
        return s_currentAggregator.getAnswer(roundId);
    }

    function getTimestamp(
        uint256 roundId
    ) external view returns (uint256 timestamp) {
        return s_currentAggregator.getTimestamp(roundId);
    }

    function getRoundData(
        uint80 _roundId
    )
        external
        view
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return s_currentAggregator.getRoundData(_roundId);
    }

    function latestRoundData()
        external
        view
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return s_currentAggregator.latestRoundData();
    }

    /// @notice returns the current aggregator address.
    function aggregator() external view returns (address) {
        return address(s_currentAggregator);
    }

    /// @notice represents the number of decimals the aggregator responses represent.
    function decimals() external view returns (uint8) {
        return s_currentAggregator.decimals();
    }

    /// @notice the version number representing the type of aggregator the proxy
    /// points to.
    function version() external view override returns (uint256 aggregatorVersion) {
        return s_currentAggregator.version();
    }

    /// @notice returns the description of the aggregator the proxy points to.
    function description() external view returns (string memory aggregatorDescription) {
        return s_currentAggregator.description();
    }

    /// @notice returns the current proposed aggregator
    function proposedAggregator()
        external
        view
        returns (address proposedAggregatorAddress)
    {
        return address(s_proposedAggregator);
    }

    /// @notice Allows the owner to propose a new address for the aggregator
    /// @param aggregatorAddress The new address for the aggregator contract
    function proposeAggregator(address aggregatorAddress) external onlyOwner {
        s_proposedAggregator = IAggregator(aggregatorAddress);
        emit AggregatorProposed(
            address(s_currentAggregator),
            aggregatorAddress
        );
    }

    /// @notice Allows the owner to confirm and change the address
    /// to the proposed aggregator
    /// @dev Reverts if the given address doesn't match what was previously proposed
    /// @param aggregatorAddress The new address for the aggregator contract
    function confirmAggregator(address aggregatorAddress) external onlyOwner {
        if (aggregatorAddress != address(s_proposedAggregator)) {
            revert AggregatorNotProposed(aggregatorAddress);
        }
        address previousAggregator = address(s_currentAggregator);
        delete s_proposedAggregator;
        s_currentAggregator = IAggregator(aggregatorAddress);
        emit AggregatorConfirmed(previousAggregator, aggregatorAddress);
    }
}
