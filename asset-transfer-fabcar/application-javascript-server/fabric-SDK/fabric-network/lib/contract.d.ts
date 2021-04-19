/// <reference types="node" />
import { Transaction } from './transaction';
import { NetworkImpl } from './network';
import { ContractListener, ListenerOptions } from './events';
import { DiscoveryHandler } from 'fabric-common';
import { Gateway } from './gateway';
export interface DiscoveryInterest {
    name: string;
    collectionNames?: string[];
    noPrivateReads?: boolean;
}
export interface Contract {
    readonly chaincodeId: string;
    readonly namespace: string;
    createTransaction(name: string): Transaction;
    evaluateTransaction(name: string, ...args: string[]): Promise<Buffer>;
    submitTransaction(name: string, ...args: string[]): Promise<Buffer>;
    addContractListener(listener: ContractListener, options?: ListenerOptions): Promise<ContractListener>;
    removeContractListener(listener: ContractListener): void;
    addDiscoveryInterest(interest: DiscoveryInterest): Contract;
    resetDiscoveryInterests(): Contract;
}
/**
 * <p>Represents a smart contract (chaincode) instance in a network.
 * Applications should get a Contract instance using the
 * networks's [getContract]{@link module:fabric-network.Network#getContract} method.</p>
 *
 * <p>The Contract allows applications to:</p>
 * <ul>
 *   <li>Submit transactions that store state to the ledger using
 *       [submitTransaction]{@link module:fabric-network.Contract#submitTransaction}.</li>
 *   <li>Evaluate transactions that query state from the ledger using
 *       [evaluateTransaction]{@link module:fabric-network.Contract#evaluateTransaction}.</li>
 *   <li>Listen for new chaincode events and replay previous chaincode events emitted by the smart contract using
 *       [addContractListener]{@link module:fabric-network.Contract#addContractListener}.</li>
 * </ul>
 *
 * <p>If more control over transaction invocation is required, such as including transient data,
 * [createTransaction]{@link module:fabric-network.Contract#createTransaction} can be used to build a transaction
 * request that is submitted to or evaluated by the smart contract.</p>
 * @interface Contract
 * @memberof module:fabric-network
 */
/**
 * Create an object representing a specific invocation of a transaction
 * function implemented by this contract, and provides more control over
 * the transaction invocation. A new transaction object <strong>must</strong>
 * be created for each transaction invocation.
 * @method Contract#createTransaction
 * @memberof module:fabric-network
 * @param {string} name Transaction function name.
 * @returns {module:fabric-network.Transaction} A transaction object.
 */
/**
 * Submit a transaction to the ledger. The transaction function <code>name</code>
 * will be evaluated on the endorsing peers and then submitted to the ordering service
 * for committing to the ledger.
 * This function is equivalent to calling <code>createTransaction(name).submit()</code>.
 * @method Contract#submitTransaction
 * @memberof module:fabric-network
 * @param {string} name Transaction function name.
 * @param {...string} [args] Transaction function arguments.
 * @returns {Buffer} Payload response from the transaction function.
 * @throws {module:fabric-network.TimeoutError} If the transaction was successfully submitted to the orderer but
 * timed out before a commit event was received from peers.
 */
/**
 * Evaluate a transaction function and return its results.
 * The transaction function <code>name</code>
 * will be evaluated on the endorsing peers but the responses will not be sent to
 * the ordering service and hence will not be committed to the ledger.
 * This is used for querying the world state.
 * This function is equivalent to calling <code>createTransaction(name).evaluate()</code>.
 * @method Contract#evaluateTransaction
 * @memberof module:fabric-network
 * @param {string} name Transaction function name.
 * @param {...string} [args] Transaction function arguments.
 * @returns {Buffer} Payload response from the transaction function.
 */
