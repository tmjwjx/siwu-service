package test

import (
	"fmt"
	"forum/observer"
	"forum/pkg/globals"
	"testing"
)

func TestPushMessage(t *testing.T) {

	userId := "1"
	eventName := "hello"

	notifyChan := make(chan string)
	globals.SubscriberChannels[userId] = notifyChan
	go func() {
		for {
			select {
			case message := <-notifyChan:
				fmt.Println("message ;;;;;;;;", message)
			}
		}
	}()

	observer.PushMessage(userId, eventName)
}