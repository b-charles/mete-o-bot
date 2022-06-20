package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type meteopnm struct {
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

var METEOPNM_DEFAULT = meteopnm{
	"2", "Frais mais pas trop",
	"99", "Chaud bouillant",
	"21", "Tièdasse",
	"-9999999", "Sibérique",
	"Il fera si chaud aujourd'hui que nous conseillons de ne pas porter de vêtement."}

func (self *meteopnm) extractTemp(doc *goquery.Document, idx int, def string) string {
	path := fmt.Sprintf("div.forecasts > div.forecasts__item:nth-of-type(%d) > div.forecasts__item--temp", idx)
	return strings.TrimSpace(doc.Find(path).First().Text())
}

func (self *meteopnm) extractWeather(doc *goquery.Document, idx int, def string) string {
	path := fmt.Sprintf("div.forecasts > div.forecasts__item:nth-of-type(%d) > svg.forecasts__item--picto", idx)
	return strings.TrimSpace(doc.Find(path).First().AttrOr("data-tippy", def))
}

func (self *meteopnm) parse(doc *goquery.Document) *meteopnm {

	self.morningT = self.extractTemp(doc, 1, METEOPNM_DEFAULT.morningT)
	self.morningW = self.extractWeather(doc, 1, METEOPNM_DEFAULT.morningW)
	self.afternoonT = self.extractTemp(doc, 2, METEOPNM_DEFAULT.afternoonT)
	self.afternoonW = self.extractWeather(doc, 2, METEOPNM_DEFAULT.afternoonW)
	self.eveningT = self.extractTemp(doc, 3, METEOPNM_DEFAULT.eveningT)
	self.eveningW = self.extractWeather(doc, 3, METEOPNM_DEFAULT.eveningW)
	self.nightT = self.extractTemp(doc, 4, METEOPNM_DEFAULT.nightT)
	self.nightW = self.extractWeather(doc, 4, METEOPNM_DEFAULT.nightW)

	self.comment = strings.TrimSpace(doc.Find("div.day-legend__text").First().Text())

	return self

}

const METEOPNM_FORMAT = "Météo Paris &#2947; " +
	"%s &#2947; Temps de la journée &#2947; " +
	"Matin: %s° %s &#2947; " +
	"Après-midi : %s° %s &#2947; " +
	"Soirée : %s° %s &#2947; " +
	"Nuit : %s° %s &#2947; " +
	"%s &#2947; Bonne journée !"

func (self *meteopnm) format() Message {
	msg := fmt.Sprintf(METEOPNM_FORMAT,
		today(),
		self.morningT, self.morningW,
		self.afternoonT, self.afternoonW,
		self.eveningT, self.eveningW,
		self.nightT, self.nightW,
		self.comment)
	return Message{Title: "Mété-O-PNM", Body: msg}
}

func Meteopnm() (Message, error) {

	resp, err := http.Get("https://www.meteo-paris.com")
	if err != nil {
		return METEOPNM_DEFAULT.format(), err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return METEOPNM_DEFAULT.format(), err
	}

	return new(meteopnm).parse(doc).format(), nil

}
