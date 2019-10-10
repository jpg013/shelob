package proxy

// import (
// 	"errors"
// 	"time"
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

// // ProxyRotator struct
// type ProxyRotator struct {
// 	refreshTicker bool
// }

// func NewProxyRotator() *ProxyRotator {
// 	return &ProxyRotator{
// 		refreshInProgress: false,
// 	}
// }

// func (p *ProxyRotator) refreshProxyList() {

// }

// func (p *ProxyRotator) ScheduleRefreshTask(delay time.Duration) (<-chan bool, error) {
// 	if p.refreshInProgress == true {
// 		return nil, errors.New("refresh task already scheduled")
// 	}

// 	stop := make(<-chan bool)
// 	p.refreshInProgress = true

// 	go func() {
// 		p.refreshProxyList()
// 		for {
// 			select {
// 			case <-time.After(delay):
// 				p.refreshProxyList()
// 			case <-stop:
// 				return
// 			}
// 		}
// 	}()

// 	return stop, nil
// }
