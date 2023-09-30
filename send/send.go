/*
Package send exposes primitives to send emails by responding to [CloudEvents].

[CloudEvents]: https://cloudevents.io/
*/
package send

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"log"
	"time"
)

// EmailEvent creates a function to send an email by responding to a [CloudEvent].
//
// [CloudEvent]: https://cloudevents.io/
func EmailEvent(app *App) func(context.Context, cloudevents.Event) error {
	return func(ctx context.Context, event cloudevents.Event) error {
		defer func() {
			if err := app.flusher.Flush(); err != nil {
				// Not appropriate to rely on the info or error logger here as that is probably the thing that is being
				// flushed and errored.
				log.Printf(fmt.Sprintf("failed to flush: %v", err))
			}
		}()

		_, msgData, err := extractEventData(event)
		if err != nil {
			app.errorLogger.Printf("failed to extract event data - %v", err)
			// Being unable to unmarshal the event data likely means this message will never succeed. We log the error
			// and don't return it to prevent the event from replaying.
			return nil
		}
		err = validateMessageData(msgData)
		if err != nil {
			// Invalid data likely means this message will never succeed. We log the error and don't return it to
			// prevent the event from replaying.
			app.errorLogger.Printf("invalid event data - %v", err)
			return nil
		}

		emailBody, err := determineEmailBody(ctx, app, msgData)
		if err != nil {
			app.errorLogger.Printf("failed to determine email body %v", err)
			if errors.As(err, &ReadTemplateError{}) {
				// In this specific case we should return the error to trigger the event to replay.
				// Being unable to read the template could be a network issue or the developer simply has not
				// put the template in place yet.
				return err
			}

			// If the error is not a ReadTemplateError then it is probably related the email data not working
			// with the body variables. We do not return the error to prevent the event from replaying.
			return nil
		}

		domain := ""
		sender, hasDomain := app.domainSenders[domain]
		if !hasDomain {
			err = errors.New(fmt.Sprintf(
				"domain: \"%s\" from \"sender\": \"%s\" does not match any registered domain to send emails from",
				domain,
				msgData.Sender,
			))
			app.errorLogger.Printf("%v", err)
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		_, _, err = sender.Send(ctx, Message{
			Sender:  msgData.Sender,
			Subject: msgData.Subject,
			Body:    emailBody,
			To:      msgData.To,
		})
		if err != nil {
			app.errorLogger.Printf("failed to send email: %v\n", err)
			return err
		}
		app.infoLogger.Printf(
			"email sent: sender: %s, subject: %s, to: %s\n",
			msgData.Sender,
			msgData.Subject,
			msgData.To,
		)

		return nil
	}
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
