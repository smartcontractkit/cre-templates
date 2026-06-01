import { EVMClient, type Runtime } from "@chainlink/cre-sdk";
import { IConfig } from "../interfaces/IConfig";

export interface BridgeInputProps {
    runtime: Runtime<IConfig>;
    evmClient: EVMClient;
    from: {
        chain: string;
        token: string;
        depositor: string;
    };
    to: {
        chain: string;
        token: string;
        recipient: string;
    };
    amount: string | BigInt;
    sourceConfigContract: string;
}

export interface BridgeFactoryCCIPProps extends BridgeInputProps {
    destinationChainSelector: string;
}

export interface IBridge {
    bridgeTokenService(props: BridgeInputProps): any;
    addExtraParams(props: IConfig): any;
    bridgeName(): string;
}