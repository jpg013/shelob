package proxy

import (
	"bytes"
	"errors"
	"log"
	"shelob/db"
	"strconv"
	"strings"
	"sync"
	"time"

	dom "github.com/jpg013/go_dom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/html"
)

type IPFrag struct {
	Value     string
	Styles    map[string]string
	Classname string
}

type ProxyBuilder struct {
	root            *html.Node
	tds             []*html.Node
	styles          *html.Node
	location        string
	protocol        string
	ipAddress       string
	portImageSource string
	portHashState   string
	port            int
}

type Executor func(*ProxyBuilder) (*ProxyBuilder, error)

type Predicate func(*ProxyBuilder) bool

type Pipeline interface {
	Pipe(executor Executor) Pipeline
	Merge() chan *ProxyBuilder
}

type Collector struct {
	executors []Executor
	dataCh    chan *ProxyBuilder
}

func (c *Collector) Merge() chan *ProxyBuilder {
	for i := 0; i < len(c.executors); i++ {
		fn := c.executors[i]
		c.dataCh = runExecutor(fn, c.dataCh)
	}

	return c.dataCh
}

func (c *Collector) Pipe(fn Executor) Pipeline {
	c.executors = append(c.executors, fn)
	return c
}

func NewCollector(in chan *ProxyBuilder) Pipeline {
	return &Collector{
		dataCh:    in,
		executors: make([]Executor, 0),
	}
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

func convertPortImageSrc(p *ProxyBuilder) (*ProxyBuilder, error) {
	if p.port > 0 {
		return p, nil
	}

	return p, nil
}

func parsePortImgSource(p *ProxyBuilder) (*ProxyBuilder, error) {
	n := p.tds[2]

	if n == nil {
		return p, errors.New("could not get port image source node")
	}

	p.portImageSource = dom.ParseImageSrc(n)
	p.portHashState = hashString(p.portImageSource)

	return p, nil
}

func waitPortOCRTask(taskKey string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)

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

func processPortImageSource(p *ProxyBuilder) (*ProxyBuilder, error) {
	port, _ := lookupPortByHashState(p.portHashState)

	if port > 0 {
		p.port = port
		return p, nil
	}

	// process port image and extract txt
	taskKey, err := Base64ImgToText(p.portImageSource)

	if err != nil {
		return p, err
	}

	// wait for task to be completed
	portStr := <-waitPortOCRTask(taskKey)

	// convert port image characters to int
	port, err = strconv.Atoi(portStr)

	if err == nil {
		p.port = port
	}
	filter := &bson.D{{"hash_state", p.portHashState}}
	update := &bson.D{
		{"$set", bson.D{
			{"port", p.port},
			{"hash_state", p.portHashState},
			{"base64_image", p.portImageSource},
		}}}

	// insert proxy port hash
	db.UpdateOne("proxy_port_hash", filter, update)
	return p, nil
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

func selectTableRowData(p *ProxyBuilder) (*ProxyBuilder, error) {
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

func parseIPAddress(p *ProxyBuilder) (*ProxyBuilder, error) {
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
	p.ipAddress = buildIPAddress(fs, dom.ParseStyleNodeBody(p.styles))

	return p, nil
}

func parseLocation(p *ProxyBuilder) (*ProxyBuilder, error) {
	// Location is the third table cell
	n := p.tds[3]

	if n == nil {
		return p, errors.New("could not get location node")
	}

	c := n.FirstChild

	if c != nil && c.NextSibling != nil {
		p.location = strings.TrimSpace(c.NextSibling.Data)
	} else {
		return p, errors.New("could not get location node")
	}

	return p, nil
}

func parseProtocol(p *ProxyBuilder) (*ProxyBuilder, error) {
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
		p.protocol = "https"
	} else {
		p.protocol = "http"
	}

	return p, nil
}

func ScrapeProxyList(doc *html.Node, outCh chan<- *Proxy) {
	// Get table rows and loop over each tr in list
	ch := NewCollector(producer(doc)).
		Pipe(selectTableRowData).
		Pipe(parseLocation).
		Pipe(parseProtocol).
		Pipe(parseIPAddress).
		Pipe(parsePortImgSource).
		Pipe(processPortImageSource).
		Merge()

	for p := range ch {
		if p.port == 0 {
			continue
		}

		proxy := NewProxy(p.ipAddress, p.port, p.protocol, p.location)

		go func(proxy *Proxy) {
			VerifyProxy(proxy)
		}(proxy)
	}
}

func runExecutor(fn Executor, in <-chan *ProxyBuilder) chan *ProxyBuilder {
	out := make(chan *ProxyBuilder)

	go func() {
		var wg sync.WaitGroup

		for p := range in {
			wg.Add(1)
			go func(p *ProxyBuilder) {
				val, _ := fn(p)
				out <- val
				wg.Done()
			}(p)
		}

		go func() {
			defer close(out)
			wg.Wait()
		}()
	}()

	return out
}

func producer(n *html.Node) chan *ProxyBuilder {
	outCh := make(chan *ProxyBuilder)

	go func() {
		defer close(outCh)
		for _, n := range getTableRows(n) {
			outCh <- &ProxyBuilder{root: n}
		}
	}()

	return outCh
}
