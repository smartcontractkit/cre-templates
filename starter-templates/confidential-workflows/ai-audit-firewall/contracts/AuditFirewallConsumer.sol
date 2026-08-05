// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {ReceiverTemplate} from "./ReceiverTemplate.sol";

contract AuditFirewallConsumer is ReceiverTemplate {
  enum Verdict {
    Unknown,
    Allow,
    Deny,
    ManualReview
  }

  Verdict public lastVerdict;
  uint8 public lastRiskMask;
  uint64 public lastChainSelector;
  bytes public lastReport;

  event VerdictReceived(Verdict verdict, uint8 riskMask, uint64 chainSelector);

  constructor(address _forwarderAddress) ReceiverTemplate(_forwarderAddress) {}

  function _processReport(bytes calldata report) internal override {
    (uint8 verdictCode, uint8 riskMask, uint64 chainSelector) = abi.decode(report, (uint8, uint8, uint64));

    Verdict verdict = Verdict.Unknown;
    if (verdictCode == 1) {
      verdict = Verdict.Allow;
    } else if (verdictCode == 2) {
      verdict = Verdict.Deny;
    } else if (verdictCode == 3) {
      verdict = Verdict.ManualReview;
    }

    lastVerdict = verdict;
    lastRiskMask = riskMask;
    lastChainSelector = chainSelector;
    lastReport = report;

    emit VerdictReceived(verdict, riskMask, chainSelector);
  }
}
