package main_test

import (
	_ "embed"
	"fmt"

	. "github.com/b-charles/mete-o-bot"
	"github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ServiceA struct{}

func (self *ServiceA) Order() int {
	return 10
}

func (self *ServiceA) Name() string {
	return "Service A"
}

func (self *ServiceA) Message() (*Message, error) {
	return &Message{
		Title: "Service A",
		Body:  "My beautiful service A.",
	}, nil
}

type ServiceB struct{}

func (self *ServiceB) Order() int {
	return 20
}

func (self *ServiceB) Name() string {
	return "Service B"
}

func (self *ServiceB) Message() (*Message, error) {
	return &Message{
		Title: "Service B",
		Body:  "My gorgeous service B.",
	}, nil
}

var _ = Describe("Services", func() {

	It("should generate a compiled message", func() {

		ioc.TestPut(&ServiceA{}, func(Service) {})
		ioc.TestPut(&ServiceB{}, func(Service) {})

		ioc.CallInjected(func(services *Services) {
			msg := services.CompileMessages()
			Expect(msg.Body).ShouldNot(BeEmpty())
			fmt.Printf(">> %s\n", msg.Body)
		})

	})

})
