# Standalone
Run [package email][package-email] as a standalone function/server. This is the preferred way to deploy this package
and is how the [package author][package-author] uses it in production.

```go
package main

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/itmayziii/email/send"
	"github.com/mailgun/mailgun-go/v4"
	"gocloud.dev/blob"
	"log"
	"os"
)

import (
	_ "github.com/joho/godotenv/autoload"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/gcsblob"
)

func main() {
	ctx := context.Background()
	// You could register this at any path but the root path is appropriate is your running this as a microservice
	// with only one responsibility.
	if err := funcframework.RegisterCloudEventFunctionContext(ctx, "/", registerFunction(ctx)); err != nil {
		log.Fatalf("funcframework.RegisterCloudEventFunctionContext: %v\n", err)
	}

	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

func registerFunction(ctx context.Context) func(context.Context, cloudevents.Event) error {
	loggingClient, err := logging.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal("failed to create logging client", err)
	}
	logger := loggingClient.Logger("email", logging.RedirectAsJSON(os.Stdout))
	infoLogger := logger.StandardLogger(logging.Info)
	errorLogger := logger.StandardLogger(logging.Error)

	bucketName := os.Getenv("BUCKET")
	// Do not worry about closing the bucket, this is a no-op for all the implementations that we use i.e. file, GCS, S3.
	// It would be more difficult to close the bucket because function invocations sharing the same cloud function would
	// not be able to share resources.
	bucket, err := blob.OpenBucket(ctx, bucketName)
	if err != nil {
		logger.LogSync(ctx, logging.Entry{
			Severity: logging.Alert,
			Payload:  fmt.Sprintf("failed to open bucket %s - %v", bucketName, err),
		})
		os.Exit(1)
	}

	app := send.NewApp(
		send.AppWithFlusher(logger),
		send.AppWithInfoLogger(infoLogger),
		send.AppWithErrorLogger(errorLogger),
		send.AppWithFileStorage(bucket),
		send.AppWithDomainSender("tommymay.dev", send.NewMailgunSender(
			mailgun.NewMailgun("mg.tommymay.dev", os.Getenv("MG_API_KEY_MG_TOMMYMAY_DEV")),
		)),
	)

	return send.EmailEvent(app)
}
```

## Deployment
The easiest way to deploy is via [GCP Cloud Functions.][gcp-cloud-functions] This package was specifically designed
with [Cloud Functions 2nd generation][gcp-cloud-functions-2-gen] in mind. That being said, there is nothing stopping
you from running this code anywhere you want, it is simply an HTTP server that responds to [CloudEvents over HTTP]
[cloud-events-http]. The HTTP server can live anywhere and the CloudEvents can be delivered by anything that can
perform an HTTP request.

_Example deploy script with [gcloud CLI][gcloud] to deploy as a Cloud Function._
```shell
PROJECT_ID="itmayziii"
REGION="us-central1"
FUNCTION_SA="app-email-func@itmayziii.iam.gserviceaccount.com"
RUN_SA="app-email-func@itmayziii.iam.gserviceaccount.com"
TRIGGER_SA="eventarc-trigger@itmayziii.iam.gserviceaccount.com"
# Pub/sub topic will need created separately.
TOPIC="send-email"
VERSION="v1_0_0"
APP=email
```
_Create pub/sub topic_
```shell
gcloud pubsub topics create send-email --labels="managed_by=manual,app=email"
```
_Deploy cloud function_
```shell
gcloud functions deploy SendEmail \
  --project="$PROJECT_ID"
  --gen2 \
  --trigger-topic="$TOPIC" \
  --runtime="go121" \
  --entry-point="SendEmail" \
  --region="$REGION" \
  --source="." \
  --ingress-settings="internal-only" \
  --no-allow-unauthenticated \
  --retry \
  --trigger-service-account="$TRIGGER_SA" \
  --run-service-account="$RUN_SA" \
  --service-account="$FUNCTION_SA" \
  --set-secrets="MG_API_KEY_MG_TOMMYMAY_DEV=mailgun-api-key-mg-tommymay-dev:latest" \
  --set-env-vars="PROJECT_ID=$PROJECT_ID,BUCKET=gs://itmayziii-email-templates" \
  --clear-labels \
  --update-labels="managed_by=manual,version=$VERSION"
```

[package-author]: https://github.com/itmayziii
[package-email]: https://github.com/itmayziii/email
