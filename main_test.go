package main

import (
	"fmt"
	"testing"
)

func TestCompleteMessage(t *testing.T) {

	message, err := CompleteMessage()

	fmt.Println(message.String())

	if err != nil {
		fmt.Printf("%v (%T): %t", err, err, err == nil)
		t.Fatal(err)
	}

}
