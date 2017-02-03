var hyperledger = {};
hyperledger.URL = 'http://0.0.0.0:7050/';
hyperledger.APPNAME = 'DecodedBlockChain';
hyperledger.BLOCKCHAIN = [];

// Specify the base structure for the options object.
function getOptions(path, requestMethod, functionName, arg, ID) {
    var opt = {
        method: "POST",
        url: hyperledger.URL + path,
        body: {
            jsonrpc: "2.0",
            method: requestMethod,
            params: {
                type: 1, // Golang
                chaincodeID: { name: hyperledger.APPNAME },
                ctorMsg: {
                    function: functionName,
                    args: arg
                }
            },
            id: ID,
        },
        headers: { 
            "Content-Type": "application/json", 
            "User-Agent": "request", // needed?
            "Cache-Control": "no-cache" // needed?
        },
        json: true
    };
    return opt;
} // end of getOptions()

hyperledger.init = function (requestpromise) {
    var optionsCheck = getOptions("chaincode", "query", "read", [ "Parties" ], 1337);
    requestpromise(optionsCheck)
        .then( function (dataCheck) {
            // Even error will come back as a successful api-request.
            if ( 'error' in dataCheck && dataCheck.error.code == -32003 && dataCheck.error.message == 'Query failure' ){
                // Assuming it has not been deployed. If deploying fails it means the blockchain backend is probably down.
                var optionsDeploy = getOptions("chaincode", "deploy", "init", [], 9999);
                requestpromise(optionsDeploy)
                    .then( function (dataDeploy) {
                        // This check does not work... it doesnt get an error message if the chaincode is down.
                        // This only shows the blockchain cluster is up and running, not that the chaincode is running.
                        if ( dataDeploy.result.status == 'OK' && dataDeploy.result.message == hyperledger.APPNAME ) {
                            console.log(' *** Successfully deployed the blockchain app');
                        }
                        else {
                            console.log(' *** Failed to deploy the blockchain app', dataDeploy);
                        }
                    })
                    .catch( function (errDeploy) {
                        console.log(' *** ERROR in deploying the blockchain', err);
                    })
            }
            else if ( 'result' in dataCheck && dataCheck.result.status == 'OK' ) {
                console.log(' *** Blockchain is seems to be ready and deployed.');
            }
        })
        .catch( function (errCheck) {
            console.log(' *** ERROR in checking for blockchain availability: ', errCheck);
        })
}

hyperledger.deploy = function () {
    var opt = getOptions("chaincode", "deploy", "init", [], 9999);
    return opt;
}

hyperledger.createParty = function (id, name, voter, candidate) {
    return getOptions("chaincode", "invoke", "createParty", [ id, name, voter, candidate ], 2600);
}

hyperledger.readParty = function (id) {
    return getOptions("chaincode", "query", "readParty", [ id ], 2600);
}

hyperledger.readAllParties = function () {
    return getOptions("chaincode", "query", "readAllParties", [], 1337);
}

hyperledger.createVotesAndAssignToAll = function () {
    return getOptions("chaincode", "invoke", "createVotesAndAssignToAll", [], 2600)
}

hyperledger.giveVoteToCandidate = function (voteId, candidateId) {
    return getOptions("chaincode", "invoke", "createVotesAndAssignToAll", [ voteId, candidateId ], 2600)
}

hyperledger.getChain = function () {
    return { method: "GET", url: hyperledger.URL + 'chain', json: true };
}

hyperledger.getBlock = function (blockId) {
    return { method: "GET", url: hyperledger.URL + 'chain/blocks/' + blockId, json: true };
}

hyperledger.getFullBlockChain = function (requestpromise) {
    console.log('\t*** Fetching the entire blockchain.');
    var optionsChain = hyperledger.getChain();
    // Return the promise itself.
    return requestpromise(optionsChain)
        .then( function (dataChain) {
            var theBlockChain = [];
            var chainPromise = new Promise( function (resolve, reject) {
                if ( 'height' in dataChain ) {
                    // These promises will not be resolved in order...
                    for ( var b=0; b<dataChain.height; b++ ) {
                        var optB = hyperledger.getBlock(b);
                        requestpromise(optB)
                            .then( function (dataBlock) { 
                                theBlockChain.push(dataBlock);
                                if ( theBlockChain.length == (dataChain.height) ) {
                                    resolve(theBlockChain); // This resolves the entire thing.
                                }
                            })
                            .catch( function (errBlock) {
                                console.log(errBlock);
                            });
                    }
                }
                else {
                    reject('Can\'t get dataChain');
                }
            });                
            return chainPromise;
        })
        .then( function (chain) {
            console.log('\t*** SUCCESS: Fetching the entire blockchain.');
            var chainPromise = new Promise( function (resolve, reject) {
                // Blockchain ready, we now need to decode the payloads.
                for ( var b=0; b<chain.length; b++ ) {
                    var block = chain[b];
                    if ( 'transactions' in block && block.transactions.length > 0 ) {
                        for ( var t=0; t<block.transactions.length; t++ ) {
                            var tx = block.transactions[t];
                            tx.payloadDecoded = hyperledger.base64ToString(tx.payload);
                            chain[b].transactions[t] = tx;
                        }
                    }
                }
                resolve(chain.sort(hyperledger.compare));
            });
            return chainPromise;
        })
        .catch( function (err) {
            console.log(' *** ERROR: ' + err)
        });
}

hyperledger.compare = function (a,b) {
    if (a.nonHashData.localLedgerCommitTimestamp.seconds < b.nonHashData.localLedgerCommitTimestamp.seconds)
        return -1;
    if (a.nonHashData.localLedgerCommitTimestamp.seconds > b.nonHashData.localLedgerCommitTimestamp.seconds)
        return 1;
    return 0;
}

hyperledger.base64ToString = function (base64Blob) {
    return Buffer.from(base64Blob, 'base64').toString('ascii');
}


hyperledger.stringToBase64 = function (string) {
    return Buffer(string).toString('base64')
}


module.exports = hyperledger;

// Done.
