package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/b-charles/pigs/ioc"
)

// Main handle

func main() {

	ioc.CallInjected(func(services *Services, mailSender *MailSender) {

		lambda.Start(func() {
			message := services.CompileMessages()
			mailSender.send(message)
		})

	})

}
