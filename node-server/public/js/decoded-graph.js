var svgObj = d3.select('#svg-container'),
	width = +svgObj.attr('width'),
	height = +svgObj.attr('height');
margin = { left: 15, top: 15, bottom: 15, right: 15, spacing: 8 },
    fontHeight = 12,
    lbls = ['Genesis', 'Party created', 'Vote created', 'Voting', 'Activity'],
    rowHeight = (height - margin.top) / lbls.length,
    radius = 10;

function initBlockChainGraph(svg) {
    // Create the groups.
	labels = svg.append('g').attr('class', 'labels');
	edges = svg.append('g').attr('class', 'edges');
	nodes = svg.append('g').attr('class', 'nodes');
    // Create labels...
	var label = labels.selectAll('.label')
        .data(lbls)
        .enter().append('g')
        .attr('class', 'label')
        .attr('transform', function (d,i) {
	return 'translate(' + margin.left + ',' + (3+labelCentre(d)) + ')';
});
    // Add text.
	label.append('text').text(function(d) { return d; });
    // Create label-grid-lines.
	label.append('line')
        .attr('class', 'gridline')
        .attr('x1', 80)
        .attr('y1', -3)
        .attr('x2', width - margin.right)
        .attr('y2', -3)
        .attr('stroke', 'gray')
        .attr('stroke-width', '2')
        .attr('opacity', '0.1');
}

function labelCentre (thisLabelName) {
	var labelIx = lbls.indexOf(thisLabelName);
	return (margin.top + rowHeight * labelIx + 0.5 * (rowHeight-fontHeight));
}

function updateBlockChainGraph(data) {
	var t = d3.transition().duration(750);
    // JOIN data with old data.
	var node = nodes.selectAll('.node')
        .data(data, function(d) { return d; });
    // EXIT old elements not present in new data.
	node.exit().transition(t).style('fill-opacity', 1e-6).remove();
	node.enter().append('rect')
        .attr('class', function (d,i) { return d.css; })
        .attr('rx', 4) // rounded corners
        .attr('ry', 4) // rounded corners
        .attr('x', function (d,i) { return d.cx-11; })
        .attr('y', function (d,i) { return d.cy-11; })
        .attr('width', 22)
        .attr('height', 22)
        .attr('data-label', function (d,i) { return d.id; })
        .attr('data-num', function (d,i) { return window.btoa(d); })
        .on('click', function (d,i) {
	if ( d.own ) {
                // Toggle popup and populate it.
		$('#div-blockchain-inspector').find('.header').html('BLOCK ' + d.ix);
		$('#div-blockchain-inspector').find('.content').html(displayBlock(d));
		$('#qrcode1').qrcode(d.previousBlockHash);
		$('#qrcode2').qrcode(d.stateHash);
                // Show
		$('#div-blockchain-inspector').modal('show');
	}
});
    // Draw edges
	var edge = edges.selectAll('.edge')
        .data(data.slice(1, data.length), function(d) { return d; });
	edge.exit().remove();
	edge.enter().append('line')
        .attr('class','edge')
        .attr('id', function (d,i) { return d.id + '---' + i; })
        .attr('x1', function (d,i) { return d.line.x1; })
        .attr('x2', function (d,i) { return d.line.x2; })
        .attr('y1', function (d,i) { return d.line.y1; })
        .attr('y2', function (d,i) { return d.line.y2; })
        .attr('stroke', 'purple')
        .attr('stroke-width', '2')
        .attr('opacity', '0.6');
}

function displayBlock(blockData) {
    // Create table
	var tableBody = $('<table class="ui very basic table" id="table-block-inspector"><tbody></tbody></table>');

    // add transactions
	if ( blockData.transactions && blockData.transactions.length > 0 ) {
        // tableBody.append('<tr><td colspan=2>TRANSACTIONS (' + blockData.transactions.length + ')</td></tr>');
        // tableBody.append('<tr><td colspan=2 style="text-align:center;">TRANSACTIONS (' + blockData.transactions.length + ')</td></tr>');
		for ( var t=0; t<blockData.transactions.length; t++ ) {
			tableBody.append('<tr><td>Transaction id:</td><td>' + blockData.transactions[t].txid + '</td></tr>');
            // tableBody.append('<tr><td>Payload</td><td>' + blockData.transactions[t].payload + '</td></tr>');
			tableBody.append('<tr><td>Payload:</td><td>' + blockData.transactions[t].payloadDecoded + '</td></tr>');
		}
	}
    // Add block information
	tableBody.append('<tr><td>Previous Block Hash: </td><td>' + ( blockData.previousBlockHash || 'None' ) + '</td></tr>');
	tableBody.append('<tr><td>State Hash: </td><td>' + ( blockData.stateHash || 'None' ) + '</td></tr>');
    // in the onclick: $('#qrcode').qrcode(d.stateHash);
    // tableBody.append('<tr style="color:#ccc; font-weight: 200;"><td>State:</td><td style="overflow-wrap: break-word;">' + blockData.stateHash + '</td></tr>');
    // tableBody.append('<tr style="color:#ccc; font-weight: 200;"><td>Previous Block:</td><td style="overflow-wrap: break-word;">' + blockData.previousBlockHash + '</td></tr>');
	var date = new Date(blockData.nonHashData.localLedgerCommitTimestamp.seconds*1000); // add nano seconds?
	tableBody.append('<tr><td>Date/Time:</td><td>' + date + '</td></tr>');
	return tableBody;
}

