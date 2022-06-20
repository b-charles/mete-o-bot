package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type Message struct {
	Title string
	Body  string
}

func bodyMessage() string {

	var messages = []Message{
		Meteopnm(),
		Citocin(),
		Bonjourmadamevibdy()}

	var builder strings.Builder
	for _, msg := range messages {
		fmt.Fprintf(&builder, "%s <br>~~~~~~~~~<br><br> %s <br><br>", msg.Title, msg.Body)
	}
	return builder.String()

}

type Recipient struct {
	Name    string
	Address string
}

func mailTo() []Recipient {

	var r []Recipient
	if err := json.Unmarshal([]byte(os.Getenv("MAIL_TO")), &r); err != nil {
		panic(err)
	}

	return r

}

func main() {

	subject := fmt.Sprintf("Mété-O-BOT %s", today())
	body := bodyMessage()

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		fmt.Println(subject)
		fmt.Println("---")
		fmt.Println(body)
		return
	}

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}

	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	dialer := gomail.NewDialer(host, port, user, password)
	sender, err := dialer.Dial()
	if err != nil {
		panic(err)
	}

	from := os.Getenv("MAIL_FROM")

	message := gomail.NewMessage()

	for _, recipient := range mailTo() {

		message.SetHeader("From", from)
		message.SetAddressHeader("To", recipient.Address, recipient.Name)
		message.SetHeader("Subject", subject)
		message.SetBody("text/html", body)

		if err := gomail.Send(sender, message); err != nil {
			log.Printf("Could not send email to %v (%v): %v", recipient.Name, recipient.Address, err)
		}

		message.Reset()

	}

}
