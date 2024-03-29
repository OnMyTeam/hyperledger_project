import Long = require('long');
export interface Checkpointer {
    addTransactionId(transactionId: string): Promise<void>;
    getBlockNumber(): Promise<Long | undefined>;
    getTransactionIds(): Promise<Set<string>>;
    setBlockNumber(blockNumber: Long): Promise<void>;
}
