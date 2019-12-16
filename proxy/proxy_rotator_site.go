package proxy

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"shelob/db"
	"strconv"
	"strings"
	"time"

	dom "github.com/jpg013/go_dom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/html"
)

var (
	urls = []string{
		"https://www.proxyrotator.com/free-proxy-list/1",
		// "https://www.proxyrotator.com/free-proxy-list/2",
		// "https://www.proxyrotator.com/free-proxy-list/3",
		// "https://www.proxyrotator.com/free-proxy-list/4",
		// "https://www.proxyrotator.com/free-proxy-list/5",
		// "https://www.proxyrotator.com/free-proxy-list/6",
		// "https://www.proxyrotator.com/free-proxy-list/7",
		// "https://www.proxyrotator.com/free-proxy-list/8",
		// "https://www.proxyrotator.com/free-proxy-list/9",
		// "https://www.proxyrotator.com/free-proxy-list/10",
	}
)

// IPFrag represents a proxy rotator site data required to construct the
// entire ip address. The Value property includes multiple faux / invalid data items
// in the DOM to prevent scraping. Only by looking at the visible
// styles / classnames can we determine the actual IP Address.
type IPFrag struct {
	Value     string
	Styles    map[string]string
	Classname string
}

type proxyRotatorSite struct {
	requestThrottle *time.Ticker
	collector       Pipeline
	in              chan *proxyBuilder
	stop            chan bool
}

func (prs *proxyRotatorSite) doScrape(url string) {
	// Wait for request throttle
	<-prs.requestThrottle.C

	resp, err := http.Get(url)

	if err != nil {
		log.Println(fmt.Sprintf("error fetching proxy list: %v", err))
	} else {
		defer resp.Body.Close()
		doc, err := dom.ParseHTMLDocument(resp.Body)

		if err != nil {
			log.Fatal(fmt.Sprintf("error parsing html document: %v", err))
		}

		// Parse the table rows and send each item through the pipeline
		rows := getTableRows(doc)

		for _, n := range rows {
			data := make(map[string]interface{})
			data["source"] = "proxyrotator.com"

			prs.in <- &proxyBuilder{
				root:     n,
				siteName: "proxyrotator.com",
				data:     data,
			}
		}
	}
}

func (prs *proxyRotatorSite) Stop() error {
	// send a signal to the stop channel
	prs.stop <- true

	return nil
}

func (prs *proxyRotatorSite) Run(out chan<- *Proxy) error {
	// Call merge on the collector to create a producer channel
	in := prs.collector.Merge()

	// wait for collector results
	go func() {
		// Loop over data channel
		for p := range in {
			proxy, err := NewProxy(p.data)

			if err != nil {
				log.Fatalf("could not create proxy: %v", err.Error())
			}

			out <- proxy
		}
	}()

	// call make generator
	go prs.makeGenerator(out)

	return nil
}

func (prs *proxyRotatorSite) makeGenerator(out chan<- *Proxy) {
	// Create the request throttle for urls
	prs.requestThrottle = time.NewTicker(time.Second)

	// Make this configurable	// Make this configurable
	taskDelay := time.NewTicker(time.Minute * 2)
	defer taskDelay.Stop()

	// Generate function iterates through proxy url list and
	// scrapes proxies for each page
	generate := func() {
		fmt.Println("Starting generator")

		for _, url := range urls {
			// Wait for request throttle
			go prs.doScrape(url)
		}
	}

	// run initial generator
	generate()

	for {
		select {
		case <-taskDelay.C:
			generate()
		// listen for stop signal and break out of loop
		case <-prs.stop:
			return
		}
	}
}

// NewProxyRotatorSite factory returns a proxy rotator site
func NewProxyRotatorSite() Site {
	in := make(chan *proxyBuilder)
	collector := NewCollector(in).
		Pipe(selectTableRowData).
		Pipe(parseLocation).
		Pipe(parseProtocol).
		Pipe(parseIPAddress).
		Pipe(parsePortImgSource).
		Pipe(determinePortFromImageSource)

	return &proxyRotatorSite{
		in:        in,
		collector: collector,
		stop:      make(chan bool),
	}
}

func getTableRows(n *html.Node) (trs []*html.Node) {
	// find document body
	body := dom.GetDocumentBody(n)

	if body == nil {
		log.Println("could not get document body")
		return trs
	}

	// get table body for proxy list
	tbody := dom.QuerySelector("tbody", body)

	if tbody == nil {
		log.Println("could not get table body")
		return trs
	}

	// query all tr elements from table body
	trs = dom.QuerySelectorAll("tr", tbody)
	return trs
}

func selectTableRowData(p *proxyBuilder) (*proxyBuilder, error) {
	// Select all tds for tr
	tds := dom.QuerySelectorAll("td", p.root)

	// select tr style
	styles := dom.QuerySelector("style", p.root)

	if tds == nil {
		return p, errors.New("could not query td elements")
	}

	if styles == nil {
		return p, errors.New("unable to find style node")
	}

	p.tds = tds
	p.styles = styles

	return p, nil
}

