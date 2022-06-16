package main

import (
	"fmt"
	"strings"
)

type Message struct {
	Title string
	Body  string
}

func formatMessages(messages []Message) string {
	var builder strings.Builder
	for _, msg := range messages {
		fmt.Fprintf(&builder, "%s <br>~~~~~~~~~<br><br> %s <br><br>", msg.Title, msg.Body)
	}
	return builder.String()
}

var messages = []Message{
	Meteopnm(),
	Citocin(),
	Bonjourmadamevibdy()}

func main() {
	fmt.Printf(formatMessages(messages))
}
