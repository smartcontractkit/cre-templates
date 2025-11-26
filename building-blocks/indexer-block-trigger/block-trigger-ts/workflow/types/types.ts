export type Config = {
  watchedAddresses: string[];
};

export type AlchemyWebhookPayload = {
  webhookId: string;
  id: string;
  createdAt: string;
  type: string;
  event: {
    data: {
      block: {
        hash: string;
        number: number;
        timestamp: number;
        logs: Array<{
          data: string;
          topics: string[];
          index: number;
          account: {
            address: string;
          };
          transaction: {
            hash: string;
            nonce: number;
            index: number;
            from: {
              address: string;
            };
            to: {
              address: string;
            } | null;
            value: string;
            gasPrice: string;
            maxFeePerGas: string | null;
            maxPriorityFeePerGas: string | null;
            gas: number;
            status: number;
            gasUsed: number;
            cumulativeGasUsed: number;
            effectiveGasPrice: string;
            createdContract?: {
              address: string;
            };
          };
        }>;
      };
    };
  };
};

export type Transaction = {
  hash: string;
  nonce: number;
  index: number;
  from: string;
  to: string | null;
  value: string;
  gasPrice: string;
  gas: number;
  status: number;
  gasUsed: number;
  blockNumber: number;
  blockHash: string;
  timestamp: number;
};

export type TransactionStore = {
  watchedAddresses: Set<string>;
  transactions: Map<string, Transaction>;
};
