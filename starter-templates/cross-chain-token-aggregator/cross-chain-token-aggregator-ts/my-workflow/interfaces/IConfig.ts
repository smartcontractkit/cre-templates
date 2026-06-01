import { getNetwork } from "@chainlink/cre-sdk";

export interface IConfig {
    networks: {
        [key: string]: {
            chainId: string;
            creNetworkConfig: Parameters<typeof getNetwork>;
            configContract: string;
            targetUserAddress: string;
            tokenArr: Array<string>;
            tokenMap: {
                [key: `0x${string}`]: {
                    unichainToken: `0x${string}`;
                    bridge: string;
                };
            },
        };
    };
    unichain: {
        chainId: string;
        unichainDestinationAddress: string;
        configContract: string;
        chainlinkCCIPSelector: string;
    };
    telegramMessageApi: string;
    acrossApiUrl: string;
}