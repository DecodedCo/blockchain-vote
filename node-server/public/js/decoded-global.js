var BLOCKCHAINNAME = 'DecodedBlockChain';
var DELAY = 1200;

function wait(ms){
	var start = new Date().getTime();
	var end = start;
	while(end < start + ms) {
		end = new Date().getTime();
	}
}

function guid() {
	function s4() {
		return Math.floor((1 + Math.random()) * 0x10000)
          .toString(16)
          .substring(1);
	}
	return s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();
}
