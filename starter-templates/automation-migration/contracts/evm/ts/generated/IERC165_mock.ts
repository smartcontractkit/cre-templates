// Code generated — DO NOT EDIT.
import type { Address } from 'viem'
import { addContractMock, type ContractMock, type EvmMock } from '@chainlink/cre-sdk/test'

import { IERC165ABI } from './IERC165'

export type IERC165Mock = {
  supportsInterface?: (interfaceId: `0x${string}`) => boolean
} & Pick<ContractMock<typeof IERC165ABI>, 'writeReport'>

export function newIERC165Mock(address: Address, evmMock: EvmMock): IERC165Mock {
  return addContractMock(evmMock, { address, abi: IERC165ABI }) as IERC165Mock
}

