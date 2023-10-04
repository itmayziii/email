package send

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"net/mail"
	"strings"
)

const pubSubType = "google.cloud.pubsub.topic.v1.messagePublished"

// PubSubPayload represents GCP pub/sub [MessagePublishedData format].
//
// [MessagePublishedData format]: https://googleapis.github.io/google-cloudevents/examples/binary/pubsub/MessagePublishedData-complex.json
type PubSubPayload struct {
	// Subscription name that this event is associated with.
	Subscription string        `json:"subscription"`
	Message      PubSubMessage `json:"message"`
}

// PubSubMessage is the [PubsubMessage format] when the message comes from Google's Pub/Sub.
//
// [PubsubMessage format]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Attributes  map[string]string `json:"attributes"`
	MessageId   string            `json:"messageId"`
	PublishTime string            `json:"publishTime"`
	// Data is automatically decoded from base64.
	Data []byte `json:"data"`
}

// EventData is this email packages specific event payload data needed to actually send an email.
type EventData struct {
	// Sender is who the email is from.
	Sender string `json:"sender"`
	// Subject is the email subject line.
	Subject string `json:"subject"`
	// Body is the email body and can be HTML. It will be parsed as a [Go HTML template] and bound to the variables
	// provided by Data. You should not pass both [EventData.Body] and [EventData.Template] at the same
	// time as they are both meant to represent the email body.
	//
	// [Go HTML template]: https://pkg.go.dev/html/template
	Body string `json:"body"`
	// To represents who the email should go to and can be provided as an array of strings or just a string
	To MessageTo `json:"to"`
	// Template is a path to the email template to use as the email [EventData.Body]. It will be parsed as
	// a [Go HTML template] and bound to the variables provided by Data. You should not pass both
	// [EventData.Template] and [EventData.Body] at the same time as they are both meant to represent
	// the email body.
	//
	// [Go HTML template]: https://pkg.go.dev/html/template
	Template string `json:"template"`
	// Data is an arbitrary map of variables to values that will be used in the [EventData.Template] or
	// [EventData.Body] using [Go HTML templates].
	//
	// [Go HTML templates]: https://pkg.go.dev/html/template
	Data map[string]interface{} `json:"data"`
}

// MessageTo represents who an email should be sent to.
type MessageTo []string

func (to *MessageTo) UnmarshalJSON(data []byte) error {
	rawTo := string(data)
	if strings.HasPrefix(rawTo, "[") {
		var emails []string
		dec := json.NewDecoder(bytes.NewReader(data))
		err := dec.Decode(&emails)
		if err != nil {
			return err
		}
		*to = emails
		return nil
	}

	*to = []string{rawTo}

	return nil
}

// extractEventData unmarshals the event payload into our expected [EventData] format.
func extractEventData(event cloudevents.Event) (EventData, error) {
	// Handling GCP pub/sub format is a convenience for package users. Not sure if we want to handle every custom
	// format, but we can cross that bridge if people start requesting them. It might not be a bad idea for this package
	// to play nicely with the major event producers.
	if event.Type() == pubSubType {
		var pubSubPayload PubSubPayload
		if err := event.DataAs(&pubSubPayload); err != nil {
			return EventData{}, err
		}
		var eventData EventData
		if err := json.Unmarshal(pubSubPayload.Message.Data, &eventData); err != nil {
			return EventData{}, err
		}

		return eventData, nil
	}

	var eventData EventData
	if err := event.DataAs(&eventData); err != nil {
		return EventData{}, err
	}
	return eventData, nil
}

// validateEventData ensures that [EventData] contains appropriate values such as having a valid sender, subject, etc...
func validateEventData(msgData EventData) error {
	if msgData.Sender == "" {
		return errors.New("missing \"sender\"")
	}
	if _, err := mail.ParseAddress(msgData.Sender); err != nil {
		return fmt.Errorf("invalid \"sender\" - %v", err)
	}

	if msgData.Subject == "" {
		return errors.New("missing \"subject\"")
	}

	if msgData.Body == "" && msgData.Template == "" {
		return errors.New("either \"body\" or \"template\" should be defined")
	}

	if len(msgData.To) == 0 {
		return errors.New("missing \"to\"")
	}
	for _, to := range msgData.To {
		if _, err := mail.ParseAddress(to); err != nil {
			return fmt.Errorf("invalid \"to\" - %v", err)
		}
	}

	return nil
}
