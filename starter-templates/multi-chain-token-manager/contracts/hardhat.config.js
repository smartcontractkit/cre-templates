require('@nomicfoundation/hardhat-toolbox');
require('@nomicfoundation/hardhat-foundry');
require('solidity-coverage');

module.exports = {
  solidity: {
    version: '0.8.26',
    settings: {
      optimizer: {
        enabled: true,
        runs: 1000000
      },
      viaIR: true,
      evmVersion: 'paris'
    }
  },
  paths: {
    artifacts: './artifacts',
    cache: './cache_hardhat',
    sources: './src',
    tests: './test'  
  }
};
