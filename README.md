# blockchain-vote

### Install dependencies

1. Download and install Docker: https://www.docker.com/products/docker#/mac (or `brew install docker`)
2. Download and install Go: https://golang.org/doc/install (or `brew install golang`)
3. Download and install Node.js: https://nodejs.org/en/download/

### Set up your GO workspace

1. Create a directory to act as your Go workspace. Suggestion: `~/dev/go`.
2. In your Go workspace `mkdir -p src/github.com/`
3. In `.../go/src/github.com` add the hyperledger Go code:

```
mkdir hyperledger
cd hyperledger
git clone git@github.com:hyperledger/fabric.git
cd fabric
git fetch
git checkout v0.6
```

4. In your `~/.zshrc` file add `export GOPATH=$HOME/dev/go` or whatever that path is.

### Run the application 

In your dev directory. Assumed for this guide: `~/dev`

```
git clone git@github.com:DecodedCo/blockchain-vote.git
cd blockchain-vote
```

#### Docker

Make sure you have Docker running on your mac, then:

```
cd ~/dev/blockchain-vote/docker-compose
docker pull hyperledger/fabric-peer:latest
docker pull hyperledger/fabric-membersrvc:latest
```

Wait a few moments then

```
cd one-peer-one-ca
docker-compose up
```

#### Go

Open up a new shell, then:

```
cd ~/dev/blockchain-vote/golang-chaincode
go build .
```

Wait a moment then:

```
CORE_CHAINCODE_ID_NAME=DecodedBlockChain CORE_PEER_ADDRESS=0.0.0.0:7051 ./golang-chaincode
```

#### Node

Open up another new shell, then:

```
cd ~/dev/blockchain-vote/node-server
npm install
npm start
```

Open a web browser and navigate to `localhost:3000`

### Shutting down

#### Docker

Open up a fourth shell, then:

```
cd ~/dev/blockchain-vote/docker-compose/one-peer-one-ca
docker-compose stop
```

That stops the process and preserves the state of the blockchain. To delete the blockchain, run:

```
docker rm -f $(docker ps -aq)
```

#### Go and Node

The Go chaincode should stop running when the docker container stops. If not, the Go and Node processes can be stopped using `ctrl+C`
