package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {

	// Function to check if we are passing any args (Atleast one args)
	args := os.Args

	if len(args) < 2 {
		fmt.Printf("Not Proper Elements")
		os.Exit(1)
	}

	fmt.Printf("Args:%v",args[1:])

	//https://10.26.10.44:8889/app/health-check
	if _, err := url.ParseRequestURI(args[1]); err != nil {
		fmt.Printf("Invalid URL %v\n", err)
		os.Exit(1)
	}

	response, err := http.Get(args[1])

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