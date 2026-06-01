import { cre, Runner } from "@chainlink/cre-sdk";
import { getEvmClient } from "./utils/getEvmClient";
import { IConfig } from "./interfaces/IConfig"
import { ethereumOnReceiveToken } from "./ethereumOnReceiveToken";
import { toHex, keccak256 } from "viem";

const initWorkflow = (config: IConfig) => {
  const networksConfig = config.networks;
  const { evmClient: ethereumClient } = getEvmClient(...networksConfig['eth'].creNetworkConfig);
  const TRANSFER_EVENT = keccak256(toHex("Transfer(address,address,uint256)"));
  const ethereumLogTrigger = getLogTrigger(ethereumClient, TRANSFER_EVENT, networksConfig['eth'].targetUserAddress, networksConfig['eth'].tokenArr);

  return [
    cre.handler(
      ethereumLogTrigger,
      ethereumOnReceiveToken,
    ),
  ];
};

const getLogTrigger = (evmClient: InstanceType<typeof cre.capabilities.EVMClient>, transferEvent: `0x${string}`, toAddress: string, token: Array<string>) => {
  return evmClient.logTrigger({
    addresses: [...token],
    topics: [
      { values: [transferEvent] },
      { values: [] },
      { values: [toAddress] },
    ],
    confidence: 'CONFIDENCE_LEVEL_FINALIZED',
  });
}

export async function main() {
  const runner = await Runner.newRunner<IConfig>();
  await runner.run(initWorkflow);
}

main();
