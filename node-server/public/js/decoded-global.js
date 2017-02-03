var BLOCKCHAINNAME = 'DecodedBlockChain';
var DELAY = 1200;

function wait(ms){
    var start = new Date().getTime();
    var end = start;
    while(end < start + ms) {
         end = new Date().getTime();
    }
}

function storeTransaction(transactionHash) {
    var storageKey = 'transactions'
    // If it does not exist, initialise
    if ( localStorage.getItem(storageKey) == null ) {
        var userTransactions = [ transactionHash ];
        localStorage.setItem(storageKey, JSON.stringify(userTransactions));
    }
    // Otherwise add the transaction to the existing array with transactions.
    else {
        var userTransactions = JSON.parse(localStorage.getItem(storageKey));
        userTransactions.push(transactionHash);
        localStorage.setItem(storageKey, JSON.stringify(userTransactions));
    }
}

// Create a unique browser id
if ( localStorage.getItem('dcdId') == null ) {
    console.log('Creating the Decoded ID...');
    localStorage.setItem('dcdId', Math.random().toString(36).substring(16));
}

function guid() {
    function s4() {
        return Math.floor((1 + Math.random()) * 0x10000)
          .toString(16)
          .substring(1);
    }
    return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
}
