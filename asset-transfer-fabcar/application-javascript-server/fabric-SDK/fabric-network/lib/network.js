"use strict";
/*
 * Copyright 2018, 2019 IBM All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.NetworkImpl = void 0;
const contract_1 = require("./contract");
const blockeventsource_1 = require("./impl/event/blockeventsource");
const commitlistenersession_1 = require("./impl/event/commitlistenersession");
const eventservicemanager_1 = require("./impl/event/eventservicemanager");
const isolatedblocklistenersession_1 = require("./impl/event/isolatedblocklistenersession");
const listeners_1 = require("./impl/event/listeners");
const listenersession_1 = require("./impl/event/listenersession");
const sharedblocklistenersession_1 = require("./impl/event/sharedblocklistenersession");
const gatewayutils_1 = require("./impl/gatewayutils");
const Logger = require("./logger");
const logger = Logger.getLogger('Network');
async function listenerOptionsWithDefaults(options) {
    var _a;
    const defaultOptions = {
        type: 'full'
    };
    const result = Object.assign(defaultOptions, options);
    const checkpointBlock = await ((_a = options.checkpointer) === null || _a === void 0 ? void 0 : _a.getBlockNumber());
    if (checkpointBlock) {
        result.startBlock = checkpointBlock;
    }
    return result;
}
/**
 * <p>A Network represents the set of peers in a Fabric network.
 * Applications should get a Network instance using the
 * gateway's [getNetwork]{@link module:fabric-network.Gateway#getNetwork} method.</p>
 *
 * <p>The Network object provides the ability for applications to:</p>
 * <ul>
 *   <li>Obtain a specific smart contract deployed to the network using [getContract]{@link module:fabric-network.Network#getContract},
 *       in order to submit and evaluate transactions for that smart contract.</li>
 *   <li>Listen to new block events and replay previous block events using
 *       [addBlockListener]{@link module:fabric-network.Network#addBlockListener}.</li>
 * </ul>
 * @interface Network
 * @memberof module:fabric-network
 */
/**
 * Get the owning Gateway connection.
 * @method Network#getGateway
 * @memberof module:fabric-network
 * @returns {module:fabric-network.Gateway} A Gateway.
 */
/**
 * Get an instance of a contract (chaincode) on the current network.
 * @method Network#getContract
 * @memberof module:fabric-network
 * @param {string} chaincodeId - the chaincode identifier.
 * @param {string} [name] - the name of the contract.
 * @param {string[]} [collections] - the names of collections defined for this chaincode.
 * @returns {module:fabric-network.Contract} the contract.
 */
/**
 * Get the underlying channel object representation of this network.
 * @method Network#getChannel
 * @memberof module:fabric-network
 * @returns {Channel} A channel.
 */
/**
 * Add a listener to receive transaction commit and peer disconnect events for a set of peers. This is typically used
 * only within the implementation of a custom [transaction commit event handler]{@tutorial transaction-commit-events}.
 * @method Network#addCommitListener
 * @memberof module:fabric-network
 * @param {module:fabric-network.CommitListener} listener A transaction commit listener callback function.
 * @param {Endorser[]} peers The peers from which to receive events.
 * @param {string} transactionId A transaction ID.
 * @returns {module:fabric-network.CommitListener} The added listener.
 * @example
 * const listener: CommitListener = (error, event) => {
 *     if (error) {
 *         // Handle peer communication error
 *     } else {
 *         // Handle transaction commit event
 *     }
 * }
 * const peers = network.channel.getEndorsers();
 * await network.addCommitListener(listener, peers, transactionId);
 */
/**
 * Remove a previously added transaction commit listener.
 * @method Network#removeCommitListener
 * @memberof module:fabric-network
 * @param {module:fabric-network.CommitListener} listener A transaction commit listener callback function.
 */
