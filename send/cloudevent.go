package send

import (
	"bytes"
	"encoding/json"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"strings"
)

// EventData is the highest level of the payload an event should contain. It models Google's Pub/Sub
// [MessagePublishedData format].
//
// [MessagePublishedData format]: https://googleapis.github.io/google-cloudevents/examples/binary/pubsub/MessagePublishedData-complex.json
type EventData struct {
	// Subscription name that this event is associated with.
	Subscription string       `json:"subscription"`
	Message      EventMessage `json:"message"`
}

// EventMessage is the event payload message which contains data following Google's Pub/Sub [PubsubMessage format].
//
// [PubsubMessage format]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type EventMessage struct {
	Attributes  map[string]string `json:"attributes"`
	MessageId   string            `json:"messageId"`
	PublishTime string            `json:"publishTime"`
	// Data is automatically decoded from base64.
	Data []byte `json:"data"`
}

// EventMessageData is the event payload data needed to actually send an email.
type EventMessageData struct {
	// Sender is who the email is from.
	Sender string `json:"sender"`
	// Subject is the email subject line.
	Subject string `json:"subject"`
	// Body is the email body and can be HTML. It will be parsed as a [Go HTML template] and bound to the variables
	// provided by Data. You should not pass both [EventMessageData.Body] and [EventMessageData.Template] at the same
	// time as they are both meant to represent the email body.
	//
	// [Go HTML template]: https://pkg.go.dev/html/template
	Body string `json:"body"`
	// To represents who the email should go to and can be provided as an array of strings or just a string
	To MessageTo `json:"to"`
	// Template is a path to the email template to use as the email [EventMessageData.Body]. It will be parsed as
	// a [Go HTML template] and bound to the variables provided by Data. You should not pass both
	// [EventMessageData.Template] and [EventMessageData.Body] at the same time as they are both meant to represent
	// the email body.
	//
	// [Go HTML template]: https://pkg.go.dev/html/template
	Template string `json:"template"`
	// Data is an arbitrary map of variables to values that will be used in the [EventMessageData.Template] or
	// [EventMessageData.Body] using [Go HTML templates].
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

// extractEventData unmarshals the event payload into our expected [EventData] and [EventMessageData] format.
func extractEventData(event cloudevents.Event) (EventData, EventMessageData, error) {
	var eventData EventData
	if err := event.DataAs(&eventData); err != nil {
		return EventData{}, EventMessageData{}, err
	}
	var msgData EventMessageData
	if err := json.Unmarshal(eventData.Message.Data, &msgData); err != nil {
		return EventData{}, EventMessageData{}, err
	}

	return eventData, msgData, nil
}

// validateMessageData ensures that [EventMessageData] contains appropriate values such as having a valid sender, subject,
// etc...
func validateMessageData(msgData EventMessageData) error {
	if msgData.Sender == "" {
		return errors.New("missing \"sender\"")
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

	return nil
}
