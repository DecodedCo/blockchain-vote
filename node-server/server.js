const express = require('express');    
const bodyParser = require('body-parser');
const requestpromise = require('request-promise');
const log = require('./app/utils/log');

// express app =======================================================
const port = process.env.PORT || 3000;
const app = express();
const ip = '0.0.0.0'; //'127.0.0.1'
app.listen(port, ip, () => {
	log.info('App listening at http://%s:%s', ip, port);
});

// hyperledger =======================================================
const hyperledger = require('./app/hyperledger.js');
// Check if the blockchain application is up, running and deployed. If not, deploy.
// Return the current blockchain state. This is to avoid making annoying requests...
hyperledger.init(requestpromise);
hyperledger.getFullBlockChain(requestpromise)
    .then( function (data) {
	hyperledger.BLOCKCHAIN = data;
	log.info('Success: blockchain retrieved and stored.');
	log.info('The blockchain currently has ' + hyperledger.BLOCKCHAIN.length + ' blocks.');
})
    .catch( function (err) {
	log.error('Failed to grab the current blockchain:', err);
});


// app config ========================================================
app.use(express.static('public')); // Directory to serve files from
app.set('view engine', 'ejs'); // Set up ejs for templating
app.use(require('morgan')('dev')); // HTTP req/res logging
app.use(bodyParser.urlencoded({ extended: true })); // get information from html forms
app.use(bodyParser.json()); // get information from html forms
app.use(function (req, res, next) { // update internal state with each request
	hyperledger.getFullBlockChain(requestpromise)
        .then( function (data) {
	hyperledger.BLOCKCHAIN = data;
	log.info('SUCCESS: Updated internal blockchain state.');
	log.info('The BLOCKCHAIN currently has ' + hyperledger.BLOCKCHAIN.length + ' blocks.');
	next();
})
        .catch( function (err) {
	res.status(500).json(err);
});
});

// routes ============================================================
require('./app/routes/routes.js')(app, requestpromise, hyperledger);

if (process.env.SCREENSHOTLAYER_KEY === undefined) {
	log.warn('SCREENSHOTLAYER_KEY is undefined');
}

// EOF
