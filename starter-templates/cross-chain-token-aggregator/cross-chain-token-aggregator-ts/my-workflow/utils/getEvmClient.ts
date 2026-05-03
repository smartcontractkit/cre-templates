import { cre, getNetwork } from "@chainlink/cre-sdk";

export const getEvmClient = (...args: Parameters<typeof getNetwork>): { network: ReturnType<typeof getNetwork>, evmClient: InstanceType<typeof cre.capabilities.EVMClient> } => {
    const network = getNetwork(...args);
    if (!network) {
        throw new Error("Network config incorrect");
    }
    return {
        network,
        evmClient: new cre.capabilities.EVMClient(network.chainSelector.selector),
    }
}