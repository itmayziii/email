package main

import (
	"context"
	"encoding/json"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	emailapi "github.com/itmayziii/email-api"
	"log"
)

func main() {
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	// Create an Event.
	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")

	messageData, err := json.Marshal(emailapi.PubSubMessageData{
		Sender:   "no-reply@tommymay.dev",
		Subject:  "Hello World",
		To:       "tommymay37@gmail.com",
		Template: "contact.gohtml",
		Data: map[string]interface{}{
			"name":    "Johnny",
			"subject": "hello subject",
			"email":   "johnny@example.com",
			"message": "Hello message",
		},
	})
	if err != nil {
		log.Fatalf("failed to marshal pub sub message data, %v", messageData)
	}
	pubSubData := emailapi.PubSubData{
		Subscription: "projects/test-project/subscriptions/my-subscription",
		Message: emailapi.PubSubMessage{
			Data: messageData,
		},
		Attributes: map[string]string{
			"app": "itmayziii-api",
		},
		MessageId:   "123kjsadjk",
		PublishTime: "asdf asf",
		OrderingKey: "askfj",
	}
	data, err := json.Marshal(pubSubData)
	if err != nil {
		log.Fatalf("failed to marshal pub sub data, %v", pubSubData)
	}

	event.SetData(cloudevents.ApplicationJSON, data)
	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), "http://localhost:8080/")

	// Send that Event.
	if result := c.Send(ctx, event); cloudevents.IsUndelivered(result) {
		log.Fatalf("failed to send, %v", result)
	} else {
		log.Printf("sent: %v", event)
		log.Printf("result: %v", result)
	}
}
