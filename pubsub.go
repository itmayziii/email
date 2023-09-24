package email_api

import (
	"encoding/json"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// PubSubData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type PubSubData struct {
	Subscription string            `json:"subscription"`
	Message      PubSubMessage     `json:"message"`
	Attributes   map[string]string `json:"attributes"`
	MessageId    string            `json:"messageId"`
	PublishTime  string            `json:"publishTime"`
	OrderingKey  string            `json:"orderingKey"`
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type PubSubMessageData struct {
	Sender   string                 `json:"sender"`
	Subject  string                 `json:"subject"`
	Body     string                 `json:"body"`
	To       string                 `json:"to"`
	Template string                 `json:"template"`
	Data     map[string]interface{} `json:"data"`
}

func extractEventData(app *app, event cloudevents.Event) (PubSubData, PubSubMessageData, error) {
	var msg PubSubData
	if err := event.DataAs(&msg); err != nil {
		return PubSubData{}, PubSubMessageData{}, err
	}
	var msgData PubSubMessageData
	if err := json.Unmarshal(msg.Message.Data, &msgData); err != nil {
		return PubSubData{}, PubSubMessageData{}, err
	}
	if msgData.Sender == "" {
		return PubSubData{}, PubSubMessageData{}, errors.New("missing \"sender\"")
	}
	if msgData.Subject == "" {
		return PubSubData{}, PubSubMessageData{}, errors.New("missing \"subject\"")
	}
	if msgData.Body == "" && msgData.Template == "" {
		return PubSubData{}, PubSubMessageData{}, errors.New("either \"body\" or \"template\" should be defined")
	}
	if msgData.To == "" {
		return PubSubData{}, PubSubMessageData{}, errors.New("missing \"to\"")
	}

	return msg, msgData, nil
}
