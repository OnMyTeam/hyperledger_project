/**
 * Copyright 2020 IBM All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
'use strict';
Object.defineProperty(exports, "__esModule", { value: true });
exports.TransactionError = void 0;
const fabricerror_1 = require("./fabricerror");
/**
 * Base type for Fabric-specific errors.
 * @memberof module:fabric-network
 * @property {string} [transactionId] ID of the associated transaction.
 * @property {string} [transactionCode] The transaction validation code of the associated transaction.
 */
class TransactionError extends fabricerror_1.FabricError {
    /*
     * Constructor.
     * @param {(string|object)} [info] Either an error message (string) or additional properties to assign to this
     * instance (object).
     */
    constructor(info) {
        super(info);
        this.name = TransactionError.name;
    }
}
exports.TransactionError = TransactionError;
//# sourceMappingURL=transactionerror.js.map