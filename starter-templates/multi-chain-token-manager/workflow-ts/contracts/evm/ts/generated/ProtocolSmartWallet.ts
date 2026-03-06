// Code generated — DO NOT EDIT.
import {
  decodeEventLog,
  decodeFunctionResult,
  encodeEventTopics,
  encodeFunctionData,
  zeroAddress,
} from 'viem'
import type { Address, Hex } from 'viem'
import {
  bytesToHex,
  encodeCallMsg,
  EVMClient,
  hexToBase64,
  LAST_FINALIZED_BLOCK_NUMBER,
  prepareReportRequest,
  type EVMLog,
  type Runtime,
} from '@chainlink/cre-sdk'

export interface DecodedLog<T> extends Omit<EVMLog, 'data'> { data: T }





/**
 * Filter params for Deposit. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type DepositTopics = {
  aset?: `0x${string}`
}

/**
 * Decoded Deposit event data.
 */
export type DepositDecoded = {
  aset: `0x${string}`
  amount: bigint
}


/**
 * Filter params for KeystoneForwarderRemoved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type KeystoneForwarderRemovedTopics = {
  keystoneForwarder?: `0x${string}`
}

/**
 * Decoded KeystoneForwarderRemoved event data.
 */
export type KeystoneForwarderRemovedDecoded = {
  keystoneForwarder: `0x${string}`
}


/**
 * Filter params for KeystoneForwarderSet. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type KeystoneForwarderSetTopics = {
  keystoneForwarder?: `0x${string}`
}

/**
 * Decoded KeystoneForwarderSet event data.
 */
export type KeystoneForwarderSetDecoded = {
  keystoneForwarder: `0x${string}`
}


/**
 * Filter params for MessageReceived. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type MessageReceivedTopics = {
  messageId?: `0x${string}`
  sourceChainSelector?: bigint
}

/**
 * Decoded MessageReceived event data.
 */
export type MessageReceivedDecoded = {
  messageId: `0x${string}`
  sourceChainSelector: bigint
  sender: `0x${string}`
  tokenAmount: { token: `0x${string}`; amount: bigint }
}


/**
 * Filter params for MessageSent. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type MessageSentTopics = {
  messageId?: `0x${string}`
  destinationChainSelector?: bigint
}

/**
 * Decoded MessageSent event data.
 */
export type MessageSentDecoded = {
  messageId: `0x${string}`
  destinationChainSelector: bigint
  receiver: `0x${string}`
  tokenAmount: { token: `0x${string}`; amount: bigint }
  fees: bigint
}


/**
 * Filter params for OwnershipTransferRequested. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type OwnershipTransferRequestedTopics = {
  from?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded OwnershipTransferRequested event data.
 */
export type OwnershipTransferRequestedDecoded = {
  from: `0x${string}`
  to: `0x${string}`
}


/**
 * Filter params for OwnershipTransferred. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type OwnershipTransferredTopics = {
  from?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded OwnershipTransferred event data.
 */
export type OwnershipTransferredDecoded = {
  from: `0x${string}`
  to: `0x${string}`
}


/**
 * Filter params for ReportReceived. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type ReportReceivedTopics = {
  workflowOwner?: `0x${string}`
  workflowName?: `0x${string}`
}

/**
 * Decoded ReportReceived event data.
 */
export type ReportReceivedDecoded = {
  workflowOwner: `0x${string}`
  workflowName: `0x${string}`
  params: { asset: `0x${string}`; amount: bigint; destinationChainSelector: bigint; destinationProtocolSmartWallet: `0x${string}` }
}


/**
 * Filter params for SenderForSourceChainRemoved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type SenderForSourceChainRemovedTopics = {
  sourceChainSelector?: bigint
  sender?: `0x${string}`
}

/**
 * Decoded SenderForSourceChainRemoved event data.
 */
export type SenderForSourceChainRemovedDecoded = {
  sourceChainSelector: bigint
  sender: `0x${string}`
}


/**
 * Filter params for SenderForSourceChainSet. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type SenderForSourceChainSetTopics = {
  sourceChainSelector?: bigint
  sender?: `0x${string}`
}

/**
 * Decoded SenderForSourceChainSet event data.
 */
export type SenderForSourceChainSetDecoded = {
  sourceChainSelector: bigint
  sender: `0x${string}`
}


