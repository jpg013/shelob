package proxy

import (
	"bytes"
	"log"
	"strconv"
	"strings"

	dom "github.com/jpg013/go_dom"
	"golang.org/x/net/html"
)

type IPFrag struct {
	Value     string
	Styles    map[string]string
	Classname string
}

func parseProxyPort(n *html.Node) int {
	// // Extract the base64 image src for port image
	imgSrc := dom.ParseImageSrc(n)

	// Save image to disk
	filePath, err := SaveBase64ImageToDisk(imgSrc)

	// Always cleanup image
	defer DeleteImage(filePath)

	if err != nil {
		log.Fatal(err)
	}

	txt, err := Base64ImgToText(filePath)

	if err != nil {
		log.Fatal(err)
	}

	// Attempt to convert text to int
	port, err := strconv.Atoi(txt)

	if err != nil || port == 0 {
		log.Println("unable to convert imgSrc to port " + txt)
	}

	return port
}

func getTableRows(n *html.Node) []*html.Node {
	// find document body
	body := dom.GetDocumentBody(n)

	// get table body for proxy list
	tbody := dom.QuerySelector("tbody", body)

	// pull all trs from table body
	return dom.QuerySelectorAll("tr", tbody)
}

func constructIPAddressFromFragments(fs []*IPFrag, styleMap dom.StyleDeclarations) string {
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

	// construc the ip address from the fragments
	return constructIPAddressFromFragments(fs, dom.ParseStyleNodeBody(s))
}

func parseProxyLocation(n *html.Node) (loc string) {
	c := n.FirstChild

	if c != nil && c.NextSibling != nil {
		loc = strings.TrimSpace(c.NextSibling.Data)
	}

	return loc
}

func parseProxyProtocol(n *html.Node) (protocol string) {
	c := n.FirstChild

	if c != nil {
		protocol = strings.TrimSpace(c.Data)
	}

	return protocol
}

func ScrapeProxyRotatorList(doc *html.Node) {
	// Get table rows && loop over each tr in list
	for _, tr := range getTableRows(doc) {
		// Select all tds for tr
		tds := dom.QuerySelectorAll("td", tr)

		port := parseProxyPort(tds[2])

		// If unable to parse port, continue
		if port == 0 {
			continue
		}

		loc := parseProxyLocation(tds[3])
		protocol := parseProxyProtocol(tds[5])
		nstyle := dom.QuerySelector("style", tr)

		if nstyle == nil {
			log.Println("unable to find style node")
			continue
		}

		ipaddr := parseIPAddress(tds[1], nstyle)

		if ipaddr == "" {
			log.Println("unable to parse ip address")
			continue
		}

		proxy := NewProxy(ipaddr, port, protocol, loc)
		proxy.Insert()
	}
}
