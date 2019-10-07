package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	dom "github.com/jpg013/go_dom"
)

var base64Prefixes = []string{
	"data:image/jpeg;charset=utf-8;base64,",
	"data:image/png;charset=utf-8;base64",
}

type Base64Payload struct {
	Base64String string `json:"base64_string,omitempty"`
	ImageType    string `json:"image_type,omitempty"`
}

func convertBase64SrcToText(base64Src string, imageType string) string {
	requestBody, err := json.Marshal(&Base64Payload{Base64String: base64Src, ImageType: imageType})

	if err != nil {
		panic("invalid request body")
	}

	// Convert base64 png image to text
	resp, err := http.Post("http://localhost:5000/base64_source_to_text", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		panic("cannot convert base64 png to text")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic("cannot read response body")
	}

	return string(body)
}

func refreshProxyList() {
	resp, err := http.Get("https://www.proxyrotator.com/free-proxy-list/1/#free-proxy-list")

	if err != nil {
		panic("error getting proxy list")
	}

	defer resp.Body.Close()

	doc, err := dom.ParseHTMLDocument(resp.Body)

	if err != nil {
		panic("error parsing document")
	}

	body := dom.GetDocumentBody(doc)
	tbody := dom.QuerySelector("tbody", body)
	trs := dom.QuerySelectorAll("tr", tbody)

	for _, tr := range trs {
		// Select all tds for row
		tds := dom.QuerySelectorAll("td", tr)

		// Extract image data for port number
		imgSrc := dom.ParseImageSrc(tds[2])

		for _, p := range base64Prefixes {
			imgSrc = strings.Replace(imgSrc, p, "", -1)
		}

		txt := convertBase64SrcToText(imgSrc, "png")
		txt2 := convertBase64SrcToText(imgSrc, "jpeg")

		port, err := strconv.Atoi(txt)
		port2, err := strconv.Atoi(txt)

		fmt.Println(port)
		fmt.Println(port2)

		style := dom.QuerySelector("style", tr)

		log.Println(style.FirstChild.Data)
	}

	// style := dom.QuerySelector("style", trs[0])

	// Css
	// fmt.Println(style.FirstChild.Data)

	// 3 - Text
	// 1 - Element

	// Last Checked
	// fmt.Println(tds[0].FirstChild.Data)
	// IP Address
	// fmt.Println(tds[1])
	// Port
	// parser(tds[2])
	// fmt.Println(tds[2].FirstChild.NextSibling.Data)
	// fmt.Println(tds[2].Type == html.ElementNode)
	// fmt.Println(tds[2].Data)
	// fmt.Println(tds[2].FirstChild.FirstChild)

	// x := dom.QuerySelectorAll("tr", body)
	// y := x[0]

	// fmt.Println(y)

	// fmt.Println(doc)

	// builder := dom.
	//  NewBuilder().
	//  QuerySelector("tbody").Build

	// result := builder.Build(body)

	// fmt.Println(result)
	// QuerySelectorAll("tr")

	// builder.Build(body)

	// val := builder.Join()

	// fmt.Println(val)
}

func main() {
	refreshProxyList()
}
