---
title: Customizing Guide
description: Learn how to customize the email package.
---

## Custom Logger
You can provide a custom `info` and `error` logger to this email package to use. By default, nothing will be logged
if you do not provide one. That is to say if you only want errors to be logged then only provide an `error` logger.

This package does not take an opinionated stance when it comes to logging. Logging is one of those things that has
become a holy war and nobody can agree what a "logger" should look like. In the spirit of not deviating from the Go
standard library too much we use `*log.Logger`.

### [Standard Logger][standard-logger]
```go
package main

import (
    "log"
)

func main() {
    infoLogger := log.New(os.Stdout, "info - ", log.Ltime)
    errorLogger := log.New(os.Stdout, "error - ", log.Ltime)
    
    app := send.NewApp(
    send.AppWithInfoLogger(infoLogger),
    send.AppWithErrorLogger(errorLogger)
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
		send.AppWithFlusher(ZapFlusher{
			zap: logger,
        }),
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
		sned.AppWithFlusher(logger),
		send.AppWithInfoLogger(infoLogger),
		send.AppWithErrorLogger(errorLogger),
	)

	send.EmailEvent(app)
}
```

[standard-logger]: https://pkg.go.dev/log
[zap]: https://pkg.go.dev/go.uber.org/zap
[gcp-logging]: https://cloud.google.com/logging/docs/setup/go
