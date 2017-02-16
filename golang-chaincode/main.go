/*

DECODED HYPERLEDGER APPLICATION

DecodedChainCode functions:
- main, Init, Invoke, Query (standard and required)
- read: reads the contents of a specific key.
- readAll: reads all the primary keys.
- getDataArrayStrings - private function
- saveStringToDataArray
- saveLedger

*/

package main


import (
    "encoding/json"
    "errors"
    "fmt"

    "github.com/hyperledger/fabric/core/chaincode/shim"

    utils "github.com/DecodedCo/blockchain-vote/golang-chaincode/utils"
)


// ============================================================================================================================


// Create the struct to tie the methods to.
type DecodedChainCode struct {
}

var PRIMARYKEY = [3]string{ "Parties", "Votes", "Candidates" }


// ============================================================================================================================
// Main


func main() {
    err := shim.Start(new(DecodedChainCode))
    if err != nil {
        fmt.Printf("Error starting Simple chaincode: %s", err)
    }
}


// ============================================================================================================================


// Init resets all the things
func (dcc *DecodedChainCode) Init(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    if len(args) != 0 {
        err = errors.New("{\"Error\":\"Incorrect number of arguments\", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Initialise the empty datastores. The owners and assets datastores are going to store an array of all keys/ids.
    var blank []string
    blankBytes, err := json.Marshal(&blank)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    //
    if err = stub.PutState(PRIMARYKEY[0], blankBytes); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    if err = stub.PutState(PRIMARYKEY[1], blankBytes); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    if err = stub.PutState(PRIMARYKEY[2], blankBytes); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Done.
    utils.PrintSuccess("Initialisation complete")
    return nil, nil
} 


// Invoke is our entry point to invoke a chaincode function
func (dcc *DecodedChainCode) Invoke(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    // Handle different functions
    if fn == "init" { //initialize the chaincode state, used as reset
        return dcc.Init(stub, fn, args)
    } else if fn == "createParty" {
        return dcc.createParty(stub, fn, args)
    } else if fn == "createVotesAndAssignToAll" {
        return dcc.createVotesAndAssignToAll(stub, fn, args)
    } else if fn == "updateParty" {
        return dcc.updateParty(stub, fn, args)
    }
    // In any other case.
    fmt.Println("\t*** ERROR: Invoke function did not find ChainCode function: " + fn) // Error handling.
    return nil, errors.New(" --- INVOKE ERROR: Received unknown function invocation")
}


// Query is our entry point for queries
func (dcc *DecodedChainCode) Query(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    // Handle different functions
    if fn == "read" { // read a variable
        return dcc.read(stub, fn, args) 
    } else if fn == "readParty" {
        return dcc.readParty(stub, fn, args)
    } else if fn == "readAllParties" {
        return dcc.readAllParties(stub, fn, args)
    } else if fn == "readAllCandidates" {
        return dcc.readAllCandidates(stub, fn, args)
    }
    fmt.Println("\t*** ERROR: Query function did not find ChainCode function: " + fn)
    return nil, errors.New(" --- QUERY ERROR: Received unknown function query")
}


// ============================================================================================================================


// Function that reads the bytes associated with a data-key and returns the byte-array.
func (dcc *DecodedChainCode) read(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    if len(args) != 1 { // needs a data key to read.
        err = errors.New("{\"Error\":\"Incorrect number of arguments\", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    dataKey := args[0]
    dataBytes, err := stub.GetState(dataKey)
    if dataBytes == nil { // deals with non existing data keys.
        err = errors.New("{\"Error\":\"State " + dataKey + " does not exist\", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    if err != nil {
        err = errors.New("{\"Error\":\"Failed to get state for " + dataKey + "\", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    } 
    fmt.Println("\t--- Read the ledger: " + dataKey)
    return dataBytes, nil
}


func (dcc *DecodedChainCode) readAll(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    var emptyArgs []string
    if len(args) != 0 {
        err = errors.New("{\"Error\":\"Incorrect number of arguments\", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // get all parties - returns an array of strings.
    partiesLedger, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[0], emptyArgs)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // get all votes
    votesLedger, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[1], emptyArgs)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Create a map of both
    m := map[string][]string{ 
        PRIMARYKEY[0]: partiesLedger, 
        PRIMARYKEY[1]: votesLedger,
    }
    // Cast to JSON
    mStr, err := json.Marshal(m)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Return as bytes.
    fmt.Println("\t--- Read the main ledgers ")
    out := []byte(string(mStr))
    return out, nil
}


func (dcc *DecodedChainCode) getDataArrayStrings(stub shim.ChaincodeStubInterface, dataKey string, args []string) ([]string, error) {
    var err error
    var empty []string
    if len(args) != 0 {
        err = errors.New("{\"Error\":\"Incorrect number of arguments\", \"Function\":\"getDataArrayStrings\"}")
        fmt.Printf("\t *** %s", err)
        return empty, err
    }
    arrayBytes, err := stub.GetState(dataKey)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return empty, err
    }
    var outputArray []string
    err = json.Unmarshal(arrayBytes, &outputArray)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return empty, err
    }
    return outputArray, nil
}


func (dcc *DecodedChainCode) saveStringToDataArray(stub shim.ChaincodeStubInterface, dataKey string, addString string, ledger []string) ([]byte, error) {
    var err error
    // Add the string to the array
    ledger = append(ledger, addString)
    // err = dcc.saveLedger(stub, dataKey, ledger)
    if err = dcc.saveLedger(stub, dataKey, ledger); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    return nil, nil
}


func (dcc *DecodedChainCode) saveLedger(stub shim.ChaincodeStubInterface, dataKey string, ledger []string) (error) {
    var err error
    // Marshall the ledger to bytes
    bytesToWrite, err := json.Marshal(&ledger)
    if err != nil {
        return err
    }
    // Save the array.
    err = stub.PutState(dataKey, bytesToWrite)
    if err != nil {
        return err
    }
    return nil
}


// ============================================================================================================================
