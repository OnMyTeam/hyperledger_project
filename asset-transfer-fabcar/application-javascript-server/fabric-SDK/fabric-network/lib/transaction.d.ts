/// <reference types="node" />
import { Endorser } from 'fabric-common';
import { ContractImpl } from './contract';
import { TxEventHandlerFactory } from './impl/event/transactioneventhandler';
export interface TransientMap {
    [key: string]: Buffer;
}
/**
 * Represents a specific invocation of a transaction function, and provides
 * flexibility over how that transaction is invoked. Applications should
 * obtain instances of this class by calling
 * [Contract#createTransaction()]{@link module:fabric-network.Contract#createTransaction}.
 * <br><br>
 * Instances of this class are stateful. A new instance <strong>must</strong>
 * be created for each transaction invocation.
 * @memberof module:fabric-network
 * @hideconstructor
 */
export declare class Transaction {
    private readonly name;
    private readonly contract;
    private transientMap?;
    private readonly gatewayOptions;
    private eventHandlerStrategyFactory;
    private readonly queryHandler;
    private endorsingPeers?;
    private endorsingOrgs?;
    private readonly identityContext;
    private readonly transactionId;
    constructor(contract: ContractImpl, name: string);
    /**
     * Get the fully qualified name of the transaction function.
     * @returns {string} Transaction name.
     */
    getName(): string;
    /**
     * Set transient data that will be passed to the transaction function
     * but will not be stored on the ledger. This can be used to pass
     * private data to a transaction function.
     * @param {Object} transientMap Object with String property names and
     * Buffer property values.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setTransient(transientMap: TransientMap): Transaction;
    /**
     * Get the ID that will be used for this transaction invocation.
     * @returns {string} A transaction ID.
     */
    getTransactionId(): string;
    /**
     * Set the peers that should be used for endorsement when this transaction
     * is submitted to the ledger.
     * Setting the peers will override the use of discovery and the submit will
     * send the proposal to these peers.
     * This will override the setEndorsingOrganizations if previously called.
     * @param {Endorser[]} peers - Endorsing peers.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setEndorsingPeers(peers: Endorser[]): Transaction;
    /**
     * Set the organizations that should be used for endorsement when this
     * transaction is submitted to the ledger.
     * Peers that are in the organizations will be used for the endorsement.
     * This will override the setEndorsingPeers if previously called. Setting
     * the endorsing organizations will not override discovery, however it will
     * filter the peers provided by discovery to be those in these organizatons.
     * @param {string[]} orgs - Endorsing organizations.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setEndorsingOrganizations(...orgs: string[]): Transaction;
    /**
     * Set an event handling strategy to use for this transaction instead of the default configured on the gateway.
     * @param strategy An event handling strategy.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setEventHandler(strategy: TxEventHandlerFactory): Transaction;
    /**
     * Submit a transaction to the ledger. The transaction function <code>name</code>
     * will be evaluated on the endorsing peers and then submitted to the ordering service
     * for committing to the ledger.
     * @async
     * @param {...string} [args] Transaction function arguments.
     * @returns {Buffer} Payload response from the transaction function.
     * @throws {module:fabric-network.TimeoutError} If the transaction was successfully submitted to the orderer but
     * timed out before a commit event was received from peers.
     */
    submit(...args: string[]): Promise<Buffer>;

    submitAggregator(tag: JSON, ...args: string[]): Promise<Buffer>;    
    /**
     * Evaluate a transaction function and return its results.
     * The transaction function will be evaluated on the endorsing peers but
     * the responses will not be sent to the ordering service and hence will
     * not be committed to the ledger.
     * This is used for querying the world state.
     * @async
     * @param {...string} [args] Transaction function arguments.
     * @returns {Promise<Buffer>} Payload response from the transaction function.
     */
    evaluate(...args: string[]): Promise<Buffer>;
    private newBuildProposalRequest;
}
