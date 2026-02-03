package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	ipv4Regex   = regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?))$`)
	serviceURLs = []string{
		"https://whatismyip.akamai.com",
		"https://checkip.amazonaws.com",
		"https://icanhazip.com",
		"https://api.ipify.org",
		"https://ident.me",
		"http://checkip.amazonaws.com",
		"http://icanhazip.com",
		"http://api.ipify.org",
		"http://ipecho.net/plain",
		"http://ident.me",
		"http://ipinfo.io/ip",
	}
)

func getIp(client *http.Client, url string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: received status code %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body from %s: %w", url, err)
	}

	ip := strings.TrimSpace(string(body))
	if !ipv4Regex.MatchString(ip) {
		return "", fmt.Errorf("invalid IPv4 address received from %s: %s", url, ip)
	}

	return ip, nil
}

func getIpMany(client *http.Client, urls ...string) (string, error) {
	var lastErr error
	for _, url := range urls {
		ip, err := getIp(client, url)
		if err == nil {
			return ip, nil
		}
		lastErr = err
	}
	return "", fmt.Errorf("failed to fetch IP from all sources: %w", lastErr)
}

func main() {
	ip, err := getIpMany(http.DefaultClient, serviceURLs...)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(ip)
}
