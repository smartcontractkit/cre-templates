import { cre, type Runtime, type EVMLog, bytesToHex, getNetwork } from "@chainlink/cre-sdk";
import { IConfig } from "./interfaces/IConfig"
import { parseAbi, decodeEventLog } from "viem";
import { sendTelegramMessage } from "./telegramMessageService";
import { BridgeFactory } from "./bridge/factory"
import { getEvmClient } from "./utils/getEvmClient";
import { Uniflow } from "./evm/generated/Uniflow";
import { getERC20Decimals, getERC20Allowance } from "./utils/erc20Utils";

const TRANSFER_EVENT_ABI = parseAbi(["event Transfer(address indexed from, address indexed to, uint256 amount)"]);

export const ethereumOnReceiveToken = (runtime: Runtime<IConfig>, evmLog: EVMLog): any => {
    const topics = evmLog.topics.map((t: Uint8Array) => bytesToHex(t)) as [
        `0x${string}`,
        ...`0x${string}`[],
    ];

    const data = bytesToHex(evmLog.data);
    const decodedLog = decodeEventLog({abi: TRANSFER_EVENT_ABI, topics: topics, data: data});
    
    runtime.log(decodedLog.args.to);
    runtime.log(`${decodedLog.args.amount}`);
    runtime.log(bytesToHex(evmLog.eventSig));

    const networksConfig = runtime.config.networks;
    const { evmClient } = getEvmClient(...networksConfig['eth'].creNetworkConfig);

    const tokenDecimals = getERC20Decimals(bytesToHex(evmLog.address).toString(), runtime, evmClient);
    sendTelegramMessage(runtime, `Token sent by: ${decodedLog.args.from} of amount ${decodedLog.args.amount}, decimals ${tokenDecimals} on ${bytesToHex(evmLog.address)}`);
    
    // @todo check for approval given by from address to the config contract
    const allowances = getERC20Allowance(bytesToHex(evmLog.address).toString(), runtime, evmClient, { owner: decodedLog.args.to, spender: networksConfig['eth'].configContract });

    runtime.log(`Allowance to config contract for token: ${allowances.toString()}`);

    if (allowances < decodedLog.args.amount) {
        return { message: "error", error: "Insufficient allowance for config contract" };
    }
    
    const uniflowContract = new Uniflow(evmClient, networksConfig['eth'].configContract as `0x${string}`);
    const owner = uniflowContract.owner(runtime);
    runtime.log(owner);

    const tokenAddress =  bytesToHex(evmLog.address).toString() as `0x${string}`;
    const unichainTokenConfig = runtime.config.networks['eth'].tokenMap[tokenAddress];
    const bridge = BridgeFactory.getBridge(unichainTokenConfig.bridge);

    runtime.log(`Bridging via: ${bridge.bridgeName()}`);

    const extraParams = bridge.addExtraParams(runtime.config);
    const { success, txHash, error } = bridge.bridgeTokenService({
        runtime,
        evmClient,
        from: {
            chain: "11155111",
            token: tokenAddress,
            depositor: networksConfig['eth'].configContract,
        },
        to: {
            chain: runtime.config.unichain.chainId,
            token: unichainTokenConfig.unichainToken,
            recipient: runtime.config.unichain.unichainDestinationAddress,
        },
        amount: decodedLog.args.amount,
        sourceConfigContract: networksConfig['eth'].configContract,
        ...extraParams
    });

    if (success) {
        return { message: "success", txHash: txHash || "", error: error || "" };
    }

    runtime.log(`error: ${error}`);
    return { message: "error", error: error };
}