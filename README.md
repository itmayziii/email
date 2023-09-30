# email Package
Package `email` responds to [CloudEvents][cloud-events] to send emails over HTTP by implementing the 
[CloudEvents HTTP specification.][cloud-events-http]

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

## Why?
* Storing secrets necessary to send emails  in many projects is tedious, it is easier to centralize that secret access.
* Sending emails in apps often involves setting up some sort of messaging queue and templating engine to handle
  variables in emails. Doing this for every app is again, tedious.
* Centralizing email templates gives transparency of what exists already.
* Removing vendor lock-in from email SaaS products is ideal. Switching between Mailgun, Sendgrid, Mandrill, etc...

## Features
This package is designed to facilitate email sending in a vendor-neutral manner, avoiding any vendor lock-in.
We achieve this goal through the following features:

### Sending Options:
Providing users with the flexibility to choose their preferred method of sending emails, whether it's through a SaaS
solution like [Mailgun][mailgun] or [SendGrid][sendgrid], or by utilizing their own SMTP server.

### Template Storage:
Allowing users to decide where they want to store email templates, whether it's on their local disk or in cloud 
storage services like [S3][s3] or [Google Cloud Storage (GCS).][gcs]

### Standardization
Adhering to the [CloudEvent specification][cloud-events] to ensure a consistent and opinionated structure for event 
data.

## Roadmap
* Add support for some sort of database to keep track of emails already sent to ensure only once delivery. This
  current implementation relies on [GCP pub/sub][gcp-pub-sub] which does not support [exactly once delivery]
  [gcp-exactly-once-delivery] for push subscriptions and cloud functions rely on push subscriptions.
* Provide ready to go adapters for other SaaS providers besides Mailgun.
  * Sendgrid
  * Mandrill

[mailgun]: https://www.mailgun.com/
[mailgun-api]: https://documentation.mailgun.com/en/latest/api_reference.html
[sendgrid]: https://sendgrid.com/
[go-html-template]: https://pkg.go.dev/html/template
[cloud-events]: https://cloudevents.io/
[cloud-events-http]: https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/bindings/http-protocol-binding.md
[cloud-dev-blob]: https://gocloud.dev/howto/blob/
[gcs]: https://cloud.google.com/storage
[s3]: https://aws.amazon.com/s3/
[gcp-cloud-functions]: https://cloud.google.com/functions
[gcp-cloud-functions-2-gen]: https://cloud.google.com/blog/products/serverless/cloud-functions-2nd-generation-now-generally-available
[gcloud]: https://cloud.google.com/sdk/gcloud
[gcp-pub-sub]: https://cloud.google.com/pubsub
[gcp-exactly-once-delivery]: https://cloud.google.com/pubsub/docs/exactly-once-delivery
