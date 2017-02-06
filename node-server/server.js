const express = require('express');    
const bodyParser = require('body-parser');
const requestpromise = require('request-promise');
const randomcolor= require('randomcolor');

// express app =======================================================
const port = process.env.PORT || 3000;
const app = express();
const ip = "0.0.0.0" //"127.0.0.1"
const server = app.listen(port, ip, () => {
    console.log(" *** Example app listening at http://%s:%s", ip, port);
});


// web sockets =======================================================
const socket = require('socket.io')(server);


// hyperledger =======================================================
const hyperledger = require('./app/hyperledger.js');
// Check if the blockchain application is up, running and deployed. If not, deploy.
// Return the current blockchain state. This is to avoid making annoying requests...
hyperledger.init(requestpromise);
hyperledger.getFullBlockChain(requestpromise)
    .then( function (data) {
        hyperledger.BLOCKCHAIN = data;
        console.log('\t*** SUCCESS: BLOCKCHAIN retrieved and stored.');
        console.log('\t*** The BLOCKCHAIN currently has ' + hyperledger.BLOCKCHAIN.length + ' blocks.')
    })
    .catch( function (err) {
        console.log(' *** ERROR: Failed to grab the current BLOCKCHAIN');
    });


// app config ========================================================
app.use(express.static('public')); // Directory to serve files from
app.set('view engine', 'ejs'); // Set up ejs for templating
app.use(require('morgan')('dev')); // Logging
app.use(bodyParser.urlencoded({ extended: true })); // get information from html forms
app.use(bodyParser.json()); // get information from html forms


// routes ============================================================
require('./app/routes/routes.js')(app, requestpromise, hyperledger, socket, randomcolor);


// EOF
