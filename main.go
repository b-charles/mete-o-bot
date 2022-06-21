package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/gomail.v2"
)

// Error collection

type ErrorCollection struct {
	errors []error
}

func (self *ErrorCollection) AddError(err error) {
	self.errors = append(self.errors, err)
}

func (self *ErrorCollection) Add(format string, a ...any) {
	self.AddError(errors.New(fmt.Sprintf(format, a...)))
}

func (self *ErrorCollection) NilIfEmpty() error {
	if len(self.errors) == 0 {
		return nil
	} else {
		return self
	}
}

func (self *ErrorCollection) Error() string {

	var builder strings.Builder

	l := len(self.errors)

	fmt.Fprintf(&builder, "Collected errors (%d):\n", l)
	for i, e := range self.errors {
		fmt.Fprintf(&builder, "\tError %d/%d: %s\n", i+1, l, e.Error())
	}

	return builder.String()

}

// Message / Mail body

type Message struct {
	Title string
	Body  string
}

func (self Message) String() string {
	return fmt.Sprintf("%s\n---\n%s\n", self.Title, self.Body)
}

func CompleteMessage() (Message, error) {

	messages := []Message{}
	errors := new(ErrorCollection)

	for _, messager := range []func() (Message, error){Meteopnm, Citocin, Bonjourmadamevibdy} {
		msg, err := messager()
		messages = append(messages, msg)
		if err != nil {
			errors.AddError(err)
		}
	}

	var builder strings.Builder
	for _, msg := range messages {
		fmt.Fprintf(&builder, "%s <br>~~~~~~~~~<br><br> %s <br><br>", msg.Title, msg.Body)
	}

	return Message{
			Title: fmt.Sprintf("Mété-O-BOT %s", today()),
			Body:  builder.String()},
		errors.NilIfEmpty()

}

// Mail recipients

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

// Main handle

func handle() error {

	errors := new(ErrorCollection)

	sendMessage, err := CompleteMessage()
	if err != nil {
		errors.AddError(err)
	}

	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		errors.AddError(err)
		return errors
	}

	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	log.Printf("Dial SMTP server..\n.")
	dialer := gomail.NewDialer(host, port, user, password)
	sender, err := dialer.Dial()
	if err != nil {
		errors.AddError(err)
		return errors
	}

	from := os.Getenv("MAIL_FROM")

	log.Printf("Sending...\n")

	message := gomail.NewMessage()
	for _, recipient := range mailTo() {

		log.Printf("Sending to %s...\n", recipient.Name)

		message.SetHeader("From", from)
		message.SetAddressHeader("To", recipient.Address, recipient.Name)
		message.SetHeader("Subject", sendMessage.Title)
		message.SetBody("text/html", sendMessage.Body)

		if err := gomail.Send(sender, message); err != nil {
			errors.Add("Could not send email to %v (%v): %v", recipient.Name, recipient.Address, err)
		}

		message.Reset()

	}

	log.Printf("Sent!\n")

	return errors.NilIfEmpty()

}

func main() {
	lambda.Start(handle)
}
