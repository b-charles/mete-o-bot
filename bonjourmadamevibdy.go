package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type bonjourmadamevibdy struct {
	link string
}

var BONJOURMADAMEVIBDY_DEFAULT = bonjourmadamevibdy{
	"CHARGEMENT EN COURS... ... ... Erreur 403: Forbidden. Dommage, elle était rigolote, celle là..."}

func (self *bonjourmadamevibdy) parse(doc *goquery.Document) *bonjourmadamevibdy {

	src, exist := doc.Find("div[class=\"post-content\"] img").First().Attr("src")
	if exist {
		self.link = fmt.Sprintf("<img src=\"%s\"/>", strings.TrimSpace(src))
	} else {
		self.link = BONJOURMADAMEVIBDY_DEFAULT.link
	}

	return self

}

func (self *bonjourmadamevibdy) format() Message {
	return Message{Title: "BonjourMadame-Vi-BDY", Body: self.link}
}

func Bonjourmadamevibdy() Message {

	resp, err := http.Get("http://dites.bonjourmadame.fr")
	if err != nil {
		log.Print(err)
		return BONJOURMADAMEVIBDY_DEFAULT.format()
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Print(err)
		return BONJOURMADAMEVIBDY_DEFAULT.format()
	}

	return (&bonjourmadamevibdy{}).parse(doc).format()

}