/**
 * Add a listener to receive block events for this network. Blocks will be received in order and without duplication.
 * The default is to listen for full block events from the current block position.
 * @method Network#addBlockListener
 * @memberof module:fabric-network
 * @async
 * @param {module:fabric-network.BlockListener} listener A block listener callback function.
 * @param {module:fabric-network.ListenerOptions} [options] Listener options.
 * @returns {Promise<module:fabric-network.BlockListener>} The added listener.
 * @example
 * const listener: BlockListener = async (event) => {
 *     // Handle block event
 *
 *     // Listener may remove itself if desired
 *     if (event.blockNumber.equals(endBlock)) {
 *         network.removeBlockListener(listener);
 *     }
 * }
 * const options: ListenerOptions = {
 *     startBlock: 1
 * };
 * await network.addBlockListener(listener, options);
 */
/**
 * Remove a previously added block listener.
 * @method Network#removeBlockListener
 * @memberof module:fabric-network
 * @param listener {module:fabric-network.BlockListener} A block listener callback function.
 */
/**
 * A callback function that will be invoked when a block event is received.
 * @callback BlockListener
 * @memberof module:fabric-network
 * @async
 * @param {module:fabric-network.BlockEvent} event A block event.
 * @returns {Promise<void>}
 */
/**
 * A callback function that will be invoked when either a peer communication error occurs or a transaction commit event
 * is received. Only one of the two arguments will have a value for any given invocation.
 * @callback CommitListener
 * @memberof module:fabric-network
 * @param {module:fabric-network.CommitError} [error] Peer communication error.
 * @param {module:fabric-network.CommitEvent} [event] Transaction commit event from a specific peer.
 */
/**
 * @interface CommitError
 * @extends Error
 * @memberof module:fabric-network
 * @property {Endorser} peer The peer that raised this error.
 */
/**
 * @interface CommitEvent
 * @extends {module:fabric-network.TransactionEvent}
 * @memberof module:fabric-network
 * @property {Endorser} peer The endorsing peer that produced this event.
 */
