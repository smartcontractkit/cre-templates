/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {IPool, DataTypes} from "./IPool.sol";

import {IERC20} from
    "@openzeppelin/contracts@5.0.2/token/ERC20/IERC20.sol";
import {SafeERC20} from
    "@openzeppelin/contracts@5.0.2/token/ERC20/utils/SafeERC20.sol";
import {OwnerIsCreator} from "@chainlink/contracts/src/v0.8/shared/access/OwnerIsCreator.sol";

contract MockPool is OwnerIsCreator, IPool{
    using SafeERC20 for IERC20;

    event CurrentLiquidityRateUpdated(address indexed asset, uint256 currentLiquidityRate);

    error AmountZero();
    error NoBalance();
    error InsufficientBalance();

    mapping(address => uint128) public currentLiquidityRate;

    // User balances (no interest growth in mock)
    mapping(address => mapping(address => uint256)) public balanceOf; // user => asset => amount

    // Set per-asset fixed APR in RAY, e.g. 0.05 * 1e27 for 5%.
    function setCurrentLiquidityRate(address asset, uint128 rate) external onlyOwner {
        currentLiquidityRate[asset] = rate;
        emit CurrentLiquidityRateUpdated(asset, rate);
    }

    // Implement IPool interface.
    // TODO: mint aTokens
    function supply(address asset, uint256 amount, address onBehalfOf, uint16) external override {
        if (amount == 0) {
            revert AmountZero();
        }
        IERC20(asset).safeTransferFrom(msg.sender, address(this), amount);
        balanceOf[onBehalfOf][asset] += amount;
        emit Supply(asset, msg.sender, onBehalfOf, amount, 0);
    }

    // Implement IPool interface.
    // TODO: burn aTokens
    function withdraw(address asset, uint256 amount, address to) external returns (uint256) {

        uint256 bal = balanceOf[msg.sender][asset];
        if (bal == 0) {
            revert NoBalance();
        }
        if (amount == type(uint256).max) {
            amount = bal;
        }
        if (amount > bal) {
            revert InsufficientBalance();
        }

        balanceOf[msg.sender][asset] = bal - amount;
        IERC20(asset).safeTransfer(to, amount);

        emit Withdraw(asset, msg.sender, to, amount);
        return amount;
    }

    function getReserveData(address asset) external view returns (DataTypes.ReserveDataLegacy memory) {
        DataTypes.ReserveDataLegacy memory data;
        data.currentLiquidityRate = currentLiquidityRate[asset];
        return data;
    }
}
