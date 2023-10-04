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
        "log"
        "os"
    )
    
    func main() {
        ctx := context.Background()
        // Register our function to be called when someone makes a POST request to the root path.
        if err := funcframework.RegisterCloudEventFunctionContext(ctx, "/", emailEvent(ctx)); err != nil {
            log.Fatalf("funcframework.RegisterCloudEventFunctionContext: %v\n", err)
        }
    
        // Use PORT environment variable, or default to 8080.
        port := "8080"
        if envPort := os.Getenv("PORT"); envPort != "" {
            port = envPort
        }
    
        // Start the HTTP server
        if err := funcframework.Start(port); err != nil {
            log.Fatalf("funcframework.Start: %v\n", err)
        }
    }
    
    // emailEvent returns a function responsible for sending emails based on CloudEvent data.
    func emailEvent(ctx context.Context) func(context.Context, cloudevents.Event) error {
        infoLogger := log.New(os.Stdout, "info - ", log.Ltime)
    
        app := send.NewApp(
            send.AppWithInfoLogger(infoLogger),
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
    You can read more about the expected event format in the [event format guide][event-format-guide].

5. The above example uses the `send.NoopSender` which is useful for testing but not so useful if you actually want to 
   send emails. Let's replace the `NoopSender` with an implementation that sends emails using [Mailgun.][mailgun]

    _function.go_
    ```go
        // Create a mailgun object, note that mg.example.com is the domain you register with Mailgun and not
        // the domain your sending from. i.e. Mailgun suggests setting Mailgun up at mg.SOME_DOMAIN.com so if you want
        // to send emails from SOME_DOMAIN.com then you should use mg.SOME_DOMAIN.com in the value here. 
        mailgun := mailgun.NewMailgun("mg.example.com", os.Getenv("MG_API_KEY"))
	    app := send.NewApp(
		    send.AppWithInfoLogger(infoLogger),
		    send.AppWithDomainSender("example.com", send.NewMailgunSender(mailgun)),
	    )
    ```
   
Check out the [Mailgun Go package][mailgun-go] for more information about using it.
   
Just like you configured Mailgun to send emails, there are [other customizations][customize-guide] you can provide
like what logger to use, where to read emails templates from, and of course other email providers.

[ff-go]: https://github.com/GoogleCloudPlatform/functions-framework-go
[cloud-events]: https://cloudevents.io/
[pubsub-message-format]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
[event-format-guide]: /guides/event-format/
[customize-guide]: /guides/customize/
[mailgun]: https://www.mailgun.com/
[mailgun-go]: https://github.com/mailgun/mailgun-go
