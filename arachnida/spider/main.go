package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Config struct {
	Recursive bool
	Depth     int
	Path      string
	URL       string
}

func parseDepth(val string) (int, error) {
	depth, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("depth must be a number, got %q", val)
	}
	if depth <= 0 {
		return 0, fmt.Errorf("depth must be greater than 0, got %d", depth)
	}
	return depth, nil
}

func parseArgs(args []string) (Config, error) {
	cfg := Config{Depth: 5, Path: "./data/"}
	depthSet := false

	if len(args) < 2 {
		return cfg, fmt.Errorf("missing URL")
	}

	last := len(args) - 1

	for i := 1; i < last; i++ {
		switch args[i] {
		case "-r":
			cfg.Recursive = true

		case "-l":
			if i+1 >= last {
				return cfg, fmt.Errorf("-l needs a depth value")
			}
			depth, err := parseDepth(args[i+1])
			if err != nil {
				return cfg, err
			}
			cfg.Depth = depth
			depthSet = true
			i++

		case "-p":
			if i+1 >= last {
				return cfg, fmt.Errorf("-p needs a path")
			}
			cfg.Path = args[i+1]
			i++

		case "-rlp":
			if i+2 >= last {
				return cfg, fmt.Errorf("-rlp needs a depth and a path: -rlp N PATH URL")
			}
			depth, err := parseDepth(args[i+1])
			if err != nil {
				return cfg, err
			}
			cfg.Recursive = true
			cfg.Depth = depth
			depthSet = true
			cfg.Path = args[i+2]
			i += 2

		default:
			return cfg, fmt.Errorf("unknown option %q", args[i])
		}
	}

	if depthSet && !cfg.Recursive {
		return cfg, fmt.Errorf("-l can only be used with -r")
	}

	raw := args[last]
	if strings.HasPrefix(raw, "-") {
		return cfg, fmt.Errorf("missing URL (last argument is %q)", raw)
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return cfg, fmt.Errorf("invalid URL %q (use http:// or https://)", raw)
	}
	cfg.URL = u.String()

	return cfg, nil
}

func fetchPage(client *http.Client, pageURL string) (*html.Node, error) {
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: bad status %s", pageURL, resp.Status)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return nil, fmt.Errorf("%s: not an HTML page", pageURL)
	}
	return html.Parse(resp.Body)
}

func walk(n *html.Node, images *[]string, links *[]string) {
	if n.Type == html.ElementNode {
		if n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					*images = append(*images, a.Val)
				}
			}
		} else if n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					*links = append(*links, a.Val)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, images, links)
	}
}

// n.Data  .jpg/jpeg, .png, .gif, .bmp

func crawl(cfg Config, client *http.Client) {
	depth := 0
	if cfg.Recursive == false {
		depth = 1
	}
	start, _ := url.Parse(cfg.URL)
	visited := make(map[string]bool)
	current := []string{cfg.URL}
	for d := 0; d <= depth && len(current) > 0; d++ {
		var next []string
		for _, pageURL := range current {
			if visited[pageURL] {
				continue
			}
			visited[pageURL] = true

			doc, err := fetchPage(client, pageURL)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				continue
			}
			var images, links []string
			walk(doc, &images, &links)

			base, _ := url.Parse(pageURL)

			for _, img := range images {
				//resolve with base.Parse (like links below),
				//check the extension, download into cfg.Path
				fmt.Println("image:", img) //test
			}
			for _, link := range links {
				u, err := base.Parse(link)
				if err != nil {
					continue
				}
				u.Fragment = ""
				if u.Host != start.Host || (u.Scheme != "http" && u.Scheme != "https") {
					continue
				}
				next = append(next, u.String())
			}
		}
		current = next
	}
}

func main() {
	cfg, err := parseArgs(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Usage: ./spider [-r] [-l N] [-p PATH] URL")
		os.Exit(1)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	crawl(cfg, client)

	// fmt.Println("images:", images)
	// fmt.Println("links:", links)
	// fmt.Printf("%+v\n", cfg)
}
