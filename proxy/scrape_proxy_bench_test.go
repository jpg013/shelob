package proxy

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	dom "github.com/jpg013/go_dom"
	"golang.org/x/net/html"
)

var rows []*html.Node

func setup() {
	file, err := os.Open("proxy-rotator-list.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	doc, err := dom.ParseHTMLDocument(file)

	if err != nil {
		log.Fatal(err)
	}

	// find document body
	body := dom.GetDocumentBody(doc)

	if body == nil {
		log.Fatal("could not get document body")
	}

	// get table body for proxy list
	tbody := dom.QuerySelector("tbody", body)

	if tbody == nil {
		log.Fatal("could not get table body")
	}

	// query all tr elements from table body
	rows = dom.QuerySelectorAll("tr", tbody)
}

func shutdown() {
	fmt.Println("shutdown!")
}

// // Map positional td nodes
// ipNode := tds[1]
// portNode := tds[2]
// locNode := tds[3]
// protocolNode := tds[5]

func GetLocation(p *PipeArgs) *PipeArgs {
	// Select all tds for tr
	n := dom.QuerySelectorAll("td", p.node)[3]
	c := n.FirstChild

	if c != nil && c.NextSibling != nil {
		p.location = strings.TrimSpace(c.NextSibling.Data)
	}

	return p
}

func GetPort(p *PipeArgs) *PipeArgs {
	// Select all tds for tr
	n := dom.QuerySelectorAll("td", p.node)[2]
	p.port = dom.ParseImageSrc(n)

	return p
}

func dataSource(dataCh chan<- *PipeArgs) {
	defer close(dataCh)
	for i := 0; i < len(rows); i++ {
		dataCh <- &PipeArgs{node: rows[i]}
	}
}

func BenchmarkPipeline(b *testing.B) {
	inCh := make(chan *PipeArgs)
	go dataSource(inCh)

	outCh := NewCollector(inCh).
		Pipe(GetLocation).
		Merge()

	for range outCh {
		// Do nothing, just for  drain out channel
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
