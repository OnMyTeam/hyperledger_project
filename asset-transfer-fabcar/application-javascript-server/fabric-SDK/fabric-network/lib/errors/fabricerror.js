/**
 * Copyright 2018 IBM All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */
'use strict';
Object.defineProperty(exports, "__esModule", { value: true });
exports.FabricError = void 0;
/**
 * Base type for Fabric-specific errors.
 * @memberof module:fabric-network
 * @property {Error} [cause] Underlying error that caused this error.
 * @property {string} [transactionId] ID of the associated transaction.
 */
class FabricError extends Error {
    /*
     * Constructor.
     * @param {(string|object)} [info] Either an error message (string) or additional properties to assign to this
     * inctance (object).
     */
    constructor(info) {
        if (!info) {
            super();
        }
        else if (typeof info === 'string') {
            super(info);
        }
        else {
            super(info.message);
            Object.assign(this, info);
        }
        this.name = FabricError.name;
    }
}
exports.FabricError = FabricError;
//# sourceMappingURL=fabricerror.js.map