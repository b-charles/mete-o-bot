package main

import (
	"fmt"
	"time"
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

func today() string {

	now := time.Now()

	return fmt.Sprintf(
		"%s %d %s %d",
		WEEK[now.Weekday()],
		now.Day(),
		MONTHS[now.Month()],
		now.Year())

}
