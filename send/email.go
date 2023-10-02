package send

import (
	"context"
	"fmt"
	"strings"
)

// NoopSender implements the [Sender] interface but doesn't actually send any emails which is helpful for testing
// purposes.
type NoopSender struct{}

func (ns NoopSender) Send(ctx context.Context, m Message) (string, string, error) {
	return "", "", nil
}

// determineEmailBody takes the [EventMessageData.Body] or [EventMessageData.Template] and executes them as
// Go HTML templates with variables being provided by [EventMessageData.Data]. The result should be HTML appropriate to
// use as an email body.
//
// [Go HTML templates]: https://pkg.go.dev/html/template
func determineEmailBody(ctx context.Context, app *App, msgData EventMessageData) (string, error) {
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
