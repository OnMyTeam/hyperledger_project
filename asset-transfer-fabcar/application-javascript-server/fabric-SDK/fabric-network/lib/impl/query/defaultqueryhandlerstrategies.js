"use strict";
/*
 * Copyright 2018, 2019 IBM All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.PREFER_MSPID_SCOPE_ROUND_ROBIN = exports.PREFER_MSPID_SCOPE_SINGLE = exports.MSPID_SCOPE_ROUND_ROBIN = exports.MSPID_SCOPE_SINGLE = void 0;
const singlequeryhandler_1 = require("./singlequeryhandler");
const roundrobinqueryhandler_1 = require("./roundrobinqueryhandler");
function getOrganizationPeers(network) {
    const mspId = network.getGateway().getIdentity().mspId;
    return network.getChannel().getEndorsers(mspId);
}
function getNetworkPeers(network) {
    return network.getChannel().getEndorsers();
}
/**
 * @typedef DefaultQueryHandlerStrategies
 * @memberof module:fabric-network
 * @property {module:fabric-network.QueryHandlerFactory} MSPID_SCOPE_SINGLE Query any one of the peers for the connected organization. Continue
 * to use the same event service for all queries unless it fails. If the client identity's organization has no peers, this strategy will fail.
 * @property {module:fabric-network.QueryHandlerFactory} MSPID_SCOPE_ROUND_ROBIN Query any one of the peers for the connected organization.
 * Use the next available peer for each successive query. If the client identity's organization has no peers, this strategy will fail.
 * @property {module:fabric-network.QueryHandlerFactory} PREFER_MSPID_SCOPE_SINGLE Query any one of the peers for the connected organization. If the
 * connected organization has no peers, query any one of the peers in the network. Continue to use the same event service for all queries unless it
 * fails.
 * @property {module:fabric-network.QueryHandlerFactory} PREFER_MSPID_SCOPE_ROUND_ROBIN Query any one of the peers for the connected organization. If
 * the connected organization has no peers, query any one of the peers in the network. Use the next available peer for each successive query.
 */
exports.MSPID_SCOPE_SINGLE = (network) => {
    const peers = getOrganizationPeers(network);
    return new singlequeryhandler_1.SingleQueryHandler(peers);
};
exports.MSPID_SCOPE_ROUND_ROBIN = (network) => {
    const peers = getOrganizationPeers(network);
    return new roundrobinqueryhandler_1.RoundRobinQueryHandler(peers);
};
exports.PREFER_MSPID_SCOPE_SINGLE = (network) => {
    let peers = getOrganizationPeers(network);
    if (peers.length === 0) {
        peers = getNetworkPeers(network);
    }
    return new singlequeryhandler_1.SingleQueryHandler(peers);
};
exports.PREFER_MSPID_SCOPE_ROUND_ROBIN = (network) => {
    let peers = getOrganizationPeers(network);
    if (peers.length === 0) {
        peers = getNetworkPeers(network);
    }
    return new roundrobinqueryhandler_1.RoundRobinQueryHandler(peers);
};
//# sourceMappingURL=defaultqueryhandlerstrategies.js.map