func parseLocation(p *proxyBuilder) (*proxyBuilder, error) {
	// Location is the third table cell
	n := p.tds[3]

	if n == nil {
		return p, errors.New("could not get location node")
	}

	c := n.FirstChild

	if c != nil && c.NextSibling != nil {
		p.data["location"] = strings.TrimSpace(c.NextSibling.Data)
	} else {
		return p, errors.New("could not get location node")
	}

	return p, nil
}

func parseProtocol(p *proxyBuilder) (*proxyBuilder, error) {
	// protocol is the fifth table cell
	n := p.tds[5]

	if n == nil {
		return p, errors.New("could not get protocol node")
	}

	c := n.FirstChild

	if c == nil {
		return p, errors.New("could not get protocol node")
	}

	protocol := strings.TrimSpace(c.Data)

	if protocol == "HTTP(S)" {
		p.data["protocol"] = "https"
	} else {
		p.data["protocol"] = "http"
	}

	return p, nil
}

func parseIPAddress(p *proxyBuilder) (*proxyBuilder, error) {
	// ip address node is the first table cell
	n := p.tds[1]

	if n == nil {
		return p, errors.New("could not get ip address node")
	}

	fs := make([]*IPFrag, 0)
	ts := dom.GetChildrenByType(n, html.TextNode)

	makeFrag := func(n *html.Node) {
		if n.FirstChild == nil {
			return
		}

		frag := &IPFrag{
			Styles:    dom.ParseStyleAttribute(n),
			Classname: dom.GetAttribute(n, "class"),
			Value:     n.FirstChild.Data,
		}

		fs = append(fs, frag)
	}

	for _, t := range ts {
		dom.IterateSiblings(t, makeFrag)
	}

	// construct the ip address from the fragments
	p.data["ipAddress"] = buildIPAddress(fs, dom.ParseStyleNodeBody(p.styles))

	return p, nil
}

func buildIPAddress(fs []*IPFrag, styleMap dom.StyleDeclarations) string {
	buf := bytes.Buffer{}

	for _, f := range fs {
		inlineStyle := f.Styles["display"]
		globalStyle := ""

		if styleData, ok := styleMap[f.Classname]; ok {
			globalStyle = styleData["display"]
		}

		if inlineStyle == "inline" || globalStyle == "inline" {
			buf.WriteString(f.Value)
		}
	}

	return buf.String()
}

func parsePortImgSource(p *proxyBuilder) (*proxyBuilder, error) {
	n := p.tds[2]

	if n == nil {
		return p, errors.New("could not get port image source node")
	}

	p.data["portImageSource"] = dom.ParseImageSrc(n)
	p.data["portHashState"] = hashString(p.data["portImageSource"].(string))

	return p, nil
}

func determinePortFromImageSource(p *proxyBuilder) (*proxyBuilder, error) {
	// zero port
	p.data["port"] = 0

	port, _ := lookupPortByHashState(p.data["portHashState"].(string))

	if port > 0 {
		p.data["port"] = port
		return p, nil
	}

	// process port image and extract txt
	taskKey, err := Base64ImgToText(p.data["portImageSource"].(string))

	if err != nil {
		return p, err
	}

	// wait for task to be completed
	portStr := <-waitPortOCRTask(taskKey)

	// convert port image characters to int
	port, err = strconv.Atoi(portStr)

	if err == nil {
		// Weird, converting back from int to string
		p.data["port"] = port
	}

	filter := &bson.D{{Key: "hash_state", Value: p.data["portHashState"]}}
	update := &bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "port", Value: p.data["port"]},
			{Key: "hash_state", Value: p.data["portHashState"]},
			{Key: "base64_image", Value: p.data["portImageSource"]},
		}}}

	// insert proxy port hash
	db.UpdateOne("proxy_port_hash", filter, update)
	return p, nil
}

func lookupPortByHashState(hash string) (int, error) {
	filter := bson.M{"hash_state": hash, "port": bson.M{"$exists": true, "$gt": 0}}
	options := options.Find()

	// Limit by 1 document only
	options.SetLimit(1)
	docs, err := db.Find("proxy_port_hash", filter, options)

	if err != nil {
		return 0, nil
	}

	if len(docs) > 0 {
		// assert value as integer
		return int(docs[0]["port"].(int32)), nil
	}

	return 0, nil
}

func waitPortOCRTask(taskKey string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

		// wait for 5 minutes
		count := 300

		for i := 0; i < count; i++ {
			time.Sleep(time.Second * 1)
			portStr, err := db.GetCache(taskKey)

			if err != nil {
				break
			}

			if portStr != "" {
				out <- portStr
				break
			}
		}
	}()

	return out
}