/**
 * Add a listener to receive all chaincode events emitted by the smart contract as part of successfully committed
 * transactions. The default is to listen for full contract events from the current block position.
 * @method Contract#addContractListener
 * @memberof module:fabric-network
 * @param {module:fabric-network.ContractListener} listener A contract listener callback function.
 * @param {module:fabric-network.ListenerOptions} [options] Listener options.
 * @returns {Promise<module:fabric-network.ContractListener>} The added listener.
 * @example
 * const listener: ContractListener = async (event) => {
 *     if (event.eventName === 'newOrder') {
 *         const details = event.payload.toString('utf8');
 *         // Run business process to handle orders
 *     }
 * };
 * contract.addContractListener(listener);
 */
/**
 * Remove a previously added contract listener.
 * @method Contract#removeContractListener
 * @memberof module:fabric-network
 * @param {module:fabric-network.ContractListener} listener A contract listener callback function.
 */
/**
 * Provide a Discovery Interest settings to help the peer's discovery service
 * build an endorsement plan. This chaincode Id will be include by default in
 * the list of discovery interests. If this contract's chaincode is in one or
 * more collections then use this method with this chaincode Id to change the
 * default discovery interest to include those collection names.
 * @method Contract#addDiscoveryInterest
 * @memberof module:fabric-network
 * @param {DiscoveryInterest} interest - These will be added to the existing discovery interests and used when
 * {@link module:fabric-network.Transaction#submit} is called.
 * @return {Contract} This Contract instance
 */
/**
 * reset Discovery interest to default of this contracts chaincode name
 * and no collection names and no other chaincode names.
 * @method Contract#resetDiscoveryInterests
 * @memberof module:fabric-network
 * @return {Contract} This Contract instance
 */
/**
 * Retrieve the Discovery Interest settings that will help the peer's
 * discovery service build an endorsement plan.
 * @method Contract#getDiscoveryInterests
 * @memberof module:fabric-network
 * @return {DiscoveryInterest[]} - An array of DiscoveryInterest
 */
/**
 * A callback function that will be invoked when a block event is received.
 * @callback ContractListener
 * @memberof module:fabric-network
 * @async
 * @param {module:fabric-network.ContractEvent} event Contract event.
 * @returns {Promise<void>}
 */
export declare class ContractImpl {
    readonly chaincodeId: string;
    readonly namespace: string;
    readonly network: NetworkImpl;
    readonly gateway: Gateway;
    private discoveryService?;
    private readonly contractListeners;
    private discoveryInterests;
    private discoveryResultsListners;
    constructor(network: NetworkImpl, chaincodeId: string, namespace: string);
    createTransaction(name: string): Transaction;
    submitTransaction(name: string, ...args: string[]): Promise<Buffer>;
    evaluateTransaction(name: string, ...args: string[]): Promise<Buffer>;
    addContractListener(listener: ContractListener, options?: ListenerOptions): Promise<ContractListener>;
    removeContractListener(listener: ContractListener): void;
    /**
     * Internal use
     * Use this method to get the DiscoveryHandler to get the endorsements
     * needed to commit a transaction.
     * The first time this method is called, this contract's DiscoveryService
     * instance will be setup.
     * The service will make a discovery request to the same
     * target as that used by the Network. The request will include this contract's
     * discovery interests. This will enable the peer's discovery
     * service to generate an endorsement plan based on the chaincode's
     * endorsement policy, the collection configuration, and the current active
     * peers.
     * Note: It is assumed that the discovery interests will not
     * change on successive calls. The handler's DiscoveryService will use the
     * "refreshAge" discovery option after the first call to determine if the
     * endorsement plan should be refreshed by a new call to the peer's
     * discovery service.
     * @private
     * @return {DiscoveryHandler} The handler that will work with the discovery
     * endorsement plan to send a proposal to be endorsed to the peers as described
     * in the plan.
     */
    getDiscoveryHandler(): Promise<DiscoveryHandler | undefined>;
    waitDiscoveryResults(): Promise<unknown>;
    registerDiscoveryResultsListener(callback: any): void;
    notifyDiscoveryResultsListeners(hasDiscoveryResults: boolean): void;
    addDiscoveryInterest(interest: DiscoveryInterest): Contract;
    resetDiscoveryInterests(): Contract;
    getDiscoveryInterests(): DiscoveryInterest[];
    private _getQualifiedName;
}