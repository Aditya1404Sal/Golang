package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}
type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	websiteUrl := flag.String("url", "https://gophercises.com/", "input url for a website which you want to convert to xml")
	depth := flag.Int("depth", 5, "the number of times bfs gets implemented")
	flag.Parse()
	pages := bfsAlgo(*websiteUrl, *depth)
	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()

	// resp, err := http.Get(*websiteUrl)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// r := resp.Body
	// reqUrl := resp.Request.URL
	// baseUrl := &url.URL{
	// 	Scheme: reqUrl.Scheme,
	// 	Host:   reqUrl.Host,
	// }
	// base := baseUrl.String()
	// links, _ := Parse(r)
	// var hrefs []string

	// for _, l := range links {
	// 	switch {
	// 	case strings.HasPrefix(l.Href, "/"):
	// 		hrefs = append(hrefs, base+l.Href)
	// 	case strings.HasPrefix(l.Href, "http"):
	// 		hrefs = append(hrefs, l.Href)
	// 	}
	// }
	// var filteredhrefs []string
	// for _, href := range hrefs {
	// 	if strings.HasPrefix(href, base) {
	// 		filteredhrefs = append(filteredhrefs, href)
	// 	}
	// }

	// for _, fhref := range filteredhrefs {
	// 	fmt.Println(fhref)
	// }

}

type Link struct {
	Href string
	Text string
}

type emptyStruct struct {
}

func bfsAlgo(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: emptyStruct{},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = emptyStruct{}
			for _, link := range get(url) {
				nq[link] = emptyStruct{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}
func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

// Parse will take in an HTML document and will return a
// slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = text(n)
	return ret
}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}
