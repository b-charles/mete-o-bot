package main

import (
	"fmt"

	"github.com/b-charles/pigs/ioc"
	"github.com/benbjohnson/clock"
)

var WEEK = []string{
	"Dimanche",
	"Lundi",
	"Mardi",
	"Mercredi",
	"Jeudi",
	"Vendredi",
	"Samedi"}

var MONTHS = []string{"",
	"Janvier",
	"Février",
	"Mars",
	"Avril",
	"Mai",
	"Juin",
	"Juillet",
	"Août",
	"Septembre",
	"Octobre",
	"Novembre",
	"Décembre"}

type Today struct {
	Clock clock.Clock `inject:""`
}

func (self *Today) Get() string {

	now := self.Clock.Now()

	return fmt.Sprintf(
		"%s %d %s %d",
		WEEK[now.Weekday()],
		now.Day(),
		MONTHS[now.Month()],
		now.Year())

}

func init() {
	ioc.Put(&Today{})
}
