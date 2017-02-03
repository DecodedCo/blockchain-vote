package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"

    utils "./utils"
)

type Party struct {
    Id              string      `json:"id"`
    Name            string      `json:"name"`
    Voter           bool        `json:"voter"`
    Candidate       bool        `json:"candidate"`
    VotesToAssign   []string    `json:"votestoassign"`
    VotesReceived   []string    `json:"votesreceived"`
}

func (dcc *DecodedChainCode) createParty(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    var emptyArgs []string
    if len(args) != 4 { // partyId
        err = errors.New("{\"Error\":\"Expecting 4 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // The partyId needs to be unique. Check if the party does not already exist.
    partyId := args[0]
    // Get all the parties that are currently in the system.
    partyIds, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[0], emptyArgs)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Check if the partyId exists in the current ledger of parties.
    partyExists := utils.IsElementInSlice(partyIds, partyId)
    if partyExists == false {
        voter, err := strconv.ParseBool(args[2])
        if err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        candidate, err := strconv.ParseBool(args[3])
        if err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        // Create a new party
        var newParty = Party{ 
            Id: partyId,
            Name: args[1],
            Voter: voter,
            Candidate: candidate,
        }
        // Save new party
        if err = newParty.save(stub); err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        // Add party to the list.
        _, err = dcc.saveStringToDataArray(stub, PRIMARYKEY[0], partyId, partyIds)
        if err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        // Done!
        fmt.Println("\t--- Added a new party: " + partyId)
        return nil, nil
    } else {
        err = errors.New(partyId + "` already exists.")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Redundancy.
    return nil, nil
} // end of dcc.createParty

func (dcc *DecodedChainCode) readParty(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    if len(args) != 1 { // id
        err = errors.New("{\"Error\":\"Expecting 1 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    id := args[0]
    var returnSlice []Party
    party, err := dcc.getParty(stub, []string{ id })
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    returnSlice = append(returnSlice, party)
    // This gives us an slice with parties. Translate to bytes and return
    returnSliceBytes, err := json.Marshal(&returnSlice)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    fmt.Printf("\t--- Retrieved full information for Party %s", id)
    return returnSliceBytes, nil   
}

func (dcc *DecodedChainCode) readAllParties(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    var emptyArgs []string
    if len(args) != 0 {
        err = errors.New("{\"Error\":\"Expecting 0 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Get all parties - returns an slice of strings - partyIds
    partyIds, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[0], emptyArgs)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    if len(partyIds) > 0 {
        // Initialise an empty slice for the output
        var partiesLedger []Party
        // Iterate over all parties and return the party object.
        for _, partyId := range partyIds {
            thisParty, err := dcc.getParty(stub, []string{ partyId })
            if err != nil {
                fmt.Printf("\t *** %s", err)
                return nil, err
            }
            partiesLedger = append(partiesLedger, thisParty)
        }
        // This gives us an slice with parties. Translate to bytes and return
        partiesLedgerBytes, err := json.Marshal(&partiesLedger)
        if err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        fmt.Println("\t--- Retrieved full information for all Parties.")
        return partiesLedgerBytes, nil   
    } else {
        return nil, nil
    }
    return nil, nil // redundancy
} // end of dcc.readAllParties

func (dcc *DecodedChainCode) getParty(stub shim.ChaincodeStubInterface, args []string) (Party, error) {
    var party Party // We need to have an empty party ready to return in case of an error.
    var err error
    if len(args) != 1 { // Only needs a party id.
        err = errors.New("{\"Error\":\"Incorrect number of arguments\", \"Function\":\"getParty\"}")
        fmt.Printf("\t *** %s", err)
        return party, err
    }
    partyId := args[0]
    partyBytes, err := stub.GetState(partyId)
    if partyBytes == nil {
        err = errors.New("{\"Error\":\"State " + partyId + " does not exist\", \"Function\":\"getParty\"}")
        fmt.Printf("\t *** %s", err)
        return party, err
    }
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return party, err
    }
    err = json.Unmarshal(partyBytes, &party)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return party, err
    }
    return party, nil
} // end of dcc.getParty

func (dcc *DecodedChainCode) updateParty(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    if len(args) != 3 { // Id, VotesToAssign, VotesReceived
        err = errors.New("{\"Error\":\"Expecting 3 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // Load the current data
    partyId := args[0]
    party, err := dcc.getParty(stub, []string{ partyId })
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // if party is a voter, add vote uuid to VotesToAssign slice
    voteToAssign := args[1]
    if party.Voter && voteToAssign != "" {
        party.VotesToAssign = append(party.VotesToAssign, voteToAssign)
    }
    // if party is a candidate, add vote uuid to VotesReceived slice
    voteReceived := args[2]
    if party.Voter && voteReceived != "" {
        party.VotesReceived = append(party.VotesReceived, voteReceived)
    }
    // Save the new party.
    if err = party.save(stub); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    return nil, nil
} // end of dcc.assignAssetToParty

func (p *Party) save(stub shim.ChaincodeStubInterface) (error) {
    var err error
    partyBytesToWrite, err := json.Marshal(&p)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return err
    }
    err = stub.PutState(p.Id, partyBytesToWrite)
    if err != nil {
        fmt.Printf("\t *** %s", err)
        return err
    }
    fmt.Printf("\t --- Saved party %v to blockchain\n", &p.Id)
    return nil
} // end of p.save
