package internal

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	urls := make([]string, 0)

	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		return nil, err
	}

	err = traverse(doc, rawBaseURL, &urls)
	if err != nil {
		return nil, err
	}

	return urls, nil
}

func traverse(node *html.Node, rawBaseURL string, urls *[]string) error {
	if node.Type == html.ElementNode && node.DataAtom == atom.A {
		for _, a := range node.Attr {
			if a.Key == "href" {
				if strings.TrimSpace(a.Val) == "" {
					continue
				}

				hrefURL, err := url.Parse(a.Val)
				if err != nil {
					continue
				}

				base, err := url.Parse(rawBaseURL)
				if err != nil {
					return err
				}

				resolvedURL := base.ResolveReference(hrefURL)

				if resolvedURL.Fragment != "" && resolvedURL.Path == "" {
					continue
				}

				*urls = append(*urls, resolvedURL.String())
			}
		}
	} else {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if err := traverse(child, rawBaseURL, urls); err != nil {
				return err
			}
		}
	}
	return nil
}
