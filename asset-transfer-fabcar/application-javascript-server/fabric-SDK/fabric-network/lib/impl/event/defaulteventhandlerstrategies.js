"use strict";
/*
 * Copyright 2018, 2019 IBM All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.NONE = exports.PREFER_MSPID_SCOPE_ANYFORTX = exports.PREFER_MSPID_SCOPE_ALLFORTX = exports.NETWORK_SCOPE_ANYFORTX = exports.NETWORK_SCOPE_ALLFORTX = exports.MSPID_SCOPE_ANYFORTX = exports.MSPID_SCOPE_ALLFORTX = void 0;
const allfortxstrategy_1 = require("./allfortxstrategy");
const anyfortxstrategy_1 = require("./anyfortxstrategy");
const transactioneventhandler_1 = require("./transactioneventhandler");
function getOrganizationPeers(network) {
    const mspId = network.getGateway().getIdentity().mspId;
    return network.getChannel().getEndorsers(mspId);
}
function getNetworkPeers(network) {
    return network.getChannel().getEndorsers();
}
/**
 * @typedef DefaultEventHandlerStrategies
 * @memberof module:fabric-network
 * @property {module:fabric-network.TxEventHandlerFactory} MSPID_SCOPE_ALLFORTX Listen for transaction commit
 * events from all peers in the client identity's organization. If the client identity's organization has no peers,
 * this strategy will fail.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will wait until successful
 * events are received from <em>all</em> currently connected peers (minimum 1).
 * @property {module:fabric-network.TxEventHandlerFactory} MSPID_SCOPE_ANYFORTX Listen for transaction commit
 * events from all peers in the client identity's organization. If the client identity's organization has no peers,
 * this strategy will fail.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will wait until a successful
 * event is received from <em>any</em> peer.
 * @property {module:fabric-network.TxEventHandlerFactory} NETWORK_SCOPE_ALLFORTX Listen for transaction commit
 * events from all peers in the network.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will wait until successful
 * events are received from <em>all</em> currently connected peers (minimum 1).
 * @property {module:fabric-network.TxEventHandlerFactory} NETWORK_SCOPE_ANYFORTX Listen for transaction commit
 * events from all peers in the network.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will wait until a
 * successful event is received from <em>any</em> peer.
 * * @property {module:fabric-network.TxEventHandlerFactory} PREFER_MSPID_SCOPE_ALLFORTX Listen for transaction commit
 * events from all peers in the client identity's organization. If the client identity's organization has no peers, listen
 * for transaction commit events from all peers in the network.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will wait until successful
 * events are received from <em>all</em> currently connected peers (minimum 1).
 * @property {module:fabric-network.TxEventHandlerFactory} PREFER_MSPID_SCOPE_ANYFORTX Listen for transaction commit
 * events from all peers in the client identity's organization. If the client identity's organization has no peers, listen
 * for transaction commit events from all peers in the network.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will wait until a
 * successful event is received from <em>any</em> peer.
 * @property {module:fabric-network.TxEventHandlerFactory} NONE Do not wait for any commit events.
 * The [submitTransaction]{@link module:fabric-network.Contract#submitTransaction} function will return immediately
 * after successfully sending the transaction to the orderer.
 */
exports.MSPID_SCOPE_ALLFORTX = (transactionId, network) => {
    const eventStrategy = new allfortxstrategy_1.AllForTxStrategy(getOrganizationPeers(network));
    return new transactioneventhandler_1.TransactionEventHandler(transactionId, network, eventStrategy);
};
exports.MSPID_SCOPE_ANYFORTX = (transactionId, network) => {
    const eventStrategy = new anyfortxstrategy_1.AnyForTxStrategy(getOrganizationPeers(network));
    return new transactioneventhandler_1.TransactionEventHandler(transactionId, network, eventStrategy);
};
exports.NETWORK_SCOPE_ALLFORTX = (transactionId, network) => {
    const eventStrategy = new allfortxstrategy_1.AllForTxStrategy(getNetworkPeers(network));
    return new transactioneventhandler_1.TransactionEventHandler(transactionId, network, eventStrategy);
};
exports.NETWORK_SCOPE_ANYFORTX = (transactionId, network) => {
    const eventStrategy = new anyfortxstrategy_1.AnyForTxStrategy(getNetworkPeers(network));
    return new transactioneventhandler_1.TransactionEventHandler(transactionId, network, eventStrategy);
};
exports.PREFER_MSPID_SCOPE_ALLFORTX = (transactionId, network) => {
    let peers = getOrganizationPeers(network);
    if (peers.length === 0) {
        peers = getNetworkPeers(network);
    }
    const eventStrategy = new allfortxstrategy_1.AllForTxStrategy(peers);
    return new transactioneventhandler_1.TransactionEventHandler(transactionId, network, eventStrategy);
};
exports.PREFER_MSPID_SCOPE_ANYFORTX = (transactionId, network) => {
    let peers = getOrganizationPeers(network);
    if (peers.length === 0) {
        peers = getNetworkPeers(network);
    }
    const eventStrategy = new anyfortxstrategy_1.AnyForTxStrategy(peers);
    return new transactioneventhandler_1.TransactionEventHandler(transactionId, network, eventStrategy);
};
const noOpEventHandler = {
    startListening: async () => {
        // No-op
    },
    waitForEvents: async () => {
        // No-op
    },
    cancelListening: () => {
        // No-op
    }
};
exports.NONE = (transactionId, network) => {
    return noOpEventHandler;
};
//# sourceMappingURL=defaulteventhandlerstrategies.js.map