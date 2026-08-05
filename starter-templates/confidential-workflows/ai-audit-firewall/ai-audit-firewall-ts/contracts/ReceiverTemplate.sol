// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

interface IERC165 {
  function supportsInterface(bytes4 interfaceId) external view returns (bool);
}

interface IReceiver is IERC165 {
  function onReport(bytes calldata metadata, bytes calldata report) external;
}

abstract contract ReceiverTemplate is IReceiver {
  address private s_forwarderAddress;

  constructor(address _forwarderAddress) {
    require(_forwarderAddress != address(0), "Invalid forwarder address");
    s_forwarderAddress = _forwarderAddress;
  }

  function getForwarderAddress() external view returns (address) {
    return s_forwarderAddress;
  }

  function setForwarderAddress(address _forwarder) external virtual {
    s_forwarderAddress = _forwarder;
  }

  function onReport(bytes calldata, bytes calldata report) external virtual override {
    if (s_forwarderAddress != address(0) && msg.sender != s_forwarderAddress) {
      revert("Invalid sender");
    }

    _processReport(report);
  }

  function _processReport(bytes calldata report) internal virtual;

  function supportsInterface(bytes4 interfaceId) public pure virtual override returns (bool) {
    return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
  }
}
