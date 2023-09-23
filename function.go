package email_api

import (
	"cloud.google.com/go/logging"
	"context"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"os"
	"time"
)

type App struct {
	mgTommyMayDev *mailgun.MailgunImpl
	loggingClient *logging.Client
	logger        *logging.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
}

func init() {
	loggingClient, err := logging.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal("failed to create logging client", err)
	}
	logger := loggingClient.Logger("email-api", logging.RedirectAsJSON(os.Stdout))

	app := &App{
		mgTommyMayDev: mailgun.NewMailgun("mg.tommymay.dev", os.Getenv("MG_API_KEY_MG_TOMMYMAY_DEV")),
		loggingClient: loggingClient,
		logger:        logger,
		infoLogger:    logger.StandardLogger(logging.Info),
		errorLogger:   logger.StandardLogger(logging.Error),
	}
	functions.CloudEvent("SendEmail", sendEmail(app))
}

// PubSubData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type PubSubData struct {
	Subscription string                 `json:"subscription"`
	Message      PubSubMessage          `json:"message"`
	Attributes   map[string]interface{} `json:"attributes"`
	MessageId    string                 `json:"messageId"`
	PublishTime  string                 `json:"publishTime"`
	OrderingKey  string                 `json:"orderingKey"`
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type PubSubMessageData struct {
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

func sendEmail(app *App) func(context.Context, cloudevents.Event) error {
	return func(ctx context.Context, event cloudevents.Event) error {
		defer func() {
			err := app.logger.Flush()
			if err != nil {
				log.Printf(fmt.Sprintf("failed to flush logger: %v", err))
			}
		}()

		// Specifically NOT returning the error in bad data situations so that the events are not retried. These types
		// of events will never succeed and should not be retried.
		var msg PubSubData
		if err := event.DataAs(&msg); err != nil {
			app.errorLogger.Printf("failed to parse event with event.DataAs: %v", err)
			return nil
		}
		var messageData PubSubMessageData
		if err := json.Unmarshal(msg.Message.Data, &messageData); err != nil {
			app.errorLogger.Printf("failed to parse message data: %v", err)
			return nil
		}
		if messageData.Sender == "" {
			app.errorLogger.Println("failed to send email, missing \"sender\"")
			return nil
		}
		if messageData.Subject == "" {
			app.errorLogger.Println("failed to send email, missing \"subject\"")
			return nil
		}
		if messageData.Body == "" {
			app.errorLogger.Println("failed to send email, missing \"body\"")
			return nil
		}
		if messageData.To == "" {
			app.errorLogger.Println("failed to send email, missing \"to\"")
			return nil
		}

		message := app.mgTommyMayDev.NewMessage(messageData.Sender, messageData.Subject, "", messageData.To)
		message.SetHtml(messageData.Body)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		_, _, err := app.mgTommyMayDev.Send(ctx, message)
		if err != nil {
			app.errorLogger.Printf("failed to send email: %v\n", err)
			return err
		}
		app.infoLogger.Printf("email sent: sender: %s, subject: %s, to: %s\n", messageData.Sender, messageData.Subject, messageData.To)

		return nil
	}
}
