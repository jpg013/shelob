package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var requestURL = "http://ec2-184-73-148-15.compute-1.amazonaws.com/"

// VerifyProxy is anonymous
func VerifyProxy(p *Proxy) {
	proxyURL, err := url.Parse(p.Socket)
	fmt.Println(proxyURL)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

	req, _ := http.NewRequest("GET", requestURL, nil)
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error verifying proxy: ", err.Error())
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
}
