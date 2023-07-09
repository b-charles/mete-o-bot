package main_test

import (
	_ "embed"
	"fmt"

	. "github.com/b-charles/mete-o-bot"
	"github.com/b-charles/pigs/ioc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Astropix", func() {

	BeforeEach(func() {
		ioc.TestPut("test flag")
	})

	It("should generate a message", func() {

		ioc.CallInjected(func(service *Astropix) {
			msg, err := service.Message()
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Printf(">%s> %s\n", service.Name(), msg.Body)
		})

	})

})
