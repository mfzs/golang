package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {

	var (
		requestURL string
		password string
		parsedURL *url.URL
		err error
	)

	flag.StringVar(&requestURL,"url", "", "url to access")
	flag.StringVar(&password,"password", "", "password to access")
	flag.Parse()
	//https://10.26.10.44:8889/app/health-check
	if parsedURL, err := url.ParseRequestURI(requestURL); err != nil {
		fmt.Printf("Invalid URL %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	response, err := http.Get(requestURL)

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("HTTP status Code: %d\nBody: %s", response.StatusCode, body)
}