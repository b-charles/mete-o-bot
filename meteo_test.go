package main_test

import (
	_ "embed"
	"fmt"
	"strings"
	"time"

	. "github.com/b-charles/mete-o-bot"
	"github.com/b-charles/pigs/ioc"
	"github.com/benbjohnson/clock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakeClient struct{}

func (self *FakeClient) SimpleGet(url string) (string, error) {
	resp := <-self.Get(url, map[string]string{})
	return resp.Body, resp.Err
}

//go:embed meteo_test_ephemeride.json
var ephemerideJson string

//go:embed meteo_test_forecast.json
var forecastJson string

func (self *FakeClient) Get(url string, headers map[string]string) chan *HttpResponse {

	c := make(chan *HttpResponse, 1)

	if strings.Contains(url, "ephemeride") {
		c <- &HttpResponse{
			Err:  nil,
			Body: ephemerideJson,
		}
	} else if strings.Contains(url, "forecast") {
		c <- &HttpResponse{
			Err:  nil,
			Body: forecastJson,
		}
	} else {
		c <- &HttpResponse{
			Err:  fmt.Errorf("unknown url"),
			Body: "",
		}
	}

	close(c)

	return c

}

var _ = Describe("Meteo", func() {

	BeforeEach(func() {

		t, _ := time.Parse("2006-01-02", "2023-07-09")
		mockedClock := clock.NewMock()
		mockedClock.Set(t)
		ioc.TestPut(mockedClock, func(clock.Clock) {})

		ioc.TestPut(&FakeClient{}, func(HttpClient) {})

	})

	It("should process a response", func() {

		ioc.CallInjected(func(service *Meteo) {
			msg, err := service.Message()
			Expect(err).ShouldNot(HaveOccurred())
			fmt.Printf(">%s> %s\n", service.Name(), msg.Body)
		})

	})

})
