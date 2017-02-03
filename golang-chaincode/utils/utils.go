
package utils


import (
    "crypto/sha256"
    "encoding/hex"
)


// ============================================================================================================================


// types and structs


// ============================================================================================================================


func IsElementInSlice(slice []string, element string) (bool) {
    // Initialise return as false.
    check := false
    // Iterate over the list to see if the value is in there.
    for _, val := range slice {
        if val == element {
            check = true
            return check
        }
    }
    // Failsafe return.
    return check
}


// Function to return the index of an element in a slice.
func FindElementIndex(slice []string, element string) (int) {
    // Initialise return as false.
    ix := -1
    // Iterate over the list to see if the value is in there.
    for i, val := range slice {
        if val == element {
            return i
        }
    }
    // Failsafe return.
    return ix
}


// Function to return the SHA256 for
func HashSHA256(s string) (string) {
    h := sha256.New()
    h.Write([]byte(s))
    sha256Hash := hex.EncodeToString(h.Sum(nil))
    return sha256Hash
}   


// Function to delete an element from a slice.?
func DeleteElementFromSlice(slice []string, element string) ([]string) {
    var emptySlice []string
    if len(slice) == 0 {
        return emptySlice
    }
    if !IsElementInSlice(slice, element) { // It is not in the slice.
        return slice
    }
    if len(slice) == 1 {
        return emptySlice
    }
    ix := FindElementIndex(slice, element)
    slice = append(slice[:ix], slice[ix+1:]...)
    return slice
}


// ============================================================================================================================

