package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "time"
        "crypto/tls"
//        "strings"
)

type contToken struct {
        Token string `json:"continuationToken,omitempty"`
}

func getUrl(url string, uri string, contToken1 contToken) *http.Request {

	var fullUrl string

        if contToken1.Token != "" {
	  //USED FOR DEBUG!
          fmt.Printf("\nUsing TOKEN in fullUrl\n")
          tok := "continuationToken=%s"
          fullUrl = fmt.Sprintf(url + tok, contToken1.Token +"&"+ uri)
	} else {
	    //USED FOR DEBUG!
            fmt.Printf("\nNOT USING TOKEN in fullUrl\n")
	    fullUrl = fmt.Sprintf(url + uri)
          }
        //USED FOR DEBUG!!
        fmt.Printf("\n%s\n", fullUrl) 
        req, err := http.NewRequest(http.MethodGet, fullUrl, nil)
        if err != nil {
                log.Fatal(err)
        }
        return req
}

func checkBody(res *http.Response) []byte {

        body, readErr := ioutil.ReadAll(res.Body)
        if readErr != nil {
                log.Fatal(readErr)
        }
        //USED FOR DEBUG!!
	//fmt.Println(string(body))
        return body

}

func main() {

	// Set variables
        contTokens := ""
        contToken1 := contToken{}
        url := "https://registry.lappy.maddocks:8443/service/rest/v1/components?"
        uri := "repository=public"
        // Define the http transport options
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

	// Define the http client options
        baseClient := http.Client{
                Timeout: time.Second * 2, // Maximum of 2 secs
                Transport: tr,
        }

	// Check if continuation Token is received (process accordingly!)
        fmt.Println(contToken1.Token)
	//for strings.Contains(contToken1.Token, "null") {
	for contToken1.Token != "null" { 
          //USED FOR DEBUG!!
	  fmt.Println("within for loop, contToken1 = ", contToken1.Token)
          req, getErr := baseClient.Do(getUrl(url, uri, contToken1))
          if getErr != nil {
                 log.Fatal(getErr)
          }
          body := checkBody(req)
	  contToken1.Token = "null"
          jsonErr := json.Unmarshal(body, &contToken1)
          if jsonErr != nil {
                  log.Fatal(jsonErr)
          }
          contTokens += contToken1.Token
	} 
	// Print output
        fmt.Println(contTokens)
}

