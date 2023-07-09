package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/b-charles/pigs/ioc"
)

type meteoMessage struct {
	morningT   string
	morningW   string
	afternoonT string
	afternoonW string
	eveningT   string
	eveningW   string
	nightT     string
	nightW     string
	comment    string
}

func (self *meteoMessage) format() *Message {

	morningT := self.morningT
	if morningT == "" {
		morningT = "2"
	}
	morningW := self.morningW
	if morningW == "" {
		morningW = "Frais mais pas trop"
	}
	afternoonT := self.afternoonT
	if afternoonT == "" {
		afternoonT = "99"
	}
	afternoonW := self.afternoonW
	if afternoonW == "" {
		afternoonW = "Chaud bouillant"
	}
	eveningT := self.eveningT
	if eveningT == "" {
		eveningT = "21"
	}
	eveningW := self.eveningW
	if eveningW == "" {
		eveningW = "Tièdasse"
	}
	nightT := self.nightT
	if nightT == "" {
		nightT = "-9999999"
	}
	nightW := self.nightW
	if nightW == "" {
		nightW = "Sibérique"
	}
	comment := self.comment
	if comment == "" {
		comment = "Il fera si chaud aujourd'hui que nous conseillons de ne pas porter de vêtement."
	}

	msg := fmt.Sprintf("Météo Paris &#2947; "+
		"%s &#2947; Temps de la journée &#2947; "+
		"Matin: %s° %s &#2947; "+
		"Après-midi : %s° %s &#2947; "+
		"Soirée : %s° %s &#2947; "+
		"Nuit : %s° %s &#2947; "+
		"%s &#2947; Bonne journée !",
		today(),
		self.morningT, self.morningW,
		self.afternoonT, self.afternoonW,
		self.eveningT, self.eveningW,
		self.nightT, self.nightW,
		self.comment)

	return &Message{Title: "Météo Paris", Body: msg}

}

func extractTemp(doc *goquery.Document, idx int) string {
	path := fmt.Sprintf("div.forecasts > div.forecasts__item:nth-of-type(%d) > div.forecasts__item--temp", idx)
	return strings.TrimSpace(doc.Find(path).First().Text())
}

func extractWeather(doc *goquery.Document, idx int) string {
	path := fmt.Sprintf("div.forecasts > div.forecasts__item:nth-of-type(%d) > svg.forecasts__item--picto", idx)
	return strings.TrimSpace(doc.Find(path).First().AttrOr("data-tippy", "<missing>"))
}

type Meteo struct{}

func (self *Meteo) Order() int {
	return 10
}

func (self *Meteo) Name() string {
	return "meteopnm"
}

func (self *Meteo) Message() (*Message, error) {

	msg := &meteoMessage{}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get("https://www.meteo-paris.com")
	if err != nil {
		return msg.format(), err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return msg.format(), err
	}

	msg.morningT = extractTemp(doc, 1)
	msg.morningW = extractWeather(doc, 1)
	msg.afternoonT = extractTemp(doc, 2)
	msg.afternoonW = extractWeather(doc, 2)
	msg.eveningT = extractTemp(doc, 3)
	msg.eveningW = extractWeather(doc, 3)
	msg.nightT = extractTemp(doc, 4)
	msg.nightW = extractWeather(doc, 4)

	msg.comment = strings.TrimSpace(doc.Find("div.day-legend__text").First().Text())

	return msg.format(), nil

}

func init() {
	ioc.Put(&Meteo{}, func(Service) {})
}
