import { getNetwork, EVMClient } from "@chainlink/cre-sdk";

export const getEvmClient = (...args: Parameters<typeof getNetwork>): { network: ReturnType<typeof getNetwork>, evmClient: InstanceType<typeof EVMClient> } => {
    const network = getNetwork(...args);
    if (!network) {
        throw new Error("Network config incorrect");
    }
    return {
        network,
        evmClient: new EVMClient(network.chainSelector.selector),
    }
}