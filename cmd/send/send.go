package main

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
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
	//
	//messageData, err := json.Marshal(send.PubSubMessageData{
	//	Sender:   "no-reply@tommymay.dev",
	//	Subject:  "Hello World",
	//	To:       "tommymay37@gmail.com",
	//	Template: "contact.gohtml",
	//	Data: map[string]interface{}{
	//		"name":    "Johnny",
	//		"subject": "hello subject",
	//		"email":   "johnny@example.com",
	//		"message": "Hello message",
	//	},
	//})
	//if err != nil {
	//	log.Fatalf("failed to marshal pub sub message data, %v", messageData)
	//}
	//pubSubData := send.PubSubData{
	//	Subscription: "projects/test-project/subscriptions/my-subscription",
	//	Message: send.PubSubMessage{
	//		Data: messageData,
	//	},
	//	Attributes: map[string]string{
	//		"app": "itmayziii-api",
	//	},
	//	MessageId:   "123kjsadjk",
	//	PublishTime: "asdf asf",
	//	OrderingKey: "askfj",
	//}
	//data, err := json.Marshal(pubSubData)
	//if err != nil {
	//	log.Fatalf("failed to marshal pub sub data, %v", pubSubData)
	//}

	data2 := "{\"message\":{\"attributes\":{\"hello\":\"world\"},\"data\":\"eyJoZWxsbyI6IndvcmxkIiwic2VuZGVyIjoidG9tbXltYXkifQo=\",\"messageId\":\"8300998372576141\",\"message_id\":\"8300998372576141\",\"publishTime\":\"2023-09-27T18:21:16.418Z\",\"publish_time\":\"2023-09-27T18:21:16.418Z\"},\"subscription\":\"projects/itmayziii/subscriptions/eventarc-us-central1-email-func-901633-sub-154\"}"
	event.SetData(cloudevents.ApplicationJSON, []byte(data2))
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
