package main

import (
	"fmt"
	"github.com/zackshank/nookscraper/parser"
	"golang.org/x/net/html"
	//"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Villager struct {
	Name         string
	JapaneseName string
	Gender       string
	Species      string
	Personality  string
	Games        []int
	Birthday     time.Time
}

func (v *Villager) String() string {
	return fmt.Sprintf("%s (%s), %s, %s, %s, %v", v.Name, v.JapaneseName, v.Species, v.Gender, v.Personality, v.Games)
}

type VillagerParser struct{}

func (vp *VillagerParser) Parse(tr *html.Node) (bool, *Villager) {
	np := parser.NodeParser{}
	ch := make(chan bool)
	v := Villager{}
	// Start Location
	found, td := np.Find(tr, "tag", "th")
	if !found {
		fmt.Println("Could not parse Villager")
		return false, nil
	}

	// Get name and start request to get the rest of the information
	found, td = np.FindSibling(td, "tag", "th")
	if !found {
		fmt.Println("Could not find villager name")
		return false, nil
	}

	// Get url to next location
	ok, anode := np.Find(td, "tag", "a")
	if ok {
		_, attr := np.GetAttribute(anode, "href")
		url := fmt.Sprint(SiteRoot, attr.Val)
		go vp.parseAdditionalInformation(url, &v, ch)
	}

	found, v.Name = vp.parseName(td)

	if !found {
		fmt.Println("Could not find villager name")
		return false, nil
	}

	// Get JapaneseName
	found, td = np.FindSibling(td, "tag", "td")
	if !found {
		fmt.Println("Could not find villager japanese name")
		return false, nil
	}

	found, v.JapaneseName = vp.parseJapaneseName(td)

	if !found {
		fmt.Println("Could not find villager japanese name")
		return false, nil
	}

	// Get Species
	found, td = np.FindSibling(td, "tag", "td")
	if !found {
		fmt.Println("Could not find villager species")
		return false, nil
	}

	found, v.Species = vp.parseSpecies(td)

	if !found {
		fmt.Println("Could not find villager species")
		return false, nil
	}

	// Get Gender
	found, td = np.FindSibling(td, "tag", "td")
	if !found {
		fmt.Println("Could not find villager gender")
		return false, nil
	}

	found, v.Gender = vp.parseGender(td)

	if !found {
		fmt.Println("Could not find villager gender")
		return false, nil
	}

	// Get Personality
	found, td = np.FindSibling(td, "tag", "td")
	if !found {
		fmt.Println("Could not find villager personality")
		return false, nil
	}
	found, v.Personality = vp.parsePersonality(td)

	if !found {
		fmt.Println("Could not find villager personality")
		return false, nil
	}

	// Get Games
	found, td = np.FindSibling(td, "tag", "td")
	if !found {
		fmt.Println("Could not find villager games")
		return false, nil
	}
	found, v.Games = vp.parseGames(td)

	if !found {
		fmt.Println("Could not find villager games")
		return false, nil
	}

	result := <-ch

	fmt.Println(result)

	return true, &v
}

func (vp *VillagerParser) parseName(td *html.Node) (bool, string) {
	np := parser.NodeParser{}

	found, namenode := np.Find(td, "tag", "a")
	if !found {
		return false, ""
	}

	textnode := namenode.FirstChild
	if textnode.Type != html.TextNode {
		return false, ""
	}

	return true, string(textnode.Data)

}

func (vp *VillagerParser) parseJapaneseName(td *html.Node) (bool, string) {
	np := parser.NodeParser{}

	found, namenode := np.Find(td, "tag", "b")
	if !found {
		return false, ""
	}

	textnode := namenode.FirstChild
	if textnode.Type != html.TextNode {
		return false, ""
	}

	return true, string(textnode.Data)
}

func (vp *VillagerParser) parseSpecies(td *html.Node) (bool, string) {
	np := parser.NodeParser{}

	found, anode := np.Find(td, "tag", "a")
	if !found {
		return false, ""
	}

	textnode := anode.FirstChild
	if textnode.Type != html.TextNode {
		return false, ""
	}

	return true, string(textnode.Data)
}

func (vp *VillagerParser) parseGender(td *html.Node) (bool, string) {
	textnode := td.FirstChild
	if textnode.Type != html.TextNode {
		return false, ""
	}

	return true, strings.TrimSpace(string(textnode.Data))
}

func (vp *VillagerParser) parsePersonality(td *html.Node) (bool, string) {
	np := parser.NodeParser{}

	found, anode := np.Find(td, "tag", "a")
	if !found {
		return false, ""
	}

	textnode := anode.FirstChild
	if textnode.Type != html.TextNode {
		return false, ""
	}

	return true, string(textnode.Data)
}

func (vp *VillagerParser) parseGames(td *html.Node) (bool, []int) {
	np := parser.NodeParser{}

	var g []int

	for i := 1; i <= 8; i++ {
		found, _ := np.Find(td, "tag", "a")
		if found {
			g = append(g, i)
		}

		if i < 8 {
			var ok bool
			ok, td = np.FindSibling(td, "tag", "td")
			if !ok {
				return false, g
			}
		}
	}

	return true, g
}

func (vp *VillagerParser) parseAdditionalInformation(url string, v *Villager, c chan bool) {
	np := parser.NodeParser{}
	// Get page
	resp, err := http.Get(url)
	if err != nil {
		c <- false
		return
	}

	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		c <- false
		return
	}

	// Find Table with Addt. Info
	_, mw := np.Find(doc, "id", "mw-content-text")
	_, table := np.Find(mw, "tag", "table")
	_, table = np.FindSibling(table, "tag", "table")

	// Attempt to retrieve birthday
	vp.parseBirthday(table)

	c <- false

}

func (vp *VillagerParser) parseBirthday(root *html.Node) (bool, string) {
	np := parser.NodeParser{}

	found, th := np.Find(root, "html", "Birthday")
	if !found {
		return false, ""
	}

	_, td := np.FindSibling(th, "tag", "td")
	_, anode := np.Find(td, "tag", "a")
	textelement := anode.FirstChild
	fmt.Println("Found Birthday: ", string(textelement.Data))
	return true, ""
}
