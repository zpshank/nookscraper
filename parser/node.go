package parser

import (
	"golang.org/x/net/html"
	"strings"
)

// The NodeParser provides supporting methods when working with
// the 'golang.org/x/net/html' package
type NodeParser struct{}

func (np *NodeParser) match(node *html.Node, key string, value string) bool {
	switch key {
	case "id":
		if node.Type == html.ElementNode {
			match, attr := np.GetAttribute(node, key)
			if match && attr.Val == value {
				return true
			}
		}
	case "tag":
		if node.Type == html.ElementNode && string(node.Data) == value {
			return true
		}
	case "html":
		if node.Type == html.ElementNode {
			child := node.FirstChild
			if child != nil && child.Type == html.TextNode && strings.TrimSpace(child.Data) == value {
				return true
			}
		}
	}
	return false
}

// The FindNode method will return the first child Node that matches
// the given criteria in the method. Will return false and nil if no Node matches
// the given criteria.
func (np *NodeParser) Find(root *html.Node, key string, value string) (bool, *html.Node) {
	if np.match(root, key, value) {
		return true, root
	}
	// Loop through node children until main content is found
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		match, rn := np.Find(c, key, value)
		if match {
			return true, rn
		}
	}
	return false, nil

}

// The FindSibling method will return the first sibling Node that matches
// the given criteria in the method. Will return false and nil if no Node matches
// the given criteria.
func (np *NodeParser) FindSibling(root *html.Node, key string, value string) (bool, *html.Node) {
	for n := root.NextSibling; n != nil; n = n.NextSibling {
		if np.match(n, key, value) {
			return true, n
		}
	}
	return false, nil
}

// The GetAttribute method will return the Attribute which matches the given key.
// Will return false and nil if no Attribute matches the given key.
func (np *NodeParser) GetAttribute(n *html.Node, key string) (match bool, attr *html.Attribute) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return true, &attr
		}
	}
	return false, nil
}
