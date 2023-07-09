package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-charles/pigs/ioc"
)

type BonjourMadame struct {
	HttpClient HttpClient `inject:""`
}

func (self *BonjourMadame) Order() int {
	return 30
}

func (self *BonjourMadame) Name() string {
	return "bonjour-madame"
}

type bonjourMadameMessage struct {
	link string
}

func (self *bonjourMadameMessage) format() *Message {

	body := "CHARGEMENT EN COURS... ... ... Erreur 403: Forbidden. Dommage, elle était rigolote, celle là..."
	if self.link != "" {
		body = fmt.Sprintf("<img src=\"%s\"/>", self.link)
	}

	return &Message{
		Title: "Bonjour Madame",
		Body:  body,
	}

}

func (self *BonjourMadame) Message() (*Message, error) {

	msg := &bonjourMadameMessage{}

	resp, err := self.HttpClient.SimpleGet("http://dites.bonjourmadame.fr")
	if err != nil {
		return msg.format(), err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return msg.format(), err
	}

	msg.link = strings.TrimSpace(doc.Find("div[class=\"post-content\"] img").First().AttrOr("src", ""))

	return msg.format(), nil

}

func init() {
	ioc.Put(&BonjourMadame{}, func(Service) {})
}
