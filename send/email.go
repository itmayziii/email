package send

import (
	"context"
	"fmt"
	"strings"
)

// Sender sends an email with the provided Message and returns the ID identifying the request. It should be noted
// that not all email providers provide any such ID, and therefore it may be an empty string.
type Sender interface {
	Send(ctx context.Context, m Message) (string, error)
}

// Message represents an email.
type Message struct {
	Sender  string
	Subject string
	Body    string
	To      []string
}

// NoopSender implements the [Sender] interface but doesn't actually send any emails which is helpful for testing
// purposes.
type NoopSender struct{}

func (ns NoopSender) Send(ctx context.Context, m Message) (string, error) {
	return "", nil
}

// determineEmailBody takes the [EventData.Body] or [EventData.Template] and executes them as
// Go HTML templates with variables being provided by [EventData.Data]. The result should be HTML appropriate to
// use as an email body.
//
// [Go HTML templates]: https://pkg.go.dev/html/template
func determineEmailBody(ctx context.Context, app *App, msgData EventData) (string, error) {
	unparsedBody := msgData.Body
	if unparsedBody == "" {
		templateBody, err := readTemplate(ctx, app, msgData.Template)
		if err != nil {
			return "", err
		}
		unparsedBody = templateBody
	}

	body, err := executeTemplate(unparsedBody, msgData.Data)
	if err != nil {
		return "", err
	}

	return body, nil
}

// extractEmailDomain returns the email domain and gives an error if no domain was found
func extractEmailDomain(email string) (string, error) {
	at := strings.LastIndex(email, "@")
	if at == -1 {
		return "", fmt.Errorf("provided email: \"%s\" has no domain", email)
	}

	return email[at+1:], nil
}
