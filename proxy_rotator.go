package main

// import (
// 	"errors"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	dom "github.com/jpg013/go_dom"
// )

// // ProxyURLList represents list of urls to fetch proxy ip data from.
// var ProxyURLList = []string{
// 	"https://www.proxyrotator.com/free-proxy-list/1",
// 	// "https://www.proxyrotator.com/free-proxy-list/2",
// 	// "https://www.proxyrotator.com/free-proxy-list/3",
// 	// "https://www.proxyrotator.com/free-proxy-list/4",
// 	// "https://www.proxyrotator.com/free-proxy-list/5",
// 	// "https://www.proxyrotator.com/free-proxy-list/6",
// 	// "https://www.proxyrotator.com/free-proxy-list/7",
// 	// "https://www.proxyrotator.com/free-proxy-list/8",
// 	// "https://www.proxyrotator.com/free-proxy-list/9",
// 	// "https://www.proxyrotator.com/free-proxy-list/10",
// }

// type ProxyRotator struct {
// 	isRefreshTaskScheduled bool
// }

// func (p *ProxyRotator) RefreshProxyAddresses() {
// 	log.Println("refresh proxy addresses")

// 	for _, url := range ProxyURLList {
// 		resp, err := http.Get(url)

// 		if err != nil {
// 			log.Fatal(fmt.Sprintf("error fetching proxy list: %v", err))
// 		}

// 		defer resp.Body.Close()

// 		// bodyBytes, _ := ioutil.ReadAll(resp.Body)
// 		// bodyString := string(bodyBytes)
// 		// fmt.Println(bodyString)

// 		// resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

// 		// Parse response into html document
// 		doc, err := dom.ParseHTMLDocument(resp.Body)

// 		if err != nil {
// 			log.Fatal(fmt.Sprintf("error parsing html document: %v", err))
// 		}

// 		ParseProxyList(doc)

// 	}
// }

// func (p *ProxyRotator) ScheduleRefreshProxyIPTask(delay time.Duration) (<-chan bool, error) {
// 	if p.isRefreshTaskScheduled == true {
// 		return nil, errors.New("refresh task is already scheduled")
// 	}

// 	stop := make(<-chan bool)
// 	p.isRefreshTaskScheduled = true

// 	go func() {
// 		p.RefreshProxyAddresses()
// 		for {
// 			select {
// 			case <-time.After(delay):
// 				p.RefreshProxyAddresses()
// 			case <-stop:
// 				return
// 			}
// 		}
// 	}()

// 	return stop, nil
// }

// func NewProxyRotator() *ProxyRotator {
// 	return &ProxyRotator{}
// }
