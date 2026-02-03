package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	url := "http://whatismyip.akamai.com"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching URL: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "Error: received status code %d\n", resp.StatusCode)
		os.Exit(1)
	}

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		os.Exit(1)
	}
	fmt.Println()
}
