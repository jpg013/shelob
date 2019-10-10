package proxy

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	dom "github.com/jpg013/go_dom"
)

var (
	base64Prefixes = []string{
		"data:image/jpeg;charset=utf-8;base64,",
		"data:image/png;charset=utf-8;base64",
	}
	// ProxyURLList represents list of urls to fetch proxy ip data from.
	ProxyURLList = []string{
		"https://www.proxyrotator.com/free-proxy-list/1",
		"https://www.proxyrotator.com/free-proxy-list/2",
		"https://www.proxyrotator.com/free-proxy-list/3",
		"https://www.proxyrotator.com/free-proxy-list/4",
		"https://www.proxyrotator.com/free-proxy-list/5",
		"https://www.proxyrotator.com/free-proxy-list/6",
		"https://www.proxyrotator.com/free-proxy-list/7",
		"https://www.proxyrotator.com/free-proxy-list/8",
		"https://www.proxyrotator.com/free-proxy-list/9",
		"https://www.proxyrotator.com/free-proxy-list/10",
	}
)

type ProxyRefreshTask struct {
	requestThrottle <-chan time.Time
	responseCh      chan *http.Response
	inProgress      bool
	mux             sync.Mutex
}

func NewProxyRefreshTask() *ProxyRefreshTask {
	return &ProxyRefreshTask{
		inProgress: false,
	}
}

func ParseProxyResponse(resp *http.Response) {
	defer resp.Body.Close()

	doc, err := dom.ParseHTMLDocument(resp.Body)

	if err != nil {
		log.Fatal(fmt.Sprintf("error parsing html document: %v", err))
	}

	ScrapeProxyRotatorList(doc)
}

func (p *ProxyRefreshTask) fetchProxyList(url string) {
	log.Println("Fetch proxy list : ", url)

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(fmt.Sprintf("error fetching proxy list: %v", err))
	}

	// Send to response channel
	p.responseCh <- resp
}

func (p *ProxyRefreshTask) run() {
	for _, url := range ProxyURLList {
		<-p.requestThrottle
		go p.fetchProxyList(url)
	}
}

func (p *ProxyRefreshTask) handleProxyResponse() {
	go func() {
		for {
			select {
			case r := <-p.responseCh:
				go ParseProxyResponse(r)
			}
		}
	}()
}

func (p *ProxyRefreshTask) scheduleRefreshTask(delay time.Duration, stop chan bool) {
	// Create stop channel
	stop = make(chan bool)

	go func() {
		go p.run()
		for {
			select {
			case <-time.After(delay):
				go p.run()
			case <-stop:
				return
			}
		}
	}()
}

func (p *ProxyRefreshTask) Start(delay time.Duration) (stop chan bool, err error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.inProgress == true {
		return stop, errors.New("refresh task already scheduled")
	}

	// Set in progress to true
	p.inProgress = true

	// throttle rate in seconds
	rate := time.Second

	// create the rate limiter
	p.requestThrottle = time.Tick(rate)

	// create response channel
	p.responseCh = make(chan *http.Response, 100)

	// Create stop channel
	stop = make(chan bool)

	p.handleProxyResponse()
	p.scheduleRefreshTask(delay, stop)

	return stop, err
}
