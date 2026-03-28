import { consensusIdenticalAggregation, cre, HTTPSendRequester, ok, hexToBase64, TxStatus } from "@chainlink/cre-sdk";
import { BridgeInputProps, IBridge } from "../IBridge";
import { buildUrlWithParams } from "../../utils/buildUrlWithParams";
import { parseAbiParameters, encodeAbiParameters, decodeFunctionData, erc20Abi, concatHex } from "viem";
import { IConfig } from "../../interfaces/IConfig";

interface AcrossResponse {
  approvalTxns: Array<{
    chainId: number;
    to: string;
    data: string;
  }>;
  swapTx: {
    ecosystem: string;
    simulationSuccess: boolean;
    chainId: number;
    to: string;
    data: string;
  },
};

export class AcrossBridge implements IBridge {
  bridgeTokenService(props: BridgeInputProps): any {
    const httpClient = new cre.capabilities.HTTPClient();

    const acrossBridgeRequestResponse = httpClient.sendRequest(
      props.runtime,
      this.initializeAcrossBridge(props),
      consensusIdenticalAggregation<AcrossResponse>()
    )().result();

    const approvalTxnData = acrossBridgeRequestResponse.approvalTxns?.[0].data as `0x${string}`;
    const { args } = decodeFunctionData({
      abi: erc20Abi,
      data: approvalTxnData,
    });
    const [approvalContract] = args;

    props.runtime.log(`Approval Contract for Across Bridge: ${approvalContract}`);

    // send report to config contract to process token approval and bridging
    const acrossBridgeParams = parseAbiParameters("address receiver, uint256 amount, address token, address approvalContract, address depositContract, bytes memory depositData");
    const encodedBridgeParams = encodeAbiParameters(
      acrossBridgeParams,
      [
        props.to.recipient as `0x${string}`,
        BigInt(props.amount.toString()),
        props.from.token as `0x${string}`,
        approvalContract as `0x${string}`,
        acrossBridgeRequestResponse.swapTx.to as `0x${string}`,
        acrossBridgeRequestResponse.swapTx.data as `0x${string}`
      ],
    );

    const signedReport = props.runtime.report({
      encodedPayload: hexToBase64(concatHex(["0x01", encodedBridgeParams])),
      encoderName: "evm",
      signingAlgo: "ecdsa",
      hashingAlgo: "keccak256",
    }).result();

    props.runtime.log(`Report receiver: ${props.sourceConfigContract}`);

    const bridgeReportResponse = props.evmClient.writeReport(props.runtime, {
      receiver: props.sourceConfigContract,
      report: signedReport,
      gasConfig: {
        gasLimit: "1000000",
      },
    }).result();

    props.runtime.log(`${bridgeReportResponse.receiverContractExecutionStatus}`);

    if (bridgeReportResponse.txStatus == TxStatus.SUCCESS) {
      return { success: true, txHash: bridgeReportResponse.txHash, acrossBridgeRequest: acrossBridgeRequestResponse, report: concatHex(["0x01", encodedBridgeParams]) };
    }

    return { success: false, error: bridgeReportResponse.errorMessage, acrossBridgeRequest: acrossBridgeRequestResponse };
  }

  private initializeAcrossBridge(props: BridgeInputProps) {
    return (sendRequester: HTTPSendRequester): AcrossResponse => {
      const params = {
        tradeType: 'exactInput',
        amount: props.amount.toString(),
        inputToken: props.from.token,
        originChainId: props.from.chain,
        outputToken: props.to.token,
        destinationChainId: props.to.chain,
        depositor: props.from.depositor,
        recipient: props.to.recipient,
      };

      const swapUrl = `${props.runtime.config.acrossApiUrl}/api/swap/approval`;

      const req = {
        url: buildUrlWithParams(swapUrl, params),
        method: "GET" as const,
        cacheSettings: {
          store: true,
          maxAge: '60s',
        },
      };

      const resp = sendRequester.sendRequest(req).result();
      const bodyText = new TextDecoder().decode(resp.body);

      if (!ok(resp)) {
        throw new Error(`Across bridge initialize request error: ${resp.statusCode} - ${bodyText}`);
      }

      const apiResponse = JSON.parse(bodyText) as AcrossResponse;
      return apiResponse;
    }
  }

  addExtraParams(props: IConfig) {
    return {};
  }

  bridgeName(): string {
    return "Across";
  }
}