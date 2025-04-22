package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/ctiller15/crawler/internal"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error level code hit")
	}

	if !strings.Contains(resp.Header.Get("Content-type"), "text/html") {
		return "", fmt.Errorf("invalid content type, not text/html")
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(html), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()
	defer cfg.wg.Done()

	u, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if u.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	url, err := internal.NormalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	cfg.mu.Lock()
	if _, ok := cfg.pages[url]; ok {
		fmt.Println("found url, incrementing...")
		cfg.pages[url] += 1
		cfg.mu.Unlock()
		return
	} else {
		cfg.pages[url] = 1
		cfg.mu.Unlock()
	}

	finalUrl := fmt.Sprintf("%s://%s", cfg.baseURL.Scheme, url)
	fmt.Printf("getting html for %s...\n", finalUrl)
	html, err := getHTML(finalUrl)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("html found: %s\n", html)

	rawBaseURL := fmt.Sprintf("%s://%s", cfg.baseURL.Scheme, cfg.baseURL.Hostname())
	foundUrls, err := internal.GetURLsFromHTML(html, rawBaseURL)
	fmt.Println(foundUrls)
	for _, foundUrl := range foundUrls {
		fmt.Printf("crawling url %s...\n", foundUrl)
		cfg.wg.Add(1)
		go cfg.crawlPage(foundUrl)
	}

}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Print("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		log.Print("too many arguments provided")
		os.Exit(1)
	}

	baseUrl := args[0]
	log.Printf("starting crawl of: %s", baseUrl)

	parsedBase, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	pages := make(map[string]int)

	cfg := config{
		pages:              pages,
		baseURL:            parsedBase,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 5),
		wg:                 &sync.WaitGroup{},
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(baseUrl)
	cfg.wg.Wait()

	fmt.Println()
	fmt.Println("===Final report===")
	for key, val := range pages {
		fmt.Printf("%s: %d\n", key, val)
	}
}
