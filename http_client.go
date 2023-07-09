package main

import (
	"io"
	"net/http"
	"time"

	"github.com/b-charles/pigs/config"
	"github.com/b-charles/pigs/ioc"
	"github.com/b-charles/pigs/smartconfig"
)

type HttpResponse struct {
	Err  error
	Body string
}

func httpResponseError(err error) *HttpResponse {
	return &HttpResponse{
		Err:  err,
		Body: "",
	}
}

type HttpClient interface {
	SimpleGet(url string) (string, error)
	Get(url string, headers map[string]string) chan *HttpResponse
}

type RealHttpClientConfig struct {
	Timeout int
}

type RealHttpClient struct {
	client *http.Client
}

func (self *RealHttpClient) SimpleGet(url string) (string, error) {
	resp := <-self.Get(url, map[string]string{})
	return resp.Body, resp.Err
}

func (self *RealHttpClient) Get(url string, headers map[string]string) chan *HttpResponse {

	c := make(chan *HttpResponse, 1)

	if req, err := http.NewRequest("GET", url, nil); err != nil {
		c <- httpResponseError(err)
		close(c)
	} else {

		for k, v := range headers {
			req.Header.Add(k, v)
		}

		go func() {

			resp, err := self.client.Do(req)
			if err != nil {
				c <- httpResponseError(err)
				close(c)
				return
			}
			defer resp.Body.Close()

			if body, err := io.ReadAll(resp.Body); err != nil {
				c <- httpResponseError(err)
			} else {
				c <- &HttpResponse{
					Err:  nil,
					Body: string(body),
				}
			}
			close(c)

		}()

	}

	return c

}

func init() {

	config.Set("http.timeout", "3")
	smartconfig.Configure("http", &RealHttpClientConfig{})

	ioc.PutFactory(func(config *RealHttpClientConfig) *RealHttpClient {
		return &RealHttpClient{
			client: &http.Client{
				Timeout: time.Duration(config.Timeout) * time.Second,
			},
		}
	}, func(HttpClient) {})

}
