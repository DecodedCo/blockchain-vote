package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/nu7hatch/gouuid"

    // utils "./utils"
)


// ============================================================================================================================


type Vote struct {
    Uuid         string      `json:"uuid"`
}


// ============================================================================================================================


func (dcc *DecodedChainCode) createVotesAndAssignToAll(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    if len(args) != 0 { 
        err = errors.New("{\"Error\":\"Expecting 0 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    var emptyFn string
    var emptyArgs []string
    // Get all parties
    partiesLedgerBytes, err := dcc.readAllParties(stub, emptyFn, emptyArgs)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    var partiesLedger []Party
    if err := json.Unmarshal(partiesLedgerBytes, &partiesLedger); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Iterate over parties
    for _, party := range partiesLedger {
        // Filter for party.Voter == true
        if party.Voter {
            // create a new vote with uuid
            u4, err := uuid.NewV4()
            if err != nil {
                fmt.Printf("\t *** %s", err)
                return nil, err
            }
            vote := Vote{ 
                Uuid: u4.String(),
            }
            // save new vote in blockchain
            if err = vote.save(stub); err != nil {
                fmt.Printf("\t *** %s", err)
                return nil, err
            }
            // assign new vote to voting party
            var emptyFn string
            args := []string{ party.Id, vote.Uuid, "" }
            if _, err := dcc.updateParty(stub, emptyFn, args); err != nil {
                fmt.Printf("\t *** %s", err)
                return nil, err
            }
        }
    }
    return nil, nil
} // end of dcc.createvotesAndAssignToAll


func (v *Vote) save(stub shim.ChaincodeStubInterface) (error) {
    var err error
    voteBytesToWrite, err := json.Marshal(&v)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return err
    }
    err = stub.PutState(v.Uuid, voteBytesToWrite)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return err
    }
    fmt.Printf("\t --- Saved vote %+v to blockchain\n", &v)
    return nil
} // end of p.save


// ============================================================================================================================
