package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

func Bonjourmadamevibdy() (Message, error) {

	log.Printf("Loading Bonjourmadame-Vi-BDY...\n")
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get("http://dites.bonjourmadame.fr")
	if err != nil {
		return BONJOURMADAMEVIBDY_DEFAULT.format(), err
	}
	defer resp.Body.Close()

	log.Printf("Parsing Bonjourmadame-Vi-BDY...\n")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return BONJOURMADAMEVIBDY_DEFAULT.format(), err
	}

	defer log.Printf("Bonjourmadame-Vi-BDY processed!\n")
	return new(bonjourmadamevibdy).parse(doc).format(), nil

}
