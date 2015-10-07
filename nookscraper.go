package main

import (
	"fmt"
	"golang.org/x/net/html"
	//"io/ioutil"
	"net/http"
	//"strings"
	"github.com/zackshank/nookscraper/helper/node"
)

func FindVillagerListNode(r *http.Response) *html.Node {
	doc, err := html.Parse(r.Body)
	if err != nil {
		fmt.Println("There was an error parsing the response: ", err.Error())
		return nil
	}
	mc := node.FindNode(doc, "id", "mw-content-text")
	var vl *html.Node = node.FindNode(mc, "tag", "tbody")
	return vl
}

func main() {
	resp, err := http.Get("https://nookipedia.com/wiki/List_of_villagers")
	if err != nil {
		fmt.Println("There was an error with the request: ", err.Error())
	}

	defer resp.Body.Close()

	var vl *html.Node = FindVillagerListNode(resp)
	fmt.Println("Found Villager List: ", vl.Data)
}
