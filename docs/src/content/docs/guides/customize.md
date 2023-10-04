---
title: Customizing Guide
description: Learn how to customize the email package.
---

## Logger
You can provide a custom `info` and `error` logger to this email package to use. By default, nothing will be logged
if you do not provide one. That is to say if you only want errors to be logged then only provide an `error` logger.

This package does not take an opinionated stance when it comes to logging. Logging is one of those things that has
become a holy war and nobody can agree what a "logger" should look like. In the spirit of not deviating from the Go
standard library too much we use `*log.Logger`.

### [Standard Logger][standard-logger]
```go
package main

import (
	"github.com/itmayziii/email/send"
	"log"
)

func main() {
    infoLogger := log.New(os.Stdout, "info - ", log.Ltime)
    errorLogger := log.New(os.Stdout, "error - ", log.Ltime)
    
    app := send.NewApp(
        send.AppWithInfoLogger(infoLogger),
        send.AppWithErrorLogger(errorLogger),
    )
    
    send.EmailEvent(app)
}
```

### [Zap][zap]

```go
package main

import (
	"github.com/itmayziii/email/send"
	"go.uber.org/zap"
)

// ZapFlusher implements [Flusher] to clear the log buffer at the end of each email.
// This is really only useful if your running this in a serverless environment where flushing the logs before
// ending your code is necessary.
type ZapFlusher struct {
	zap *zap.Logger
}

func (z ZapFlusher) Flush() error {
	return z.zap.Sync()
}

func main() {
	logger, _ := zap.NewProduction()
	infoLogger := zap.NewStdLogAt(logger, zap.InfoLevel)
	errorLogger := zap.NewStdLogAt(logger, zap.ErrorLevel)

	app := send.NewApp(
		send.AppWithFlusher(ZapFlusher{zap: logger}),
	    send.AppWithInfoLogger(infoLogger),
		send.AppWithErrorLogger(errorLogger),
    )

	send.EmailEvent(app)
}
```

### [Google Cloud Logging][gcp-logging]
```go
package main

import (
	"cloud.google.com/go/logging"
	"github.com/itmayziii/email/send"
	"os"
)

func main() {
	loggingClient, err := logging.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatal("failed to create logging client", err)
	}
	logger := loggingClient.Logger("email-api", logging.RedirectAsJSON(os.Stdout))
	infoLogger := logger.StandardLogger(logging.Info)
	errorLogger := logger.StandardLogger(logging.Error)

	app := send.NewApp(
		send.AppWithFlusher(logger),
		send.AppWithInfoLogger(infoLogger),
		send.AppWithErrorLogger(errorLogger),
	)

	send.EmailEvent(app)
}
```

## Templating
This package relies on the [Go Cloud Development Kit blob package][blob] to abstract away handling file storage.
You can choose to store your templates on [disk][blob-disk], [GCS][blob-gcs], [S3][blob-s3], or any other option
supported by the blob package. Under the hood this package will use your configured blob to retrieve email templates
when the [template option][app-attributes] is provided.

_GCS example_
```go
package main

import (
	"context"
	"github.com/itmayziii/email/send"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/gcsblob"
	"os"
)

func main() {
	bucketName := os.Getenv("BUCKET")
	bucket, err := blob.OpenBucket(context.Background(), bucketName)
	if err != nil {
		os.Exit(1)
	}

	app := send.NewApp(send.AppWithFileStorage(bucket))

	send.EmailEvent(app)
}
```

## Email Providers
This package exposes an interface called `Sender` which can be implemented to do the actual sending of an email. 

```go
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
```

This package aims to make it easy to work with common email SaaS offerings by supplying adapters that already implement
this `Sender` interface. You have the flexibility to configure multiple email providers if you wish i.e. all emails
from `example.com` could use Mailgun and all emails from `example.org` could use Sendgrid.

### [Mailgun][mailgun] Adapter
```go
package main

import (
	"os"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/itmayziii/email/send"
)

func main() {
    mg := mailgun.NewMailgun("mg.example.com", os.Getenv("MG_API_KEY"))
    mgSender := send.NewMailgunSender(mg)
	app := send.NewApp(send.AppWithDomainSender("example.com", send.NewMailgunSender(mgSender)))

	send.EmailEvent(app)
}
```

The above example demonstrates how to use Mailgun to send emails. This package will match the ["from"][app-attributes]
attribute with the domain i.e. from: no-reply@example.com will matches the "example.com" domain which was configured
to use the provided mailgun object to send emails.

[standard-logger]: https://pkg.go.dev/log
[zap]: https://pkg.go.dev/go.uber.org/zap
[gcp-logging]: https://cloud.google.com/logging/docs/setup/go
[blob]: https://gocloud.dev/howto/blob/
[blob-disk]: https://gocloud.dev/howto/blob/#local
[blob-gcs]: https://gocloud.dev/howto/blob/#gcs
[blob-s3]: https://gocloud.dev/howto/blob/#s3
[app-attributes]: /guides/event-format/#application-specific-attributes
[mailgun]: https://www.mailgun.com/
