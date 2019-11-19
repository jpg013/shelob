package proxy

import (
	"bytes"
	"fmt"
	"log"
	"shelob/db"
	"strconv"
	"strings"

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

type Executor func(*html.Node) interface{}

type Pipeline interface {
	Pipe(executor Executor) Pipeline
	Merge() chan *html.Node
}

type Collector struct {
	executors []Executor
	dataCh    chan *html.Node
}

func (c *Collector) Merge() chan interface{} {
	for i := 0; i < len(c.executors); i++ {
		fn := c.executors[i]
		c.dataCh = run(fn, c.dataCh)
	}

	return c.dataCh
}

func run(fn Executor, inCh chan interface{}) chan interface{} {
	outCh := make(chan interface{})

	go func() {
		defer close(outCh)

		for n := range inCh {
			outCh <- fn(n)
		}
	}()

	return outCh
}

func (c *Collector) Pipe(fn Executor) Pipeline {
	c.executors = append(c.executors, fn)
	return c
}

func NewCollector(in chan interface{}) Pipeline {
	return &Collector{
		dataCh:    in,
		executors: make([]Executor, 0),
	}
}

func parsePort(n *html.Node) (port int) {
	// Extract the base64 image src for port image
	imgSrc := dom.ParseImageSrc(n)

	filter := bson.M{"hash_state": hashString(imgSrc), "port": bson.M{"$exists": true, "$gt": 0}}
	options := options.Find()

	// // Limit by 10 documents only
	options.SetLimit(1)
	docs, _ := db.Find("proxy_port_hash", filter, options)

	if len(docs) > 0 {
		fmt.Println(docs)
	}

	// process port image and extract txt
	txt, err := Base64ImgToText(imgSrc)

	if err != nil {
		log.Println("error converting base 64 image text: ", err.Error())
	}

	// convert text to int
	port, err = strconv.Atoi(txt)

	proxyPortHash := NewProxyPortHash(imgSrc, port)
	db.InsertOne("proxy_port_hash", proxyPortHash)

	return port
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

func parseIPAddress(n *html.Node, s *html.Node) (ipaddr string) {
	fs := make([]*IPFrag, 0)
	ts := dom.GetChildrenByType(n, html.TextNode)

	makeFrag := func(n *html.Node) {
		if n.FirstChild == nil {
			return
		}

		frag := &IPFrag{
			Styles:    dom.ParseStyleAttribute(n.FirstChild),
			Classname: dom.GetAttribute(n, "class"),
			Value:     n.FirstChild.Data,
		}

		fs = append(fs, frag)
	}

	for _, t := range ts {
		dom.IterateSiblings(t, makeFrag)
	}

	// construct the ip address from the fragments
	return buildIPAddress(fs, dom.ParseStyleNodeBody(s))
}

func parseLocation(n *html.Node) (loc string) {
	c := n.FirstChild

	if c != nil && c.NextSibling != nil {
		loc = strings.TrimSpace(c.NextSibling.Data)
	}

	return loc
}

func parseProtocol(n *html.Node) (protocol string) {
	c := n.FirstChild

	if c != nil {
		protocol = strings.TrimSpace(c.Data)
	}

	if protocol == "HTTP(S)" {
		return "https"
	}

	return "http"
}

func ScrapeProxy(tr *html.Node, ch chan<- *Proxy) {
	// Select all tds for tr
	tds := dom.QuerySelectorAll("td", tr)
	// select tr style
	nstyle := dom.QuerySelector("style", tr)

	if tds == nil {
		log.Println("could not query td elements")
		return
	}

	if nstyle == nil {
		log.Println("unable to find style node")
		return
	}

	// Map positional td nodes
	ipNode := tds[1]
	portNode := tds[2]
	locNode := tds[3]
	protocolNode := tds[5]

	port := parsePort(portNode)

	// If unable to parse port, continue
	if port == 0 {
		return
	}

	loc := parseLocation(locNode)
	protocol := parseProtocol(protocolNode)
	ip := parseIPAddress(ipNode, nstyle)

	if ip == "" {
		log.Println("unable to parse ip address")
		return
	}

	ch <- NewProxy(ip, port, protocol, loc)
}

func ScrapeProxyList(doc *html.Node, outCh chan<- *Proxy) {
	// Get table rows and loop over each tr in list
	// inCh := gen(getTableRows(doc))
	// c1 := doWork(inCh)
	// c2

	// for _, tr := range getTableRows(doc) {
	// 	go ScrapeProxy(tr, outCh)

	// 	// id, err := db.InsertOne("proxy", proxy)

	// 	// if err != nil {
	// 	// 	log.Fatal(err.Error())
	// 	// }

	// 	// log.SetOutput(os.Stdout)
	// 	// log.Println("inserted id: " + id)
	// }
}

func doWork(inCh <-chan *html.Node) <-chan string {
	outCh := make(chan string)

	go func() {
		for n := range inCh {
			tds := dom.QuerySelectorAll("td", n)

			// do some stuff
			outCh <- parseProtocol(tds[5])
		}
	}()

	return outCh
}

func gen(ns []*html.Node) <-chan *html.Node {
	outCh := make(chan *html.Node)

	go func() {
		for _, n := range ns {
			outCh <- n
		}
	}()

	return outCh
}
