package email_api

import (
	"cloud.google.com/go/logging"
	"context"
	"errors"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/mailgun/mailgun-go/v4"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"
	"log"
	"os"
	"time"
)

type app struct {
	mgTommyMayDev *mailgun.MailgunImpl
	loggingClient *logging.Client
	logger        *logging.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	fileStorage   *blob.Bucket
}

func init() {
	loggingClient, err := logging.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal("failed to create logging client", err)
	}
	logger := loggingClient.Logger("email-api", logging.RedirectAsJSON(os.Stdout))

	app := &app{
		mgTommyMayDev: mailgun.NewMailgun("mg.tommymay.dev", os.Getenv("MG_API_KEY_MG_TOMMYMAY_DEV")),
		loggingClient: loggingClient,
		logger:        logger,
		infoLogger:    logger.StandardLogger(logging.Info),
		errorLogger:   logger.StandardLogger(logging.Error),
	}
	functions.CloudEvent("SendEmail", sendEmail(app))
}

// Specifically not returning the error in bad data situations so that the events are not retried. These types
// of events will never succeed and should not be retried.
func sendEmail(app *app) func(context.Context, cloudevents.Event) error {
	return func(ctx context.Context, event cloudevents.Event) error {
		defer func() {
			if err := app.logger.Flush(); err != nil {
				log.Printf(fmt.Sprintf("failed to flush logger: %v", err))
			}
		}()

		_, msgData, err := extractEventData(app, event)
		if err != nil {
			app.errorLogger.Printf("failed to extract event data - %v", err)
			return nil
		}

		bucketName := os.Getenv("BUCKET")
		bucket, err := blob.OpenBucket(ctx, bucketName)
		if err != nil {
			app.errorLogger.Printf("failed to open bucket %s - %v", bucketName, err)
			return err
		}
		defer func() {
			if err := bucket.Close(); err != nil {
				app.errorLogger.Printf("failed to close bucket %s - %v", bucketName, err)
			}
		}()
		app.fileStorage = bucket

		emailBody, err := determineEmailBody(ctx, app, msgData)
		if err != nil {
			app.errorLogger.Printf("failed to determine email body %v", err)

			if errors.As(err, &ReadTemplateError{}) {
				// In this specific case we should return the error to trigger the event to fire again later.
				// Being unable to read the template could be a network issue or the developer simply has not
				// put the template in place yet.
				return err
			}

			return nil
		}

		message := app.mgTommyMayDev.NewMessage(msgData.Sender, msgData.Subject, "", msgData.To)
		message.SetHtml(emailBody)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		_, _, err = app.mgTommyMayDev.Send(ctx, message)
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

func determineEmailBody(ctx context.Context, app *app, msgData PubSubMessageData) (string, error) {
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
