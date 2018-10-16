package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "time"
        "crypto/tls"
        "strings"
)
type Items struct {
        ID string `json:"id"`
        Name string `json:"name"`
        Version string `json:"version"`	
}

type baseImage struct {
        Items []Items `json:"items"` 
	Token json.RawMessage `json:"continuationToken"`
}

func getUrl(url string, uri string, tokenString string) *http.Request {

	var fullUrl string
        if tokenString != "" {
          tok := "continuationToken=%s"
          fullUrl = fmt.Sprintf(url + tok, tokenString +"&"+ uri)
	} else {
	    fullUrl = fmt.Sprintf(url + uri)
          }

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
        return body

}

func main() {

	// Set variables
        var contTokens string
        contToken := baseImage{}
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
	for string(contToken.Token) != "null" { 

          tokenString1 := string(contToken.Token)
          tokenString := strings.Trim(tokenString1, `"`)

          req, getErr := baseClient.Do(getUrl(url, uri, tokenString))
          if getErr != nil {
                 log.Fatal(getErr)
          }

          body := checkBody(req)
	  //fmt.Println(string(body))
          jsonErr := json.Unmarshal(body, &contToken)
          if jsonErr != nil {
                  log.Fatal(jsonErr)
          }


	  if tokenString != "null" {
            contTokens += tokenString+"\n"
	    //fmt.Println(contToken.Items, string(contToken.Token))
	    fmt.Println(contToken.Items)
	  }
	} 

	// Print output
        fmt.Println(contTokens)
}