function processData(dataArray) {
	var outputArray = [];
	var maxBlocks = 300;
	var startIndex = Math.max(0, dataArray.length-maxBlocks);
	var endIndex = Math.min(dataArray.length, maxBlocks+1);    
    // For each object. Add the x-y coordinates for the circle.
    // For each object add the node-class.
    // For each object add the line to the previous (if applicable)
	for ( var i=startIndex; i<endIndex; i++ ) {
		obj = dataArray[i];
        // Initialise flags.
		obj['id'] = '';
		obj['type'] = '';
		obj['to'] = '';
		obj['own'] = true;
        // Check if it is a transaction block.
		if ( 'transactions' in obj ) { 
			obj['id'] = obj.transactions[0].txid; // needs to be an array?
			if ( obj.transactions[0].txid == BLOCKCHAINNAME ) { 
				obj['type'] = 'deploy'; // Deploy block
			}
		}
		else if ( !('stateHash' in obj) ) {
			obj['type'] = 'deploy'; // Genesis block
		}
        //
		if ( !('ix' in obj) ) {
			obj['ix'] = i - startIndex;
		}
        // Get the x-coordinate.
		obj['cx'] = 100 + margin.left + radius + obj.ix * (radius * 2 + margin.spacing);
        // extend the canvas and grid lines to accommodate blocks
		var browserWidth = window.innerWidth,
			newWidth = obj.cx < browserWidth ? browserWidth : obj.cx + 30 ; // at least browserwidth 
		svgObj.attr('width', newWidth);
		d3.selectAll('.label').attr('x2', newWidth);
		d3.selectAll('.gridline').attr('x2', newWidth - 30);

        // Get the y-coordinate and class...
		obj['cy'] = labelCentre('Activity'); // Default.
		obj['css'] = 'node';
		if ( obj.type == 'deploy' ) {
			obj['cy'] = labelCentre('Genesis');
			obj['css'] += ' genesis';
		} 
		else if ( 'transactions' in obj ) {
			var payloadDecoded = obj.transactions[0].payloadDecoded;
			if ( payloadDecoded.indexOf('createParty') !== -1 ) {
				obj['cy'] = labelCentre('Party created');
			}
			else if ( payloadDecoded.indexOf('createVotesAndAssignToAll') !== -1 ) {
				obj['cy'] = labelCentre('Vote created');
			}
			else if ( payloadDecoded.indexOf('updateParty') !== -1 ) {
				obj['cy'] = labelCentre('Voting');
			}
		}
        // Add the lines...
		if ( i > 0 ) {
			obj['line'] = {
				x1: dataArray[i-1].cx,
				x2: obj.cx,
				y1: dataArray[i-1].cy,
				y2: obj.cy
			};
		}
		outputArray.push(obj);
	}
	return outputArray;
}

function parseAddOwner(ss) {
	var output = '';
	output += 'Function: <strong>RegisterCompany</strong> as:<br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Company Name: <strong>' + ss.split('\n')[4].trim() + '</strong><br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Description: <strong>' + ss.split('\n')[6].trim() + '</strong><br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Logo URL: <strong>' + ss.split('\n')[7].trim() + '</strong><br/><br/>';
    //output += 'The company has been assigned ID `<strong>' + ss.split('\n')[3] + '</strong>` with a balance of ' + ss.split('\n')[5].trim();
	return output;
}

function parseAddAsset(ss) {
	var output = '';
	output += 'Company (ID) <strong>' + ss.split('\n')[5].trim() + '</strong> triggered the<br/>';
	output += 'Function: <strong>RegisterAsset</strong> with:<br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Asset Name: <strong>' + ss.split('\n')[4].trim() + '</strong><br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Description: <strong>' + ss.split('\n')[8].trim() + '</strong><br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Logo URL: <strong>' + ss.split('\n')[9].trim() + '</strong><br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Asset Quantity: <strong>' + ss.split('\n')[6].trim() + '</strong><br/>';
	output += '&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Asset Price: <strong>' + ss.split('\n')[7].trim() + '</strong><br/><br/>';
	if ( ss.split('\n')[10].search('true') != -1 ) {
		output += 'The asset requires approval if more than ' + ss.split('\n')[11].trim() + ' items are sold.<br/><br/>';
	}
	if ( ss.split('\n')[12].trim() != '' ) {
		output += 'For each sale a text message "' + ss.split('\n')[13].trim() + '" is sent to ' + ss.split('\n')[12].trim();
	}
	return output;
}

function parseTransactions(string){
    // 
    // \nI\b\u0001\u0012\u0013\u0012\u0011DecodedBlockChain\u001a0\n\rtransactAsset\n\u0007mangoId\n\u0003dcd\n\u0002bc\n\u000267\n\u0003990\n\u0004TRUE
    //
    // \no\b\u0001\u0012\u0013\u0012\u0011DecodedBlockChain\u001aV\n\u0012declineTransaction\n@2972c91ff7f7d3465b0da2cd4b91c3c9848ae29b849a2946e32df1b5ad808bef
}


function wrapText(string) {
	return string.match(/.{1,56}/g).join('<br/>');
}

// EOF
