// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { IERC20ABI } from './IERC20'

export type IERC20Mock = {
  allowance?: (owner: `0x${string}`, spender: `0x${string}`) => bigint
  balanceOf?: (account: `0x${string}`) => bigint
  totalSupply?: () => bigint
} & Pick<ContractMock<typeof IERC20ABI>, 'writeReport'>

export function newIERC20Mock(address: Address, evmMock: EvmMock): IERC20Mock {
  return addContractMock(evmMock, { address, abi: IERC20ABI }) as IERC20Mock
}

