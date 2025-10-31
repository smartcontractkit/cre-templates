/**
 * THIS IS AN EXAMPLE CONTRACT THAT USES HARDCODED VALUES FOR CLARITY.
 * THIS IS AN EXAMPLE CONTRACT THAT USES UN-AUDITED CODE.
 * DO NOT USE THIS CODE IN PRODUCTION.
 */

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

// Minimal definition of AAVE's DataTypes library. See original definition
// for comments and details https://github.com/aave-dao/aave-v3-origin/blob/070cd23d949d828ad78d098ef481c2049f5eabb6/src/contracts/protocol/libraries/types/DataTypes.sol#L4.
library DataTypes {
  struct ReserveDataLegacy {
    ReserveConfigurationMap configuration;
    uint128 liquidityIndex;
    uint128 currentLiquidityRate;
    uint128 variableBorrowIndex;
    uint128 currentVariableBorrowRate;
    uint128 currentStableBorrowRate;
    uint40 lastUpdateTimestamp;
    uint16 id;
    address aTokenAddress;
    address stableDebtTokenAddress;
    address variableDebtTokenAddress;
    address interestRateStrategyAddress;
    uint128 accruedToTreasury;
    uint128 unbacked;
    uint128 isolationModeTotalDebt;
  }

  struct ReserveConfigurationMap {
    uint256 data;
  }
}

// Minimal definition of AAVE's IPool interface. See original definition
// for comments and details: https://github.com/aave-dao/aave-v3-origin/blob/070cd23d949d828ad78d098ef481c2049f5eabb6/src/contracts/interfaces/IPool.sol.
interface IPool {
  event Supply(
    address indexed reserve,
    address user,
    address indexed onBehalfOf,
    uint256 amount,
    uint16 indexed referralCode
  );

  event Withdraw(address indexed reserve, address indexed user, address indexed to, uint256 amount);

  function getReserveData(address asset) external view returns (DataTypes.ReserveDataLegacy memory);
  
  function supply(address asset, uint256 amount, address onBehalfOf, uint16 referralCode) external;

  function withdraw(address asset, uint256 amount, address to) external returns (uint256);
}