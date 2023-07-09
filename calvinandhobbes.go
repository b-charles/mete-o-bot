package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-charles/pigs/ioc"
	"github.com/benbjohnson/clock"
)

type CalvinAndHobbes struct {
	Clock      clock.Clock `inject:""`
	HttpClient HttpClient  `inject:""`
}

func (self *CalvinAndHobbes) Order() int {
	return 50
}

func (self *CalvinAndHobbes) Name() string {
	return "calvin-and-hobbes"
}

type calvinAndHobbesMessage struct {
	link string
}

func (self *calvinAndHobbesMessage) format() *Message {

	body := "CHARGEMENT EN COURS... ... ... Erreur 403: Calvin is sleeping."
	if self.link != "" {
		body = fmt.Sprintf("<img src=\"%s\"/>", self.link)
	}

	return &Message{
		Title: "Calvin And Hobbes",
		Body:  body,
	}

}

func (self *CalvinAndHobbes) Message() (*Message, error) {

	msg := &calvinAndHobbesMessage{}

	now := self.Clock.Now()
	url := fmt.Sprintf("https://www.gocomics.com/calvinandhobbes/%d/%02d/%02d",
		now.Year(), now.Month(), now.Day())

	resp, err := self.HttpClient.SimpleGet(url)
	if err != nil {
		return msg.format(), err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return msg.format(), err
	}

	msg.link = strings.TrimSpace(doc.Find("picture.item-comic-image img").First().AttrOr("src", ""))

	return msg.format(), nil

}

func init() {
	ioc.Put(&CalvinAndHobbes{}, func(Service) {})
}
