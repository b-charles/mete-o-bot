package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type citocin struct {
	citation  string
	signature string
}

var CITOCIN_DEFAULT = citocin{
	"Deux, dix et onze, comme les yeux, les doigts et les orteils.",
	"La famille Adams."}

func (self *citocin) parse(doc *goquery.Document) *citocin {

	self.citation = strings.TrimSpace(doc.Find("article > blockquote > p").First().Text())
	self.signature = strings.TrimSpace(doc.Find("article > blockquote > footer").First().Text())

	return self

}

const CITOCIN_FORMAT = "\"%s\" %s"

func (self *citocin) format() Message {
	msg := fmt.Sprintf(CITOCIN_FORMAT,
		self.citation, self.signature)
	return Message{Title: "Cit-O-CIN", Body: msg}
}

func Citocin() (Message, error) {

	resp, err := http.Get("http://www.kaakook.fr")
	if err != nil {
		return CITOCIN_DEFAULT.format(), err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return CITOCIN_DEFAULT.format(), err
	}

	return new(citocin).parse(doc).format(), nil

}
