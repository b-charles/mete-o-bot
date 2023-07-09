package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/log"
)

type Message struct {
	Title string
	Body  string
}

type Service interface {
	Order() int
	Name() string
	Message() (*Message, error)
}

type Services struct {
	Logger   log.Logger `inject:""`
	Today    *Today     `inject:""`
	Services []Service  `inject:""`
}

type serviceResult struct {
	service Service
	message *Message
}

func (self *Services) CompileMessages() *Message {

	n := len(self.Services)

	c := make(chan *serviceResult, n)

	for _, service := range self.Services {
		srv := service
		go func() {
			msg, err := srv.Message()
			if err != nil {
				self.Logger.Error().Set("service", srv.Name()).Set("error", err).Log()
			}
			c <- &serviceResult{srv, msg}
		}()
	}

	results := make([]*serviceResult, n)
	for i := 0; i < n; i++ {
		results = append(results, <-c)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].service.Order() < results[j].service.Order()
	})

	var builder strings.Builder
	for _, result := range results {
		msg := result.message
		fmt.Fprintf(&builder, "%s <br>~~~~~~~~~<br><br> %s <br><br>", msg.Title, msg.Body)
	}

	return &Message{
		Title: fmt.Sprintf("Mété-O-BOT %s", self.Today.Get()),
		Body:  builder.String(),
	}

}

func init() {
	ioc.Put(&Services{})
}
