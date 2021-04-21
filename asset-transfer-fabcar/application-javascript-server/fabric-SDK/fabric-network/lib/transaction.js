"use strict";
/*
 * Copyright 2018, 2019 IBM All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Transaction = void 0;
const query_1 = require("./impl/query/query");
const EventHandlers = require("./impl/event/defaulteventhandlerstrategies");
const Logger = require("./logger");
const util = require("util");
const logger = Logger.getLogger('Transaction');
const path = require('path');
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader');
function getResponsePayload(proposalResponse) {
    const validEndorsementResponse = getValidEndorsementResponse(proposalResponse.responses);
    if (!validEndorsementResponse) {
        const error = newEndorsementError(proposalResponse);
        logger.error(error);
        throw error;
    }
    return validEndorsementResponse.response.payload;
}
function getValidEndorsementResponse(endorsementResponses) {
    return endorsementResponses.find((endorsementResponse) => endorsementResponse.endorsement);
}
function newEndorsementError(proposalResponse) {
    var _a, _b;
    const errorInfos = [];
    for (const error of proposalResponse.errors) {
        const errorInfo = {
            peer: (_a = error === null || error === void 0 ? void 0 : error.connection) === null || _a === void 0 ? void 0 : _a.name,
            status: 'grpc',
            message: error.message
        };
        errorInfos.push(errorInfo);
    }
    for (const endorsement of proposalResponse.responses) {
        const errorInfo = {
            peer: (_b = endorsement === null || endorsement === void 0 ? void 0 : endorsement.connection) === null || _b === void 0 ? void 0 : _b.name,
            status: endorsement.response.status,
            message: endorsement.response.message
        };
        errorInfos.push(errorInfo);
    }
    const messages = ['No valid responses from any peers. Errors:'];
    for (const errorInfo of errorInfos) {
        messages.push(util.format('peer=%s, status=%s, message=%s', errorInfo.peer, errorInfo.status, errorInfo.message));
    }
    return new Error(messages.join('\n    '));
}
function isInteger(value) {
    return Number.isInteger(value);
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
class Transaction {
    /*
     * @param {Contract} contract Contract to which this transaction belongs.
     * @param {String} name Fully qualified transaction name.
     * @param {function} eventStrategyFactory - A factory function that will return
     * an EventStrategy.
     */
    constructor(contract, name) {
        const method = `constructor[${name}]`;
        logger.debug('%s - start', method);
        this.contract = contract;
        this.name = name;
        this.gatewayOptions = contract.gateway.getOptions();
        this.eventHandlerStrategyFactory = this.gatewayOptions.eventHandlerOptions.strategy || EventHandlers.NONE;
        this.queryHandler = contract.network.queryHandler;
        // Store the returned copy to prevent state being modified by other code before it is used to send proposals
        this.identityContext = this.contract.gateway.identityContext.calculateTransactionId();
        this.transactionId = this.identityContext.transactionId;
    }
    /**
     * Get the fully qualified name of the transaction function.
     * @returns {string} Transaction name.
     */
    getName() {
        return this.name;
    }
    /**
     * Set transient data that will be passed to the transaction function
     * but will not be stored on the ledger. This can be used to pass
     * private data to a transaction function.
     * @param {Object} transientMap Object with String property names and
     * Buffer property values.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setTransient(transientMap) {
        const method = `setTransient[${this.name}]`;
        logger.debug('%s - start', method);
        this.transientMap = transientMap;
        return this;
    }
    /**
     * Get the ID that will be used for this transaction invocation.
     * @returns {string} A transaction ID.
     */
    getTransactionId() {
        return this.transactionId;
    }
    /**
     * Set the peers that should be used for endorsement when this transaction
     * is submitted to the ledger.
     * Setting the peers will override the use of discovery and the submit will
     * send the proposal to these peers.
     * This will override the setEndorsingOrganizations if previously called.
     * @param {Endorser[]} peers - Endorsing peers.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setEndorsingPeers(peers) {
        const method = `setEndorsingPeers[${this.name}]`;
        logger.debug('%s - start', method);
        this.endorsingPeers = peers;
        this.endorsingOrgs = undefined;
        return this;
    }
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
    setEndorsingOrganizations(...orgs) {
        const method = `setEndorsingOrganizations[${this.name}]`;
        logger.debug('%s - start', method);
        this.endorsingOrgs = orgs;
        this.endorsingPeers = undefined;
        return this;
    }
    /**
     * Set an event handling strategy to use for this transaction instead of the default configured on the gateway.
     * @param strategy An event handling strategy.
     * @returns {module:fabric-network.Transaction} This object, to allow function chaining.
     */
    setEventHandler(strategy) {
        this.eventHandlerStrategyFactory = strategy;
        return this;
    }
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
    async submit(...args) {
        const method = `submit[${this.name}]`;
        logger.debug('%s - start', method);
        const channel = this.contract.network.getChannel();
        const transactionOptions = this.gatewayOptions.eventHandlerOptions;
        // This is the object that will centralize this endorsement activities
        // with the fabric network
        const endorsement = channel.newEndorsement(this.contract.chaincodeId);
        const proposalBuildRequest = this.newBuildProposalRequest(args);
        logger.debug('%s - build and send the endorsement', method);
        // build the outbound request along with getting a new transactionId
        // from the identity context
        endorsement.build(this.identityContext, proposalBuildRequest);
        endorsement.sign(this.identityContext);
        // ------- S E N D   P R O P O S A L
        // This is where the request gets sent to the peers
        const proposalSendRequest = {};
        if (isInteger(transactionOptions.endorseTimeout)) {
            proposalSendRequest.requestTimeout = transactionOptions.endorseTimeout * 1000; // in ms;
        }
        if (this.endorsingPeers) {
            logger.debug('%s - user has assigned targets', method);
            proposalSendRequest.targets = this.endorsingPeers;
        }
        else if (this.contract.network.discoveryService) {
            logger.debug('%s - discovery handler will be used for endorsing', method);
            proposalSendRequest.handler = await this.contract.getDiscoveryHandler();
            if (this.endorsingOrgs) {
                logger.debug('%s - using discovery and user has assigned endorsing orgs %s', method, this.endorsingOrgs);
                proposalSendRequest.requiredOrgs = this.endorsingOrgs;
            }
        }
        else if (this.endorsingOrgs) {
            logger.debug('%s - user has assigned endorsing orgs %s', method, this.endorsingOrgs);
            const flatten = (accumulator, value) => {
                accumulator.push(...value);
                return accumulator;
            };
            proposalSendRequest.targets = this.endorsingOrgs.map((mspid) => channel.getEndorsers(mspid)).reduce(flatten, []);
        }
        else {
            logger.debug('%s - targets will default to all that are assigned to this channel', method);
            proposalSendRequest.targets = channel.getEndorsers();
        }
        // by now we should have targets or a discovery handler to be used
        // by the send() of the proposal instance
        
        const proposalResponse = await endorsement.send(proposalSendRequest);
        try {
            const result = getResponsePayload(proposalResponse);
            // ------- E V E N T   M O N I T O R
            const eventHandler = this.eventHandlerStrategyFactory(endorsement.getTransactionId(), this.contract.network);
            await eventHandler.startListening();
            const commit = endorsement.newCommit();
            commit.build(this.identityContext);
            commit.sign(this.identityContext);
            // -----  C O M M I T   E N D O R S E M E N T
            // this is where the endorsement results are sent to the orderer
            const commitSendRequest = {};
            if (isInteger(transactionOptions.commitTimeout)) {
                commitSendRequest.requestTimeout = transactionOptions.commitTimeout * 1000; // in ms;
            }
            if (proposalSendRequest.handler) {
                logger.debug('%s - use discovery to commit', method);
                commitSendRequest.handler = proposalSendRequest.handler;
            }
            else {
                logger.debug('%s - use the orderers assigned to the channel', method);
                commitSendRequest.targets = channel.getCommitters();
            }
            // by now we should have a discovery handler or use the target orderers
            // that have been assigned from the channel to perform the commit
            const commitResponse = await commit.send(commitSendRequest);
            logger.debug('%s - commit response %j', method, commitResponse);
            if (commitResponse.status !== 'SUCCESS') {
                const msg = `Failed to commit transaction %${endorsement.getTransactionId()}, orderer response status: ${commitResponse.status}`;
                logger.error('%s - %s', method, msg);
                eventHandler.cancelListening();
                throw new Error(msg);
            }
            else {
                logger.debug('%s - successful commit', method);
            }
            logger.debug('%s - wait for the transaction to be committed on the peer', method);
            await eventHandler.waitForEvents();
            return result;
        }
        catch (err) {
            err.responses = proposalResponse.responses;
            err.errors = proposalResponse.errors;
            throw err;
        }
    }

    async submitAggregator(...args) {
        const PROTO_PATH = path.resolve(__dirname , '../../../../../BWAggregator/protos/bwaggregator.proto');

        const packageDefinition = protoLoader.loadSync(
            PROTO_PATH,
            {keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true
            });
        const helloProto = grpc.loadPackageDefinition(packageDefinition).user;
        var client = new helloProto.User('localhost:9000',
                                            grpc.credentials.createInsecure());
        var user;
        if (process.argv.length >= 3) {
            user = process.argv[2];
        } else {
            user = 'world';
        }
        client.processProposal({user_id: "1", name: "sanggi", phone_number:"010", age:11}, function(err, response) {
            console.log('Greeting:', response);
        });        

    }    
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
    async evaluate(...args) {
        const method = `evaluate[${this.name}]`;
        logger.debug('%s - start', method);
        const queryProposal = this.contract.network.getChannel().newQuery(this.contract.chaincodeId);
        const request = this.newBuildProposalRequest(args);
        logger.debug('%s - build and sign the query', method);
        queryProposal.build(this.identityContext, request);
        queryProposal.sign(this.identityContext);
        const query = new query_1.QueryImpl(queryProposal, this.gatewayOptions.queryHandlerOptions);
        logger.debug('%s - handler will send', method);
        const results = await this.queryHandler.evaluate(query);
        logger.debug('%s - queryHandler completed', method);
        return results;
    }
    newBuildProposalRequest(args) {
        const request = {
            fcn: this.name,
            args: args,
            generateTransactionId: false
        };
        if (this.transientMap) {
            request.transientMap = this.transientMap;
        }
        return request;
    }
}
exports.Transaction = Transaction;
//# sourceMappingURL=transaction.js.map