/**
 * Filter params for Withdraw. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type WithdrawTopics = {
  asset?: `0x${string}`
  to?: `0x${string}`
}

/**
 * Decoded Withdraw event data.
 */
export type WithdrawDecoded = {
  asset: `0x${string}`
  to: `0x${string}`
  amount: bigint
}


/**
 * Filter params for WorkflowOwnerRemoved. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type WorkflowOwnerRemovedTopics = {
  workflowOwner?: `0x${string}`
}

/**
 * Decoded WorkflowOwnerRemoved event data.
 */
export type WorkflowOwnerRemovedDecoded = {
  workflowOwner: `0x${string}`
}


/**
 * Filter params for WorkflowOwnerSet. Only indexed fields can be used for filtering.
 * Indexed string/bytes must be passed as keccak256 hash (Hex).
 */
export type WorkflowOwnerSetTopics = {
  workflowOwner?: `0x${string}`
}

/**
 * Decoded WorkflowOwnerSet event data.
 */
export type WorkflowOwnerSetDecoded = {
  workflowOwner: `0x${string}`
}


export const ProtocolSmartWalletABI = [{"type":"constructor","inputs":[{"name":"_keystoneForwarders","type":"address[]","internalType":"address[]"},{"name":"_allowedWorkflowOwners","type":"address[]","internalType":"address[]"},{"name":"_pool","type":"address","internalType":"address"},{"name":"_router","type":"address","internalType":"address"},{"name":"_link","type":"address","internalType":"address"}],"stateMutability":"nonpayable"},{"type":"function","name":"acceptOwnership","inputs":[],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"allowedCcipSenders","inputs":[{"name":"chainSelector","type":"uint64","internalType":"uint64"}],"outputs":[{"name":"sender","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"allowedKeystoneForwarders","inputs":[{"name":"keystoneForwarder","type":"address","internalType":"address"}],"outputs":[{"name":"allowed","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"allowedWorkflowOwners","inputs":[{"name":"workflowOwner","type":"address","internalType":"address"}],"outputs":[{"name":"allowed","type":"bool","internalType":"bool"}],"stateMutability":"view"},{"type":"function","name":"ccipReceive","inputs":[{"name":"message","type":"tuple","internalType":"structClient.Any2EVMMessage","components":[{"name":"messageId","type":"bytes32","internalType":"bytes32"},{"name":"sourceChainSelector","type":"uint64","internalType":"uint64"},{"name":"sender","type":"bytes","internalType":"bytes"},{"name":"data","type":"bytes","internalType":"bytes"},{"name":"destTokenAmounts","type":"tuple[]","internalType":"structClient.EVMTokenAmount[]","components":[{"name":"token","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"}]}]}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"depositToPool","inputs":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"getPoolAddress","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"getRouter","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"onReport","inputs":[{"name":"metadata","type":"bytes","internalType":"bytes"},{"name":"report","type":"bytes","internalType":"bytes"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"owner","inputs":[],"outputs":[{"name":"","type":"address","internalType":"address"}],"stateMutability":"view"},{"type":"function","name":"pool","inputs":[],"outputs":[{"name":"","type":"address","internalType":"contractMockPool"}],"stateMutability":"view"},{"type":"function","name":"removeKeystoneForwarder","inputs":[{"name":"_keystoneForwarder","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"removeSenderForSourceChain","inputs":[{"name":"_sourceChainSelector","type":"uint64","internalType":"uint64"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"removeWorkflowOwner","inputs":[{"name":"_workflowOwner","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setKeystoneForwarder","inputs":[{"name":"_keystoneForwarder","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setPool","inputs":[{"name":"_pool","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setSenderForSourceChain","inputs":[{"name":"_sourceChainSelector","type":"uint64","internalType":"uint64"},{"name":"_sender","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"setWorkflowOwner","inputs":[{"name":"_workflowOwner","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"supportsInterface","inputs":[{"name":"interfaceId","type":"bytes4","internalType":"bytes4"}],"outputs":[{"name":"","type":"bool","internalType":"bool"}],"stateMutability":"pure"},{"type":"function","name":"transferOwnership","inputs":[{"name":"to","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"withdraw","inputs":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"to","type":"address","internalType":"address"}],"outputs":[{"name":"","type":"uint256","internalType":"uint256"}],"stateMutability":"nonpayable"},{"type":"function","name":"withdrawFromPool","inputs":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"to","type":"address","internalType":"address"}],"outputs":[],"stateMutability":"nonpayable"},{"type":"function","name":"withdrawFromPoolAndDepositCrossChain","inputs":[{"name":"params","type":"tuple","internalType":"structProtocolSmartWallet.RebalanceParams","components":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"destinationChainSelector","type":"uint64","internalType":"uint64"},{"name":"destinationProtocolSmartWallet","type":"address","internalType":"address"}]}],"outputs":[],"stateMutability":"nonpayable"},{"type":"event","name":"Deposit","inputs":[{"name":"aset","type":"address","indexed":true,"internalType":"address"},{"name":"amount","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"KeystoneForwarderRemoved","inputs":[{"name":"keystoneForwarder","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"KeystoneForwarderSet","inputs":[{"name":"keystoneForwarder","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"MessageReceived","inputs":[{"name":"messageId","type":"bytes32","indexed":true,"internalType":"bytes32"},{"name":"sourceChainSelector","type":"uint64","indexed":true,"internalType":"uint64"},{"name":"sender","type":"address","indexed":false,"internalType":"address"},{"name":"tokenAmount","type":"tuple","indexed":false,"internalType":"structClient.EVMTokenAmount","components":[{"name":"token","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"}]}],"anonymous":false},{"type":"event","name":"MessageSent","inputs":[{"name":"messageId","type":"bytes32","indexed":true,"internalType":"bytes32"},{"name":"destinationChainSelector","type":"uint64","indexed":true,"internalType":"uint64"},{"name":"receiver","type":"address","indexed":false,"internalType":"address"},{"name":"tokenAmount","type":"tuple","indexed":false,"internalType":"structClient.EVMTokenAmount","components":[{"name":"token","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"}]},{"name":"fees","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"OwnershipTransferRequested","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"OwnershipTransferred","inputs":[{"name":"from","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"ReportReceived","inputs":[{"name":"workflowOwner","type":"address","indexed":true,"internalType":"address"},{"name":"workflowName","type":"bytes10","indexed":true,"internalType":"bytes10"},{"name":"params","type":"tuple","indexed":false,"internalType":"structProtocolSmartWallet.RebalanceParams","components":[{"name":"asset","type":"address","internalType":"address"},{"name":"amount","type":"uint256","internalType":"uint256"},{"name":"destinationChainSelector","type":"uint64","internalType":"uint64"},{"name":"destinationProtocolSmartWallet","type":"address","internalType":"address"}]}],"anonymous":false},{"type":"event","name":"SenderForSourceChainRemoved","inputs":[{"name":"_sourceChainSelector","type":"uint64","indexed":true,"internalType":"uint64"},{"name":"_sender","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"SenderForSourceChainSet","inputs":[{"name":"_sourceChainSelector","type":"uint64","indexed":true,"internalType":"uint64"},{"name":"_sender","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"Withdraw","inputs":[{"name":"asset","type":"address","indexed":true,"internalType":"address"},{"name":"to","type":"address","indexed":true,"internalType":"address"},{"name":"amount","type":"uint256","indexed":false,"internalType":"uint256"}],"anonymous":false},{"type":"event","name":"WorkflowOwnerRemoved","inputs":[{"name":"workflowOwner","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"event","name":"WorkflowOwnerSet","inputs":[{"name":"workflowOwner","type":"address","indexed":true,"internalType":"address"}],"anonymous":false},{"type":"error","name":"InsufficientFeeTokenAmount","inputs":[]},{"type":"error","name":"InsufficientTokenAmount","inputs":[]},{"type":"error","name":"InvalidKeystoneForwarder","inputs":[]},{"type":"error","name":"InvalidRouter","inputs":[{"name":"router","type":"address","internalType":"address"}]},{"type":"error","name":"InvalidSenderAddress","inputs":[]},{"type":"error","name":"InvalidSourceChain","inputs":[]},{"type":"error","name":"InvalidWorkflowOwner","inputs":[]},{"type":"error","name":"MismatchedTokenAmount","inputs":[]},{"type":"error","name":"MustBeKeystoneForwarder","inputs":[]},{"type":"error","name":"NoSenderOnSourceChain","inputs":[{"name":"sourceChainSelector","type":"uint64","internalType":"uint64"}]},{"type":"error","name":"SafeERC20FailedOperation","inputs":[{"name":"token","type":"address","internalType":"address"}]},{"type":"error","name":"UnauthorizedWorkflowOwner","inputs":[{"name":"workflowOwner","type":"address","internalType":"address"}]},{"type":"error","name":"WrongSenderForSourceChain","inputs":[{"name":"sourceChainSelector","type":"uint64","internalType":"uint64"}]},{"type":"error","name":"ZeroAddress","inputs":[{"name":"index","type":"uint256","internalType":"uint256"}]}] as const

export class ProtocolSmartWallet {
  constructor(
    private readonly client: EVMClient,
    public readonly address: Address,
  ) {}

  allowedCcipSenders(
    runtime: Runtime<unknown>,
    chainSelector: bigint,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'allowedCcipSenders' as const,
      args: [chainSelector],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'allowedCcipSenders' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  allowedKeystoneForwarders(
    runtime: Runtime<unknown>,
    keystoneForwarder: `0x${string}`,
  ): boolean {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'allowedKeystoneForwarders' as const,
      args: [keystoneForwarder],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'allowedKeystoneForwarders' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  allowedWorkflowOwners(
    runtime: Runtime<unknown>,
    workflowOwner: `0x${string}`,
  ): boolean {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'allowedWorkflowOwners' as const,
      args: [workflowOwner],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'allowedWorkflowOwners' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  getPoolAddress(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'getPoolAddress' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'getPoolAddress' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  getRouter(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'getRouter' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'getRouter' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  owner(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'owner' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'owner' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  pool(
    runtime: Runtime<unknown>,
  ): `0x${string}` {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'pool' as const,
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'pool' as const,
      data: bytesToHex(result.data),
    }) as `0x${string}`
  }

  supportsInterface(
    runtime: Runtime<unknown>,
    interfaceId: `0x${string}`,
  ): boolean {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'supportsInterface' as const,
      args: [interfaceId],
    })

    const result = this.client
      .callContract(runtime, {
        call: encodeCallMsg({ from: zeroAddress, to: this.address, data: callData }),
        blockNumber: LAST_FINALIZED_BLOCK_NUMBER,
      })
      .result()

    return decodeFunctionResult({
      abi: ProtocolSmartWalletABI,
      functionName: 'supportsInterface' as const,
      data: bytesToHex(result.data),
    }) as boolean
  }

  writeReportFromCcipReceive(
    runtime: Runtime<unknown>,
    message: { messageId: `0x${string}`; sourceChainSelector: bigint; sender: `0x${string}`; data: `0x${string}`; destTokenAmounts: readonly { token: `0x${string}`; amount: bigint }[] },
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'ccipReceive' as const,
      args: [message],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromDepositToPool(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
    amount: bigint,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'depositToPool' as const,
      args: [asset, amount],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromOnReport(
    runtime: Runtime<unknown>,
    metadata: `0x${string}`,
    report: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'onReport' as const,
      args: [metadata, report],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromRemoveKeystoneForwarder(
    runtime: Runtime<unknown>,
    keystoneForwarder: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'removeKeystoneForwarder' as const,
      args: [keystoneForwarder],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromRemoveSenderForSourceChain(
    runtime: Runtime<unknown>,
    sourceChainSelector: bigint,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'removeSenderForSourceChain' as const,
      args: [sourceChainSelector],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromRemoveWorkflowOwner(
    runtime: Runtime<unknown>,
    workflowOwner: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'removeWorkflowOwner' as const,
      args: [workflowOwner],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetKeystoneForwarder(
    runtime: Runtime<unknown>,
    keystoneForwarder: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'setKeystoneForwarder' as const,
      args: [keystoneForwarder],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetPool(
    runtime: Runtime<unknown>,
    pool: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'setPool' as const,
      args: [pool],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetSenderForSourceChain(
    runtime: Runtime<unknown>,
    sourceChainSelector: bigint,
    sender: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'setSenderForSourceChain' as const,
      args: [sourceChainSelector, sender],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromSetWorkflowOwner(
    runtime: Runtime<unknown>,
    workflowOwner: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'setWorkflowOwner' as const,
      args: [workflowOwner],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromTransferOwnership(
    runtime: Runtime<unknown>,
    to: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'transferOwnership' as const,
      args: [to],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromWithdraw(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
    amount: bigint,
    to: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'withdraw' as const,
      args: [asset, amount, to],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromWithdrawFromPool(
    runtime: Runtime<unknown>,
    asset: `0x${string}`,
    amount: bigint,
    to: `0x${string}`,
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'withdrawFromPool' as const,
      args: [asset, amount, to],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReportFromWithdrawFromPoolAndDepositCrossChain(
    runtime: Runtime<unknown>,
    params: { asset: `0x${string}`; amount: bigint; destinationChainSelector: bigint; destinationProtocolSmartWallet: `0x${string}` },
    gasConfig?: { gasLimit?: string },
  ) {
    const callData = encodeFunctionData({
      abi: ProtocolSmartWalletABI,
      functionName: 'withdrawFromPoolAndDepositCrossChain' as const,
      args: [params],
    })

    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  writeReport(
    runtime: Runtime<unknown>,
    callData: Hex,
    gasConfig?: { gasLimit?: string },
  ) {
    const reportResponse = runtime
      .report(prepareReportRequest(callData))
      .result()

    return this.client
      .writeReport(runtime, {
        receiver: this.address,
        report: reportResponse,
        gasConfig,
      })
      .result()
  }

  /**
   * Creates a log trigger for Deposit events.
   * The returned trigger's adapt method decodes the raw log into DepositDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerDeposit(
    filters?: DepositTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'Deposit' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        aset: f.aset,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'Deposit' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          aset: f.aset,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'Deposit' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<DepositDecoded> => contract.decodeDeposit(rawOutput),
    }
  }

  /**
   * Decodes a log into Deposit data, preserving all log metadata.
   */
  decodeDeposit(log: EVMLog): DecodedLog<DepositDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as DepositDecoded }
  }

  /**
   * Creates a log trigger for KeystoneForwarderRemoved events.
   * The returned trigger's adapt method decodes the raw log into KeystoneForwarderRemovedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerKeystoneForwarderRemoved(
    filters?: KeystoneForwarderRemovedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'KeystoneForwarderRemoved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        keystoneForwarder: f.keystoneForwarder,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'KeystoneForwarderRemoved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          keystoneForwarder: f.keystoneForwarder,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'KeystoneForwarderRemoved' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<KeystoneForwarderRemovedDecoded> => contract.decodeKeystoneForwarderRemoved(rawOutput),
    }
  }

  /**
   * Decodes a log into KeystoneForwarderRemoved data, preserving all log metadata.
   */
  decodeKeystoneForwarderRemoved(log: EVMLog): DecodedLog<KeystoneForwarderRemovedDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as KeystoneForwarderRemovedDecoded }
  }

  /**
   * Creates a log trigger for KeystoneForwarderSet events.
   * The returned trigger's adapt method decodes the raw log into KeystoneForwarderSetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerKeystoneForwarderSet(
    filters?: KeystoneForwarderSetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'KeystoneForwarderSet' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        keystoneForwarder: f.keystoneForwarder,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'KeystoneForwarderSet' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          keystoneForwarder: f.keystoneForwarder,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'KeystoneForwarderSet' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<KeystoneForwarderSetDecoded> => contract.decodeKeystoneForwarderSet(rawOutput),
    }
  }

  /**
   * Decodes a log into KeystoneForwarderSet data, preserving all log metadata.
   */
  decodeKeystoneForwarderSet(log: EVMLog): DecodedLog<KeystoneForwarderSetDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as KeystoneForwarderSetDecoded }
  }

  /**
   * Creates a log trigger for MessageReceived events.
   * The returned trigger's adapt method decodes the raw log into MessageReceivedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerMessageReceived(
    filters?: MessageReceivedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'MessageReceived' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        messageId: f.messageId,
        sourceChainSelector: f.sourceChainSelector,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'MessageReceived' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          messageId: f.messageId,
          sourceChainSelector: f.sourceChainSelector,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'MessageReceived' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<MessageReceivedDecoded> => contract.decodeMessageReceived(rawOutput),
    }
  }

  /**
   * Decodes a log into MessageReceived data, preserving all log metadata.
   */
  decodeMessageReceived(log: EVMLog): DecodedLog<MessageReceivedDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as MessageReceivedDecoded }
  }

  /**
   * Creates a log trigger for MessageSent events.
   * The returned trigger's adapt method decodes the raw log into MessageSentDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerMessageSent(
    filters?: MessageSentTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'MessageSent' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        messageId: f.messageId,
        destinationChainSelector: f.destinationChainSelector,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'MessageSent' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          messageId: f.messageId,
          destinationChainSelector: f.destinationChainSelector,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'MessageSent' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<MessageSentDecoded> => contract.decodeMessageSent(rawOutput),
    }
  }

  /**
   * Decodes a log into MessageSent data, preserving all log metadata.
   */
  decodeMessageSent(log: EVMLog): DecodedLog<MessageSentDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as MessageSentDecoded }
  }

  /**
   * Creates a log trigger for OwnershipTransferRequested events.
   * The returned trigger's adapt method decodes the raw log into OwnershipTransferRequestedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerOwnershipTransferRequested(
    filters?: OwnershipTransferRequestedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'OwnershipTransferRequested' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        from: f.from,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'OwnershipTransferRequested' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          from: f.from,
          to: f.to,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'OwnershipTransferRequested' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<OwnershipTransferRequestedDecoded> => contract.decodeOwnershipTransferRequested(rawOutput),
    }
  }

  /**
   * Decodes a log into OwnershipTransferRequested data, preserving all log metadata.
   */
  decodeOwnershipTransferRequested(log: EVMLog): DecodedLog<OwnershipTransferRequestedDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferRequestedDecoded }
  }

  /**
   * Creates a log trigger for OwnershipTransferred events.
   * The returned trigger's adapt method decodes the raw log into OwnershipTransferredDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerOwnershipTransferred(
    filters?: OwnershipTransferredTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'OwnershipTransferred' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        from: f.from,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'OwnershipTransferred' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          from: f.from,
          to: f.to,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'OwnershipTransferred' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<OwnershipTransferredDecoded> => contract.decodeOwnershipTransferred(rawOutput),
    }
  }

  /**
   * Decodes a log into OwnershipTransferred data, preserving all log metadata.
   */
  decodeOwnershipTransferred(log: EVMLog): DecodedLog<OwnershipTransferredDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as OwnershipTransferredDecoded }
  }

  /**
   * Creates a log trigger for ReportReceived events.
   * The returned trigger's adapt method decodes the raw log into ReportReceivedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerReportReceived(
    filters?: ReportReceivedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'ReportReceived' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        workflowOwner: f.workflowOwner,
        workflowName: f.workflowName,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'ReportReceived' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          workflowOwner: f.workflowOwner,
          workflowName: f.workflowName,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'ReportReceived' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<ReportReceivedDecoded> => contract.decodeReportReceived(rawOutput),
    }
  }

  /**
   * Decodes a log into ReportReceived data, preserving all log metadata.
   */
  decodeReportReceived(log: EVMLog): DecodedLog<ReportReceivedDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as ReportReceivedDecoded }
  }

  /**
   * Creates a log trigger for SenderForSourceChainRemoved events.
   * The returned trigger's adapt method decodes the raw log into SenderForSourceChainRemovedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerSenderForSourceChainRemoved(
    filters?: SenderForSourceChainRemovedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'SenderForSourceChainRemoved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        sourceChainSelector: f.sourceChainSelector,
        sender: f.sender,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'SenderForSourceChainRemoved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          sourceChainSelector: f.sourceChainSelector,
          sender: f.sender,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'SenderForSourceChainRemoved' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<SenderForSourceChainRemovedDecoded> => contract.decodeSenderForSourceChainRemoved(rawOutput),
    }
  }

  /**
   * Decodes a log into SenderForSourceChainRemoved data, preserving all log metadata.
   */
  decodeSenderForSourceChainRemoved(log: EVMLog): DecodedLog<SenderForSourceChainRemovedDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as SenderForSourceChainRemovedDecoded }
  }

  /**
   * Creates a log trigger for SenderForSourceChainSet events.
   * The returned trigger's adapt method decodes the raw log into SenderForSourceChainSetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerSenderForSourceChainSet(
    filters?: SenderForSourceChainSetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'SenderForSourceChainSet' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        sourceChainSelector: f.sourceChainSelector,
        sender: f.sender,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'SenderForSourceChainSet' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          sourceChainSelector: f.sourceChainSelector,
          sender: f.sender,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'SenderForSourceChainSet' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<SenderForSourceChainSetDecoded> => contract.decodeSenderForSourceChainSet(rawOutput),
    }
  }

  /**
   * Decodes a log into SenderForSourceChainSet data, preserving all log metadata.
   */
  decodeSenderForSourceChainSet(log: EVMLog): DecodedLog<SenderForSourceChainSetDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as SenderForSourceChainSetDecoded }
  }

  /**
   * Creates a log trigger for Withdraw events.
   * The returned trigger's adapt method decodes the raw log into WithdrawDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerWithdraw(
    filters?: WithdrawTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'Withdraw' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        asset: f.asset,
        to: f.to,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'Withdraw' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          asset: f.asset,
          to: f.to,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'Withdraw' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<WithdrawDecoded> => contract.decodeWithdraw(rawOutput),
    }
  }

  /**
   * Decodes a log into Withdraw data, preserving all log metadata.
   */
  decodeWithdraw(log: EVMLog): DecodedLog<WithdrawDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as WithdrawDecoded }
  }

  /**
   * Creates a log trigger for WorkflowOwnerRemoved events.
   * The returned trigger's adapt method decodes the raw log into WorkflowOwnerRemovedDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerWorkflowOwnerRemoved(
    filters?: WorkflowOwnerRemovedTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'WorkflowOwnerRemoved' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        workflowOwner: f.workflowOwner,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'WorkflowOwnerRemoved' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          workflowOwner: f.workflowOwner,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'WorkflowOwnerRemoved' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<WorkflowOwnerRemovedDecoded> => contract.decodeWorkflowOwnerRemoved(rawOutput),
    }
  }

  /**
   * Decodes a log into WorkflowOwnerRemoved data, preserving all log metadata.
   */
  decodeWorkflowOwnerRemoved(log: EVMLog): DecodedLog<WorkflowOwnerRemovedDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as WorkflowOwnerRemovedDecoded }
  }

  /**
   * Creates a log trigger for WorkflowOwnerSet events.
   * The returned trigger's adapt method decodes the raw log into WorkflowOwnerSetDecoded,
   * so the handler receives typed event data directly.
   * When multiple filters are provided, topic values are merged with OR semantics (match any).
   */
  logTriggerWorkflowOwnerSet(
    filters?: WorkflowOwnerSetTopics[],
  ) {
    let topics: { values: string[] }[]
    if (!filters || filters.length === 0) {
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'WorkflowOwnerSet' as const,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else if (filters.length === 1) {
      const f = filters[0]
      const args = {
        workflowOwner: f.workflowOwner,
      }
      const encoded = encodeEventTopics({
        abi: ProtocolSmartWalletABI,
        eventName: 'WorkflowOwnerSet' as const,
        args,
      })
      topics = encoded.map((t) => ({ values: [hexToBase64(t)] }))
    } else {
      const allEncoded = filters.map((f) => {
        const args = {
          workflowOwner: f.workflowOwner,
        }
        return encodeEventTopics({
          abi: ProtocolSmartWalletABI,
          eventName: 'WorkflowOwnerSet' as const,
          args,
        })
      })
      topics = allEncoded[0].map((_, i) => ({
        values: [...new Set(allEncoded.map((row) => hexToBase64(row[i])))],
      }))
    }
    const baseTrigger = this.client.logTrigger({
      addresses: [hexToBase64(this.address)],
      topics,
    })
    const contract = this
    return {
      capabilityId: () => baseTrigger.capabilityId(),
      method: () => baseTrigger.method(),
      outputSchema: () => baseTrigger.outputSchema(),
      configAsAny: () => baseTrigger.configAsAny(),
      adapt: (rawOutput: EVMLog): DecodedLog<WorkflowOwnerSetDecoded> => contract.decodeWorkflowOwnerSet(rawOutput),
    }
  }

  /**
   * Decodes a log into WorkflowOwnerSet data, preserving all log metadata.
   */
  decodeWorkflowOwnerSet(log: EVMLog): DecodedLog<WorkflowOwnerSetDecoded> {
    const decoded = decodeEventLog({
      abi: ProtocolSmartWalletABI,
      data: bytesToHex(log.data),
      topics: log.topics.map((t) => bytesToHex(t)) as readonly Hex[],
    })
    const { data: _, ...rest } = log
    return { ...rest, data: decoded.args as unknown as WorkflowOwnerSetDecoded }
  }
}

