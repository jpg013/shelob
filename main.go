package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"shelob/db"
	"shelob/proxy"
	"time"

	"github.com/joho/godotenv"
)

func makeRequest() {
	proxyUrl, err := url.Parse("http://35.226.112.130:3128")

	if err != nil {
		panic(err)
	}

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	req, _ := http.NewRequest("GET", "https://www.youtube.com/channel/UC-3jIAlnQmbbVMV6gR7K8aQ", nil)
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
		log.Print("No .env file found")
	}
}

func main() {
	// Open db connection
	db.OpenConnection()

	refreshTask := proxy.NewProxyRefreshTask()
	quitCh, err := refreshTask.Start(time.Minute)

	if err != nil {
		log.Fatal("error starting proxy refresh task: ", err.Error())
	}

	log.Println(quitCh)

	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