class NetworkImpl {
    /*
     * Network constructor for internal use only.
     * @param {Gateway} gateway The owning gateway instance
     * @param {Channel} channel The fabric-common channel instance
     */
    constructor(gateway, channel) {
        this.contracts = new Map();
        this.initialized = false;
        this.commitListeners = new Map();
        this.blockListeners = new Map();
        const method = 'constructor';
        logger.debug('%s - start', method);
        this.gateway = gateway;
        this.channel = channel;
        this.eventServiceManager = new eventservicemanager_1.EventServiceManager(this);
        this.realtimeFilteredBlockEventSource = new blockeventsource_1.BlockEventSource(this.eventServiceManager, { type: 'filtered' });
        this.realtimeFullBlockEventSource = new blockeventsource_1.BlockEventSource(this.eventServiceManager, { type: 'full' });
        this.realtimePrivateBlockEventSource = new blockeventsource_1.BlockEventSource(this.eventServiceManager, { type: 'private' });
    }
    getGateway() {
        return this.gateway;
    }
    getContract(chaincodeId, name = '') {
        const method = 'getContract';
        logger.debug('%s - start - name %s', method, name);
        if (!this.initialized) {
            throw new Error('Unable to get contract as this network has failed to initialize');
        }
        const key = `${chaincodeId}:${name}`;
        let contract = this.contracts.get(key);
        if (!contract) {
            contract = new contract_1.ContractImpl(this, chaincodeId, name);
            logger.debug('%s - create new contract %s', method, chaincodeId);
            this.contracts.set(key, contract);
        }
        return contract;
    }
    getChannel() {
        return this.channel;
    }
    async addCommitListener(listener, peers, transactionId) {
        const sessionSupplier = async () => new commitlistenersession_1.CommitListenerSession(listener, this.eventServiceManager, peers, transactionId);
        return await listenersession_1.addListener(listener, this.commitListeners, sessionSupplier);
    }
    removeCommitListener(listener) {
        listenersession_1.removeListener(listener, this.commitListeners);
    }
    async addBlockListener(listener, options = {}) {
        const sessionSupplier = async () => await this.newBlockListenerSession(listener, options);
        return await listenersession_1.addListener(listener, this.blockListeners, sessionSupplier);
    }
    removeBlockListener(listener) {
        listenersession_1.removeListener(listener, this.blockListeners);
    }
    _dispose() {
        const method = '_dispose';
        logger.debug('%s - start', method);
        this.contracts.clear();
        this.commitListeners.forEach((listener) => listener.close());
        this.commitListeners.clear();
        this.blockListeners.forEach((listener) => listener.close());
        this.blockListeners.clear();
        this.realtimeFilteredBlockEventSource.close();
        this.realtimeFullBlockEventSource.close();
        this.realtimePrivateBlockEventSource.close();
        this.eventServiceManager.close();
        this.channel.close();
        this.initialized = false;
    }
    /**
     * Initialize this network instance
     * @private
     */
    async _initialize(discover) {
        const method = '_initialize';
        logger.debug('%s - start', method);
        if (this.initialized) {
            return;
        }
        await this._initializeInternalChannel(discover);
        this.initialized = true;
        // Must be created after channel initialization to ensure discovery has located the peers
        const queryOptions = this.gateway.getOptions().queryHandlerOptions;
        this.queryHandler = queryOptions.strategy(this);
        logger.debug('%s - end', method);
    }
    /**
     * initialize the channel if it hasn't been done
     * @private
     */
    async _initializeInternalChannel(options) {
        const method = '_initializeInternalChannel';
        logger.debug('%s - start', method);
        if (options === null || options === void 0 ? void 0 : options.enabled) {
            logger.debug('%s - initialize with discovery', method);
            let targets;
            logger.debug('%s - user has not specified discovery targets, check channel and client', method);
            // maybe the channel has connected endorsers with the mspid
            const mspId = this.gateway.getIdentity().mspId;
            targets = this.channel.getEndorsers(mspId);
            if (!targets || targets.length < 1) {
                // then check the client for connected peers associated with the mspid
                targets = this.channel.client.getEndorsers(mspId);
            }
            if (!targets || targets.length < 1) {
                // get any peer
                targets = this.channel.client.getEndorsers();
            }
            if (!targets || targets.length < 1) {
                throw Error('No discovery targets found');
            }
            else {
                logger.debug('%s - using channel/client targets', method);
            }
            // should have targets by now, create the discoverers from the endorsers
            const discoverers = [];
            for (const peer of targets) {
                const discoverer = this.channel.client.newDiscoverer(peer.name, peer.mspid);
                discoverer.setEndpoint(peer.endpoint);
                discoverers.push(discoverer);
            }
            this.discoveryService = this.channel.newDiscoveryService(this.channel.name);
            const idx = this.gateway.identityContext;
            // do the three steps
            this.discoveryService.build(idx);
            this.discoveryService.sign(idx);
            logger.debug('%s - will discover asLocalhost:%s', method, options.asLocalhost);
            await this.discoveryService.send({
                asLocalhost: options.asLocalhost,
                targets: discoverers
            });
            // now we can work with the discovery results
            // or get a handler later from the discoverService
            // to be used on endorsement, queries, and commits
            logger.debug('%s - discovery complete - channel is populated', method);
        }
        logger.debug('%s - end', method);
    }
    async newBlockListenerSession(listener, options) {
        options = await listenerOptionsWithDefaults(options);
        if (options.checkpointer) {
            listener = listeners_1.checkpointBlockListener(listener, options.checkpointer);
        }
        if (gatewayutils_1.notNullish(options.startBlock)) {
            return this.newIsolatedBlockListenerSession(listener, options);
        }
        else {
            return this.newSharedBlockListenerSession(listener, options.type);
        }
    }
    newIsolatedBlockListenerSession(listener, options) {
        const blockSource = new blockeventsource_1.BlockEventSource(this.eventServiceManager, options);
        return new isolatedblocklistenersession_1.IsolatedBlockListenerSession(listener, blockSource);
    }
    newSharedBlockListenerSession(listener, type) {
        if (type === 'filtered') {
            return new sharedblocklistenersession_1.SharedBlockListenerSession(listener, this.realtimeFilteredBlockEventSource);
        }
        else if (type === 'full') {
            return new sharedblocklistenersession_1.SharedBlockListenerSession(listener, this.realtimeFullBlockEventSource);
        }
        else if (type === 'private') {
            return new sharedblocklistenersession_1.SharedBlockListenerSession(listener, this.realtimePrivateBlockEventSource);
        }
        else {
            throw new Error('Unsupported event listener type: ' + type);
        }
    }
}
exports.NetworkImpl = NetworkImpl;
//# sourceMappingURL=network.js.map