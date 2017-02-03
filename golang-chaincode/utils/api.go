
package utils


import (
    "encoding/json"
    "net/http"
    "io/ioutil"
    "fmt"

    "crypto/tls"
    "crypto/x509"
)


// ============================================================================================================================


// types and structs


// ============================================================================================================================


func GetExchangeRate() {
    resp, err := http.Get("http://api.fixer.io/latest?base=USD&symbols=GBP")
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }
    var dd interface{}
    err = json.Unmarshal(body, &dd)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("Results: %v\n", dd)
} // end of GetExchangeRate


func GetBTCUSD() {
    resp, err := http.Get("http://api.coindesk.com/v1/bpi/currentprice.json")
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }
    var dd interface{}
    err = json.Unmarshal(body, &dd)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Printf("Results: %v\n", dd)
    // fmt.Printf("Results: %v\n", dd.bpi.USD.xrate)
} // end of GetBTCUSD


func SendSMS(phoneNumber string, message string) (error) {

    // Load client cert
    cert, err := tls.LoadX509KeyPair(".keys/cert.pem", ".keys/private.pem")
    if err != nil {
        return err
    }

    // Load CA cert
    caCert, err := ioutil.ReadFile(".keys/cert.pem")
    if err != nil {
        return err
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Setup HTTPS client
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs: caCertPool,
        InsecureSkipVerify: true, 
    }
    tlsConfig.BuildNameToCertificate()
    transport := &http.Transport{TLSClientConfig: tlsConfig}
    client := &http.Client{Transport: transport}

    resp, err := client.Get("https://internal.decoded.com/tech/sms/?to=" + phoneNumber + "&message=" + message)
    if err != nil {
        return err
    }

    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    return nil

} // end of SendSMS


// ============================================================================================================================

