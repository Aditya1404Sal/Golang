package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html"
)

func main() {
	htmlFilename := flag.String("html", "example1.html", "valid html filename for extracting links from")
	flag.Parse()
	htmlOpened, err := os.Open(*htmlFilename)
	if err != nil {
		fmt.Printf("trouble opening file %s", *htmlFilename)
		os.Exit(1)
	}

	htmlReadable := io.Reader(htmlOpened)

	processedDoc, err := html.Parse(htmlReadable)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	Htmldecoder(processedDoc)
}

func Htmldecoder(n *html.Node) { // here we treat entire html tag as head node
	if n.Type == html.ElementNode && n.Data == "a" { //make sure that node is a tag & not plain text
		for _, a := range n.Attr { //n. data represents the tag we are searching for
			if a.Key == "href" { // since Attributes of html tag are key/value pairs A.K.A maps , we specify the attribute key we are looking for
				fmt.Println(a.Val) // printing the value from the key/value pair
				break
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		Htmldecoder(child)
	} //searching algorithm for the next nodes , sine the parser follows breadth first approach we search for sibling nodes then move to next nodes (deeper ones)
}
