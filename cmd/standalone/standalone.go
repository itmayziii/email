/*
Package main is meant for local testing purposes only, it sets up a minimalistic standalone server to send CloudEvents
to.
*/
package main

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/itmayziii/email/send"
	"gocloud.dev/blob"
	"log"
	"os"
)

import (
	_ "github.com/joho/godotenv/autoload"
	_ "gocloud.dev/blob/fileblob"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterCloudEventFunctionContext(ctx, "/", emailEvent(ctx)); err != nil {
		log.Fatalf("funcframework.RegisterCloudEventFunctionContext: %v\n", err)
	}

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

// emailEvent returns a function responsible for sending emails based on CloudEvent data.
func emailEvent(ctx context.Context) func(context.Context, cloudevents.Event) error {
	infoLogger := log.New(os.Stdout, "info - ", log.Ltime)
	errorLogger := log.New(os.Stdout, "error - ", log.Ltime)

	bucketName := os.Getenv("BUCKET")
	bucket, err := blob.OpenBucket(ctx, bucketName)
	if err != nil {
		errorLogger.Printf("failed to open bucket %s", bucketName)
		os.Exit(1)
	}

	app := send.NewApp(
		send.AppWithInfoLogger(infoLogger),
		send.AppWithErrorLogger(errorLogger),
		send.AppWithFileStorage(bucket),
		send.AppWithDomainSender("example.com", send.NoopSender{}),
	)

	return send.EmailEvent(app)
}
