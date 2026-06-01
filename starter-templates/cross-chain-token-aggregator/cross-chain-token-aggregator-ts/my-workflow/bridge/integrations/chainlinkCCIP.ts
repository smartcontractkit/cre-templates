import { hexToBase64, TxStatus } from "@chainlink/cre-sdk";
import { BridgeFactoryCCIPProps, BridgeInputProps, IBridge } from "../IBridge";
import { parseAbiParameters, encodeAbiParameters, decodeFunctionData, erc20Abi, concatHex } from "viem";
import { IConfig } from "../../interfaces/IConfig";

export class ChainlinkCCIPBridge implements IBridge {
    bridgeTokenService(props: BridgeFactoryCCIPProps): any {
        const ccipReportParams = parseAbiParameters("address receiver, address token, uint256 amount, uint64 unichainSelector");
        const encodedReportParams = encodeAbiParameters(
            ccipReportParams,
            [
                props.to.recipient as `0x${string}`,
                props.from.token as `0x${string}`,
                BigInt(props.amount.toString()),
                BigInt(props.destinationChainSelector)
            ]
        );

        const signedReport = props.runtime.report({
            encodedPayload: hexToBase64(concatHex(["0x02", encodedReportParams])),
            encoderName: "evm",
            signingAlgo: "ecdsa",
            hashingAlgo: "keccak256"
        }).result();

        const bridgeReportResponse = props.evmClient.writeReport(props.runtime, {
            receiver: props.sourceConfigContract,
            report: signedReport,
            gasConfig: {
                gasLimit: "500000"
            }
        }).result();

        if (bridgeReportResponse.txStatus == TxStatus.SUCCESS) {
            return { success: true, txHash: bridgeReportResponse.txHash };
        }
        
        return { success: false, error: bridgeReportResponse.errorMessage };
    }

    addExtraParams(props: IConfig): Omit<BridgeFactoryCCIPProps, keyof BridgeInputProps> {
        return {
            destinationChainSelector: props.unichain.chainlinkCCIPSelector
        };
    }

    bridgeName(): string {
        return "Chainlink CCIP";
    }
}