package proxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var requestURL = "http://ec2-184-73-148-15.compute-1.amazonaws.com/"
var google = "https://google.com"

// VerifyProxy is anonymous
func VerifyProxy(p *Proxy) error {
	localAddr, err := net.ResolveIPAddr("ip", p.IPAddress)

	if err != nil {
		panic(err)
	}

	// You also need to do this to make it work and not give you a
	// "mismatched local address type ip"
	// This will make the ResolveIPAddr a TCPAddr without needing to
	// say what SRC port number to use.
	localTCPAddr := net.TCPAddr{
		IP:   localAddr.IP,
		Port: p.Port,
	}

	webclient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				LocalAddr: &localTCPAddr,
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	req, _ := http.NewRequest("GET", requestURL, nil)
	resp, err := webclient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
	}
	return nil
}
