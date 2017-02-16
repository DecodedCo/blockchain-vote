package main


import (
    "encoding/json"
    "errors"
    "fmt"
    "strconv"
    "sort"

    "github.com/hyperledger/fabric/core/chaincode/shim"

    utils "github.com/DecodedCo/blockchain-vote/golang-chaincode/utils"
)


// ============================================================================================================================


type Party struct {
    Id              string      `json:"id"`
    Name            string      `json:"name"`
    Voter           bool        `json:"voter"`
    Candidate       bool        `json:"candidate"`
    VotesToAssign   []string    `json:"votestoassign"`
    VotesReceived   []string    `json:"votesreceived"`
    CandidateUrl    string      `json:"candidateUrl"`
    ScreenshotUrl   string      `json:"screenshotUrl"`
}


type Candidates []Party // To assign the sorting functions


// ============================================================================================================================


func (slice Candidates) Len() int {
    return len(slice)
}

func (slice Candidates) Less(i, j int) bool {
    return len(slice[i].VotesReceived) > len(slice[j].VotesReceived)
}

func (slice Candidates) Swap(i, j int) {
    slice[i], slice[j] = slice[j], slice[i]
}


// ============================================================================================================================


func (dcc *DecodedChainCode) createParty(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    var emptyArgs []string
    if len(args) != 6 { // id, name, voter, candidate, votestoassign, votesreceived, candidateUrl, screenshotUrl
        err = errors.New("{\"Error\":\"Expecting 6 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    // The partyId needs to be unique. Check if the party does not already exist.
    partyId := args[0]
    // Get all the parties that are currently in the system.
    partyIds, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[0], emptyArgs)
    if err != nil {
        utils.PrintErrorFull("createParty - getDataArrayStrings", err)
        return nil, err
    }
    // Get all the candidates that are currently in the system.
    candidateIds, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[2], emptyArgs)
    if err != nil {
        utils.PrintErrorFull("createParty - getDataArrayStrings", err)
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
            CandidateUrl: args[4],
            ScreenshotUrl: args[5],
        }
        // Save new party
        if err = newParty.save(stub); err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        // Add party to the ledger.
        _, err = dcc.saveStringToDataArray(stub, PRIMARYKEY[0], partyId, partyIds)
        if err != nil {
            utils.PrintErrorFull("createParty - saveStringToDataArray", err)
            return nil, err
        }
        // If it is a candidate, add the the candidates-ledger
        if newParty.Candidate {
            _, err = dcc.saveStringToDataArray(stub, PRIMARYKEY[2], partyId, candidateIds)
            if err != nil {
                utils.PrintErrorFull("createParty - saveStringToDataArray", err)
                return nil, err
            }
        }
        // Done!
        utils.PrintSuccess("Added a new party: " + partyId)
        return nil, nil
    } else {
        err = errors.New(partyId + "` already exists.")
        utils.PrintErrorFull("createParty", err)
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
    if len(args) != 4 { // Id, VotesToAssign, VotesTransferred, VotesReceived
        err = errors.New("{\"Error\":\"Expecting 4 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
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
    // if party is a voter and there is a vote to transfer, delete from VotesToAssign slice
    voteTransferred := args[2]
    if party.Voter && voteTransferred != "" {
        // check if vote exists
        var emptyFn string
        args := []string{voteTransferred}
        _, err := dcc.readVote(stub, emptyFn, args)
        if err != nil {
            fmt.Printf("\t *** %s", err)
            return nil, err
        }
        for i, v := range party.VotesToAssign {
            if v == voteTransferred {
                party.VotesToAssign = append(party.VotesToAssign[:i], party.VotesToAssign[i+1:]...)
            }
        }
    }
    // if party is a candidate, add vote uuid to VotesReceived slice
    voteReceived := args[3]
    if party.Candidate && voteReceived != "" {
        party.VotesReceived = append(party.VotesReceived, voteReceived)
    }
    // Save the new party.
    if err = party.save(stub); err != nil {
        fmt.Printf("\t *** %s", err)
        return nil, err
    }
    fmt.Printf("\t --- Updated Party %s\n", partyId)
    return nil, nil
} // end of dcc.assignAssetToParty


func (dcc *DecodedChainCode) readAllCandidates(stub shim.ChaincodeStubInterface, fn string, args []string) ([]byte, error) {
    var err error
    var emptyArgs []string
    if len(args) != 0 {
        err = errors.New("{\"Error\":\"Expecting 0 arguments, got " + strconv.Itoa(len(args)) + ", \"Function\":\"" + fn + "\"}")
        utils.PrintErrorFull("", err)
        return nil, err
    }
    // Get candidates main ledger.
    candidateIds, err := dcc.getDataArrayStrings(stub, PRIMARYKEY[2], emptyArgs)
    if err != nil {
        utils.PrintErrorFull("readAllCandidates - getDataArrayStrings", err)
        return nil, err
    }
    // Iterate over all candidates to get the full details
    if len(candidateIds) > 0 {
        // Initialise an empty slice for the output
        var candidatesLedger []Party
        // Iterate over all parties and return the party object.
        for _, candidateId := range candidateIds {
            thisCandidate, err := dcc.getParty(stub, []string{ candidateId })
            if err != nil {
                utils.PrintErrorFull("readAllCandidates - getParty", err)
                return nil, err
            }
            candidatesLedger = append(candidatesLedger, thisCandidate)
        }
        // Sort the ledger by... number of votes received. (len(VotesReceived))
        sort.Sort(Candidates(candidatesLedger))
        // This gives us an slice with parties. Translate to bytes and return
        partiesLedgerBytes, err := json.Marshal(&candidatesLedger)
        if err != nil {
            utils.PrintErrorFull("readAllCandidates - Marshal", err)
            return nil, err
        }
        utils.PrintSuccess("Retrieved full information for all Parties.")
        return partiesLedgerBytes, nil 
    } else {
        return nil, nil
    }
} // readAllCandidates


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


// ============================================================================================================================
