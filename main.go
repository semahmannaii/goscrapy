package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

type Item struct {
	Title   string `json:"title"`
	Details string `json:"details"`
	ImgUrl  string `json:"imgUrl"`
}

func main() {
	formatJson := regexp.MustCompile(`\s+`)

	c := colly.NewCollector(colly.AllowedDomains("www.domain.com"))

	var items []Item

	c.OnHTML("div.className", func(e *colly.HTMLElement) {
		item := Item{
			Title:   formatJson.ReplaceAllString(strings.TrimSpace(e.ChildText("h2.className")), " "),
			Details: formatJson.ReplaceAllString(strings.TrimSpace(e.ChildText("p.className")), " "),
			ImgUrl:  formatJson.ReplaceAllString(strings.TrimSpace(e.ChildAttr("img", "src")), " "),
		}

		items = append(items, item)
	})

	c.OnHTML("a.className", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	c.Visit("https://www.domain.com/full-link")

	result, err := json.Marshal(items)

	if err != nil {
		log.Fatal(err)
		return
	}

	os.WriteFile("file.json", result, 0644)
}
