package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"shelob/db"
	"shelob/proxy"
	"time"

	"github.com/joho/godotenv"
)

func makeRequest() {
	// proxyUrl, err := url.Parse("http://35.226.112.130:3128")

	client := &http.Client{}
	// client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	// client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	req, _ := http.NewRequest("GET", "https://www.youtube.com/channel/UCO1cgjhGzsSYb1rsB4bFe4Q/videos", nil)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("unable to make request to youtube channel")
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
}

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}

	db.LoadConfig()
	db.OpenConnection()
	db.ExampleNewClient()
}

func main() {
	// makeRequest()
	refreshTask := proxy.NewProxyRefreshTask()
	_, err := refreshTask.Start(2 * time.Minute)

	if err != nil {
		log.Fatal("error starting proxy refresh task: ", err.Error())
	}

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":9000", nil)

	// proxyURL, err := url.Parse("http://46.247.58:3130")

	// if err != nil {
	// 	panic(err)
	// }

	// client := &http.Client{
	// 	Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	// 	Timeout:   time.Duration(5 * time.Second),
	// }

	// req, _ := http.NewRequest("GET", "https://www.youtube.com/channel/UC-3jIAlnQmbbVMV6gR7K8aQ", nil)

	// req.Header.Set("Content-type", "application/json")
	// req.Header.Set("User-Agent", " Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")

	// resp, err := client.Do(req)

	// if err != nil {
	// 	panic(err.Error())
	// }

	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)

	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Println(string(body))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
