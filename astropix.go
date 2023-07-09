package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-charles/pigs/ioc"
)

type Astropix struct {
	HttpClient HttpClient `inject:""`
}

func (self *Astropix) Order() int {
	return 40
}

func (self *Astropix) Name() string {
	return "astropix"
}

type astropixMessage struct {
	img         string
	explanation string
}

func (self *astropixMessage) format() *Message {

	img := "CHARGEMENT EN COURS... ... ... Erreur 403: Forbidden by Gods."
	if self.img != "" {
		img = fmt.Sprintf("<img src=\"%s\"/>", self.img)
	}

	explanation := "Gods was in this picture and do not like it."
	if self.explanation != "" {
		explanation = self.explanation
	}

	msg := fmt.Sprintf("%s <br> %s", img, explanation)

	return &Message{Title: "Astropix", Body: msg}

}

var (
	ASTROPIX_URL, _ = url.Parse("https://apod.nasa.gov/apod/astropix.html")
	HREF_PATTERN    = regexp.MustCompile(`href="[^"]*"`)
	SPACES_PATTERN  = regexp.MustCompile(`\s+`)
)

func (self *Astropix) Message() (*Message, error) {

	msg := &astropixMessage{}

	resp, err := self.HttpClient.SimpleGet(ASTROPIX_URL.String())
	if err != nil {
		return msg.format(), err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return msg.format(), err
	}

	img := strings.TrimSpace(doc.Find("body img").First().AttrOr("src", ""))
	if img != "" {
		link, _ := url.Parse(img)
		msg.img = ASTROPIX_URL.ResolveReference(link).String()
	}

	if explanation, err := doc.Find("body > center").Next().Next().Html(); err != nil {
		return msg.format(), err
	} else {

		cleaned := HREF_PATTERN.ReplaceAllStringFunc(explanation, func(m string) string {

			link, _ := strings.CutPrefix(m, `href="`)
			link, _ = strings.CutSuffix(link, `"`)
			if addr, err := url.Parse(link); err != nil {
				return m
			} else {
				return fmt.Sprintf(`href="%s"`, ASTROPIX_URL.ResolveReference(addr))
			}

		})
		cleaned = SPACES_PATTERN.ReplaceAllLiteralString(cleaned, " ")

		msg.explanation = strings.TrimSpace(cleaned)
	}

	return msg.format(), nil

}

func init() {
	ioc.Put(&Astropix{}, func(Service) {})
}
