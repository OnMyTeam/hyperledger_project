// ExpressJS Setup
const express = require('express');
const app = express();
var bodyParser = require('body-parser');
app.use(bodyParser.json());

// Constants
const PORT = 8000;
const HOST = "0.0.0.0";

'use strict';

const { Gateway, Wallets } = require('./fabric-SDK/fabric-network');
const path = require('path');
const fs = require('fs');

const channelName = 'mychannel';
const chaincodeName = 'fabcar';
const adminUserId = 'admin';
const adminUserPasswd = 'adminpw';
const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet');
const org1UserId = 'appUser';

function prettyJSONString(inputString) {
	return JSON.stringify(JSON.parse(inputString), null, 2);
}


async function main() {
    try {


        // Build an in memory object with the network configuration (also known as a connection profile).
        const ccpPath = path.resolve(__dirname, '..', '..', 'test-network', 'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        const fileExists = fs.existsSync(ccpPath);
        if (!fileExists) {
            throw new Error(`no such file or directory: ${ccpPath}`);
        }
        const contents = fs.readFileSync(ccpPath, 'utf8');
        const ccp = JSON.parse(contents);
        console.log(`\n=> Loaded the network configuration located at ${ccpPath}`);

        // Prepare the identity from the wallet.
        let wallet;
        if (walletPath) {
            wallet = await Wallets.newFileSystemWallet(walletPath);
            console.log(`=> Found the file system wallet at ${walletPath}`);
        } else {
            wallet = await Wallets.newInMemoryWallet();
            console.log('=> Found the in-memory wallet');
        }        

        // Check to see if we've already enrolled the user.
	    const userIdentity = await wallet.get(org1UserId);
	    if (!userIdentity) {
		    console.log(`An identity for the user ${org1UserId} does not exist in the wallet`);
		    return;
        }
        console.log(`=> The user is ${org1UserId}`);

        // Setup the gateway object.
        // The user will now be able to create connections to the fabric network and be able to
        // submit transactions and query. All transactions submitted by this gateway will be
        // signed by this user using the credentials stored in the wallet.
        const gateway = new Gateway();
        await gateway.connect(ccp, {
            wallet,
            identity: org1UserId,
            discovery: { enabled: true, asLocalhost: true } // using asLocalhost as this gateway is using a fabric network deployed locally
        });
        console.log(`=> Gateway set up`);

        // Build a network object based on the channel where the smart contract is deployed.
        const network = await gateway.getNetwork(channelName);
        console.log('=> Channel obtained');

        // Get the contract object from the network.
        const contract = network.getContract(chaincodeName);
        console.log('=> Contract obtained');

        // Invoke the chaincode function!!!
        // Now let's try to submit a transaction.
		// This will be sent to both peers and if both peers endorse the transaction, the endorsed proposal will be sent
		// to the orderer to be committed by each of the peer's to the channel ledger.
		console.log('=> Submit Transaction: AddCar, adds new car with id, make, model, colour, and owner arguments');
		await contract.submitBWTransaction('AAAA', 'BBBB', "CCCC", 11, 22, 333, 444);   
        
        console.log('=> Transaction has been submitted');
        await gateway.disconnect();

        

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
       
    }   

};

main()