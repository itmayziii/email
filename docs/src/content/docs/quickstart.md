---
title: Quickstart
description: Get up and running quickly with the email package.
---

## Functions Framework
The easiest way to get started is to run this package as a standalone HTTP server which can be done easily with the
[Functions Framework for Go.][ff-go]

1. Create a file to represent your server, lets say `function.go`.

    _function.go_
    ```go
    package main
    
    import (
        "context"
        "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
        cloudevents "github.com/cloudevents/sdk-go/v2"
        "github.com/itmayziii/email/send"
        "gocloud.dev/blob/memblob"
        "log"
        "os"
    )
    
    func main() {
        ctx := context.Background()
        if err := funcframework.RegisterCloudEventFunctionContext(ctx, "/", emailEvent(ctx)); err != nil {
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
    
    // emailEvent returns a function responsible for sending emails based on CloudEvent data.
    func emailEvent(ctx context.Context) func(context.Context, cloudevents.Event) error {
        infoLogger := log.New(os.Stdout, "info - ", log.Ltime)
        errorLogger := log.New(os.Stdout, "error - ", log.Ltime)
    
        bucket := memblob.OpenBucket(nil)
        app := send.NewApp(
            send.AppWithInfoLogger(infoLogger),
            send.AppWithErrorLogger(errorLogger),
            send.AppWithFileStorage(bucket),
            send.AppWithDomainSender("example.com", send.NoopSender{}),
        )
    
        return send.EmailEvent(app)
    }
    ```

2. Run `go mod tidy` to install dependencies.

3. Start the server with `go run function.go`

4. Send a sample [CloudEvent][cloud-events] to the server at `127.0.0.0:8080`
    ```shell
    curl --location 'localhost:8080' \
    --header 'ce-id: 1096434104173400' \
    --header 'ce-source: //pubsub.googleapis.com/projects/example-project/topics/email' \
    --header 'ce-specversion: 1.0' \
    --header 'ce-type: someType' \
    --header 'ce-time: 2020-12-20T13:37:33.647Z' \
    --header 'Content-Type: application/json' \
    --data-raw '{
      "sender": "no-reply@example.com",
      "subject": "hello world",
      "body": "some body",
      "to": ["tom@example.com"]
    }'
    ```

    You will see output from the server similar to
    `info - 21:22:41 email sent: sender: no-reply@example.com, subject: hello world, to: [tom@example.com]`.
    You will notice that none of the data in the output is present in the CloudEvent we sent to the server and this is
    because we are following the [GCP pub/sub message format][pubsub-message-format] which specifies the application
    data be base64 encoded. If you look at `message.data` you will see the base64 encoded string. You can read more
    about the expected message format in the [message format guide][message-format-guide]. 

[ff-go]: https://github.com/GoogleCloudPlatform/functions-framework-go
[cloud-events]: https://cloudevents.io/
[pubsub-message-format]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
[message-format-guide]: /guides/message-format/
