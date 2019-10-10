package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"

// 	dom "github.com/jpg013/go_dom"
// 	"golang.org/x/net/html"
// )

// type Proxy struct {
// 	IPAddress string
// 	Port      int
// 	Protocols []string
// 	Location  string
// 	Speed     int8
// }

// type StyleDeclarations map[string]map[string]string

// var base64Prefixes = []string{
// 	"data:image/jpeg;charset=utf-8;base64,",
// 	"data:image/png;charset=utf-8;base64",
// }

// type Base64Payload struct {
// 	Base64String string `json:"base64_string,omitempty"`
// 	Extension    string `json:"extension,omitempty"`
// }

// func ParseProxyList(doc *html.Node) []*Proxy {
// 	// make slice of proxy addresses
// 	ps := make([]*Proxy, 0)

// 	// find document body
// 	body := dom.GetDocumentBody(doc)
// 	// get table body for proxy list
// 	tbody := dom.QuerySelector("tbody", body)
// 	// pull all trs from table body
// 	trs := dom.QuerySelectorAll("tr", tbody)

// 	for _, tr := range trs {
// 		// Select all tds for table row
// 		tds := dom.QuerySelectorAll("td", tr)

// 		// Extract image data for port
// 		imgSrc := dom.ParseImageSrc(tds[2])

// 		// Replace prefixes
// 		for _, p := range base64Prefixes {
// 			imgSrc = strings.Replace(imgSrc, p, "", -1)
// 		}

// 		// convert the base64 image source to text
// 		txt := Base64SrcToText(imgSrc, "png")

// 		// Attempt to convert text to int
// 		port, err := strconv.Atoi(txt)

// 		if err != nil {
// 			// we were unsuccessful to convert port, save the image / text for manual review
// 			f, err := os.OpenFile("unknown_proxy_ip_ports.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 			if err != nil {
// 				panic(err)
// 			}

// 			defer f.Close()

// 			str := fmt.Sprintf("%v:%v\n", txt, imgSrc)

// 			f.WriteString(str)
// 			continue
// 		}

// 		styleData := dom.ParseStyleData(dom.QuerySelector("style", tr))
// 		addrNodes := GetTextNodeChildren(tds[1])

// 		frags := make([]*IPFrag, 0)
// 		for _, a := range addrNodes {
// 			frags = append(frags, IterateSiblings(a)...)
// 		}

// 		ip := constructIPAddressFromFragments(frags, styleData)

// 		proxy := &Proxy{
// 			IPAddress: ip,
// 			Port:      port,
// 		}

// 		ps = append(ps, proxy)
// 	}

// 	// ping each proxy to make sure it's alive.
// 	for _, p := range ps {
// 		fmt.Println(fmt.Sprintf("%v:%v", p.IPAddress, p.Port))
// 	}

// 	return ps
// }

// func constructIPAddressFromFragments(fs []*IPFrag, styleMap StyleDeclarations) string {
// 	buf := bytes.Buffer{}

// 	for _, f := range fs {
// 		inlineStyle := f.Styles["display"]
// 		globalStyle := ""

// 		if styleData, ok := styleMap[f.Classname]; ok {
// 			globalStyle = styleData["display"]
// 		}

// 		if inlineStyle != "inline" && globalStyle != "inline" {
// 			continue
// 		}

// 		buf.WriteString(f.Value)
// 	}

// 	return buf.String()
// }

// type IPFrag struct {
// 	Value     string
// 	Styles    map[string]string
// 	Classname string
// }

// func parseStyleBody(s string) map[string]string {
// 	// Split on semicolon separators
// 	split := strings.Split(s, ";")
// 	m := make(map[string]string)

// 	for _, body := range split {
// 		parts := strings.Split(strings.TrimSpace(body), ":")

// 		if len(parts) != 2 {
// 			continue
// 		}

// 		m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
// 	}

// 	return m
// }

// func IterateSiblings(n *html.Node) []*IPFrag {
// 	res := make([]*IPFrag, 0)

// 	for c := n; c != nil; c = c.NextSibling {
// 		if c.FirstChild != nil {
// 			res = append(res, &IPFrag{
// 				Styles:    parseStyleBody(dom.GetAttribute(c, "style")),
// 				Classname: dom.GetAttribute(c, "class"),
// 				Value:     c.FirstChild.Data,
// 			})
// 		}
// 	}

// 	return res
// }

// func GetTextNodeChildren(n *html.Node) []*html.Node {
// 	ns := make([]*html.Node, 0)

// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		if c.Type == html.TextNode {
// 			ns = append(ns, c)
// 		}
// 	}

// 	return ns
// }

// func Base64SrcToText(base64Src string, ext string) string {
// 	requestBody, err := json.Marshal(&Base64Payload{Base64String: base64Src, Extension: ext})

// 	if err != nil {
// 		panic("invalid request body")
// 	}

// 	// Convert base64 png image to text
// 	resp, err := http.Post("http://localhost:5000/base64_source_to_text", "application/json", bytes.NewBuffer(requestBody))

// 	if err != nil {
// 		panic("cannot convert base64 png to text")
// 	}

// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)

// 	if err != nil {
// 		panic("cannot read response body")
// 	}

// 	return string(body)
// }
