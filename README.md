# email Package
`email` package responds to [CloudEvents][cloud-events] to send emails.

Check out the [official documentation][docs] for end user documentation. This README is dedicated to information
about package development.

[![Static Badge](https://img.shields.io/badge/Documentation-green)][docs]
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](https://commitizen.github.io/cz-cli/)

## Roadmap
* Add support for some sort of database to keep track of emails already sent to ensure only once delivery. This
  current implementation relies on [GCP pub/sub][gcp-pub-sub] which does not support [exactly once delivery]
  [gcp-exactly-once-delivery] for push subscriptions and cloud functions rely on push subscriptions.
* Provide ready to go adapters for other SaaS providers besides Mailgun.
  * Sendgrid
  * Mandrill

[docs]: https://email-package.tommymay.dev
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
