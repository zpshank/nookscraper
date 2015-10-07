package node

import (
	"golang.org/x/net/html"
)

func FindNode(root *html.Node, key string, value string) (n *html.Node) {
	var find func(*html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		switch key {
		case "id":
			if n.Type == html.ElementNode {
				ok, attr := GetAttribute(n, key)
				if ok && attr.Val == value {
					return n
				}
			}
		case "tag":
			if n.Type == html.ElementNode && string(n.Data) == value {
				return n
			}
		}
		// Loop through node children until main content is found
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			rn := find(c)
			if rn != nil {
				return rn
			}
		}
		return nil
	}
	return find(root)
}

func GetAttribute(n *html.Node, key string) (ok bool, attr *html.Attribute) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return true, &attr
		}
	}
	return false, nil
}
