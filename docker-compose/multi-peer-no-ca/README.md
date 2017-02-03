# Multi-peer with CA

This needs at least 4 nodes to have Practical Byzantine Fault Tolerance (PBFT) working.

**NOTE: This is conflicting with our node-frontend at the moment.** I will fix this later this week.

### Starting the containers

Make sure your `$GOPATH` is correctly configured and you have pulled to repo according to standards.

```
docker-compose up
```

### Registering Chaincode

Documentation for HyperLedger CLI is [here](https://github.com/hyperledger/fabric/blob/master/docs/protocol-spec.md#6313-chaincode-deploy)

```bash
curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -H "Postman-Token: 74a9cd60-8c61-c0c4-ac0d-39ca203285dd" -d '{
    "jsonrpc": "2.0", 
    "method": "deploy",  
    "params": {
        "type":1, 
        "chaincodeID": {
            "name": "DecodedBlockChain",
            "path": "github.com/DecodedCo/blockchain-golang-chaincode"
        }, 
        "ctorMsg": { 
            "function":"init", 
            "args": [] 
        } 
    },
    "id": 0
}' "http://0.0.0.0:7050/chaincode"
```

The deploy will give you an ID that becomes the new `chaincodeID` name.

### Security Disabled

Security is disabled on all nodes for now.

If you now `GET` `http://0.0.0.0:7050/network/peers` you should get the following result:

```javascript
{
  "peers": [
    {
      "ID": {
        "name": "vp3"
      },
      "address": "172.17.0.4:7051",
      "type": 1
    },
    {
      "ID": {
        "name": "vp2"
      },
      "address": "172.17.0.5:7051",
      "type": 1
    },
    {
      "ID": {
        "name": "vp1"
      },
      "address": "172.17.0.6:7051",
      "type": 1
    },
    {
      "ID": {
        "name": "vp0"
      },
      "address": "172.17.0.3:7051",
      "type": 1
    }
  ]
}
```

### Security Enabled.

If you now `GET` `http://0.0.0.0:7050/network/peers` you should get the following result:

```javascript
{
  "peers": [
    {
      "ID": {
        "name": "vp1"
      },
      "address": "172.17.0.4:7051",
      "type": 1,
      "pkiID": "joEOyoOOMZwd8TcD7TmPodbUgRd1yox9JuadjHVAqSs="
    },
    {
      "ID": {
        "name": "vp2"
      },
      "address": "172.17.0.5:7051",
      "type": 1,
      "pkiID": "g51RMtKBITUO7Xaa9c+iI6S0dOp12smqz10zmCCOXzA="
    },
    {
      "ID": {
        "name": "vp3"
      },
      "address": "172.17.0.6:7051",
      "type": 1,
      "pkiID": "h+JxBz948VhwfbDJFjTERsV2/EmDAVHAEjy0LaVYDdM="
    },
    {
      "ID": {
        "name": "vp0"
      },
      "address": "172.17.0.3:7051",
      "type": 1,
      "pkiID": "o/xR3+4WEP6sgz1i/rlN/4K630xn70icNxYNZv/xDqQ="
    }
  ]
}
```