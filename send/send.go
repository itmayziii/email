/*
Package send exposes primitives to send emails by responding to [CloudEvents].

[CloudEvents]: https://cloudevents.io/
*/
package send

import (
	"context"
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
				log.Printf("failed to flush: %v", err)
			}
		}()

		eventData, err := extractEventData(event)
		if err != nil {
			app.errorLogger.Printf("failed to extract event data - %v", err)
			return err
		}
		err = validateEventData(eventData)
		if err != nil {
			app.errorLogger.Printf("invalid event data - %v", err)
			return err
		}

		emailBody, err := determineEmailBody(ctx, app, eventData)
		if err != nil {
			app.errorLogger.Printf("failed to determine email body %v", err)
			return err
		}

		domain, err := extractEmailDomain(eventData.Sender)
		if err != nil {
			app.errorLogger.Print(err)
			return err
		}
		sender, hasDomain := app.domainSenders[domain]
		if !hasDomain {
			err = fmt.Errorf(
				"domain: \"%s\" from \"sender\": \"%s\" does not match any registered domain to send emails from",
				domain,
				eventData.Sender,
			)
			app.errorLogger.Print(err)
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		id, err := sender.Send(ctx, Message{
			Sender:  eventData.Sender,
			Subject: eventData.Subject,
			Body:    emailBody,
			To:      eventData.To,
			Cc:      eventData.Cc,
			Bcc:     eventData.Bcc,
		})
		if err != nil {
			app.errorLogger.Printf("failed to send email: %v\n", err)
			return err
		}
		app.infoLogger.Printf(
			"email sent: id: %s, sender: %s, subject: %s, to: %s, cc: %s, bcc: %s\n",
			id,
			eventData.Sender,
			eventData.Subject,
			eventData.To,
			eventData.Cc,
			eventData.Bcc,
		)

		return nil
	}
}
