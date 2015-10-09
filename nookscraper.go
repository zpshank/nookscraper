package main

import (
	"fmt"
	"github.com/zackshank/nookscraper/parser"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"time"
)

const (
	SiteRoot = "https://nookipedia.com"
	ImgDir   = "images"
)

func MakeDir() {
	// Make root images directory
	os.Mkdir(ImgDir, 0777)

	// Make images subdirectories
	os.MkdirAll(fmt.Sprint(ImgDir, "/villagers"), 0777)
}

func FindVillagerListNode(r *http.Response) *html.Node {
	np := parser.NodeParser{}
	doc, err := html.Parse(r.Body)
	if err != nil {
		fmt.Println("There was an error parsing the response: ", err.Error())
		return nil
	}
	_, mc := np.Find(doc, "id", "mw-content-text")
	_, vl := np.Find(mc, "tag", "table")
	return vl
}

func main() {
	np, vp := parser.NodeParser{}, VillagerParser{}

	MakeDir()

	var vl []*Villager

	resp, err := http.Get(fmt.Sprintf("%s/wiki/List_of_villagers", SiteRoot))
	if err != nil {
		fmt.Println("There was an error with the request: ", err.Error())
	}

	defer resp.Body.Close()

	var vl *html.Node = FindVillagerListNode(resp)

	_, tr := np.Find(vl, "tag", "tr")
	for i := 0; i < 2; i++ {
		_, tr = np.FindSibling(tr, "tag", "tr")
	}

	var id int = 1
	for _, v := vp.Parse(tr); v != nil; _, v = vp.Parse(tr) {
		v.ID = id
		fmt.Printf("Found Villager: %s\n", v)
		append(vl, v)
		var next bool
		next, tr = np.FindSibling(tr, "tag", "tr")
		if !next {
			break
		}
		time.Sleep(time.Second)
		id++
	}

}
