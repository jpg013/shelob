package proxy

import (
	"errors"
	"log"
	"shelob/db"
	"sync"
)

// Site represents a basic proxy site with a Run method
// that can be called with a cancel / out channel. The implementations of
// Run will vary drastically among different proxy sites.
type Site interface {
	Run(out chan<- *Proxy) error
	Stop() error
}

// RefreshTask type
type RefreshTask struct {
	inProgress bool
	// proxyChan receives new proxies from sites
	proxyChan chan *Proxy
	stop      chan bool
	// Simple mutex for synchronizing the start phase of the refresh task
	mux sync.Mutex
	// list of proxy sites as data sources
	sites []Site
}

// NewRefreshTask factory returns a new refresh task
func NewRefreshTask() *RefreshTask {
	sites := []Site{NewProxyRotatorSite()}
	refreshTask := &RefreshTask{
		inProgress: false,
		sites:      sites,
	}

	return refreshTask
}

// Starting the refresh task by calling Run() on each ProxySite.
func (t *RefreshTask) startRefreshTask() {
	for _, site := range t.sites {
		err := site.Run(t.proxyChan)

		if err != nil {
			log.Fatalf("error starting refresh task: %v", err.Error())
		}
	}
}

// stop the refresh task by calling Stop() on each ProxySite.
func (t *RefreshTask) stopRefreshTask() {
	for _, site := range t.sites {
		err := site.Stop()

		if err != nil {
			log.Fatalf("error stopping refresh task: %v", err.Error())
		}
	}
}

// proxy handler waits for new proxies from the proxyChan
// and inserts them in to the database.
func (t *RefreshTask) proxyHandler() {
	handlerFunc := func(p *Proxy) {
		db.InsertOne("proxy", p)
	}

	go func() {
		for {
			select {
			case proxy := <-t.proxyChan:
				go handlerFunc(proxy)
			case <-t.stop:
				return
			}
		}
	}()
}

// Start will start the refresh task, returning an error if unable to do so.
func (t *RefreshTask) Start() (stop chan bool, err error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	if t.inProgress == true {
		return stop, errors.New("refresh task already in progress")
	}

	// Set in progress to true
	t.inProgress = true

	// Create the proxy channel
	t.proxyChan = make(chan *Proxy)

	// Create stop channel
	t.stop = make(chan bool)

	// Start the proxy handler
	t.proxyHandler()

	// Start the refresh task
	t.startRefreshTask()

	return t.stop, err
}
