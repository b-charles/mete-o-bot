package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-charles/pigs/ioc"
)

type Kaakook struct {
	HttpClient HttpClient `inject:""`
}

func (self *Kaakook) Order() int {
	return 20
}

func (self *Kaakook) Name() string {
	return "kaakook"
}

type kaakookMessage struct {
	citation  string
	signature string
}

func (self *kaakookMessage) format() *Message {

	citation := self.citation
	if citation == "" {
		citation = "Deux, dix et onze, comme les yeux, les doigts et les orteils."
	}

	signature := self.signature
	if signature == "" {
		signature = "La famille Adams."
	}

	msg := fmt.Sprintf("\"%s\" %s", citation, signature)

	return &Message{Title: "Kaakook", Body: msg}

}

func (self *Kaakook) Message() (*Message, error) {

	msg := &kaakookMessage{}

	resp, err := self.HttpClient.SimpleGet("http://www.kaakook.fr")
	if err != nil {
		return msg.format(), err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return msg.format(), err
	}

	msg.citation = strings.TrimSpace(doc.Find("article > blockquote > p").First().Text())
	msg.signature = strings.TrimSpace(doc.Find("article > blockquote > footer").First().Text())

	return msg.format(), nil

}

func init() {
	ioc.Put(&Kaakook{}, func(Service) {})
}
