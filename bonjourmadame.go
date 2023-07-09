package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-charles/pigs/ioc"
)

type BonjourMadame struct{}

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

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("http://dites.bonjourmadame.fr")
	if err != nil {
		return msg.format(), err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return msg.format(), err
	}

	src, exist := doc.Find("div[class=\"post-content\"] img").First().Attr("src")
	if exist {
		msg.link = strings.TrimSpace(src)
	}

	return msg.format(), nil

}

func init() {
	ioc.Put(&BonjourMadame{}, func(Service) {})
}
