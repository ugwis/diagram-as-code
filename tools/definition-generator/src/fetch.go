package main

import (
	//"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchURL() string {
	resp, err := http.Get("https://aws.amazon.com/architecture/icons/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var url string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(s.Text())
		if strings.TrimSpace(s.Text()) == "Download PPTx for Light Backgrounds" {
			url, _ = s.Attr("href")
		}
	})

	return url
}
