package main

package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "net/http"
    "time"

    "golang.org/x/net/proxy"
)

func main() {
    url := "https://example.com"
    socksAddress := "localhost:9998"

    socks, err := proxy.SOCKS5("tcp", socksAddress, nil, &net.Dialer{
        Timeout:   30 * time.Second,
        KeepAlive: 30 * time.Second,
    })
    if err != nil {
        panic(err)
    }

    client := &http.Client{
        Transport: &http.Transport{
            Dial:                socks.Dial,
            TLSHandshakeTimeout: 10 * time.Second,
        },
    }

    res, err := client.Get(url)
    if err != nil {
        panic(err)
    }
    content, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s", string(content))
}

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var requestURL = "http://ec2-184-73-148-15.compute-1.amazonaws.com/"
var ipAddress = "84.22.46.169"
var port = 8080

func main() {
	localAddr, err := net.ResolveIPAddr("ip", ipAddress)

	if err != nil {
		panic(err)
	}

	// You also need to do this to make it work and not give you a
	// "mismatched local address type ip"
	// This will make the ResolveIPAddr a TCPAddr without needing to
	// say what SRC port number to use.
	localTCPAddr := net.TCPAddr{
		IP:   localAddr.IP,
		Port: port,
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
		panic(err)
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
}
