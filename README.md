# email
Module `email` responds to [CloudEvents][cloud-events] to send transactional emails over HTTP by implementing the 
[CloudEvents HTTP specification.][cloud-events-http]

[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

## Why?
* Storing secrets necessary to send emails in many projects is tedious, it is easier to centralize that secret access.
* Sending emails in apps often involves setting up some sort of messaging queue and templating engine to handle 
  variables in emails. Doing this for every app is again, tedious.
* Centralizing email templates gives transparency of what exists already.
* Removing vendor lock-in from email SaaS products is ideal. Switching between Mailgun, Sendgrid, Mandrill, etc...

## Features
* Send email (obviously) via multiple email SaaS providers.
* Specify a body or path to a template to use as the body.
  * Body and path can be a [Go template.][go-html-template] for you to provide data to bind to.
  * Path uses [cloud development blob][cloud-dev-blob] to allow for local, in memory, and cloud file storage such as
    [GCS][gcs] or [S3][s3].

### Roadmap
* Add support for some sort of database to keep track of emails already sent to ensure only once delivery. This
  current implementation relies on [GCP pub/sub][gcp-pub-sub] which does not support [exactly once delivery]
  [gcp-exactly-once-delivery] for push subscriptions and cloud functions rely on push subscriptions.
* Support other SaaS APIs besides Mailgun.
  * Sendgrid
  * Mandrill
* Support S3 for email template storage as we already use the [cloud development blob][cloud-dev-blob] package to allow 
  for many different cloud storage options. We just need to enable S3 support and test it.

## Deployment
The easiest way to deploy is via [GCP Cloud Functions.][gcp-cloud-functions] This package was specifically designed
with [Cloud Functions 2nd generation][gcp-cloud-functions-2-gen] in mind. That being said, there is nothing stopping
you from running this code anywhere you want, it is simply an HTTP server that responds to CloudEvents. The HTTP server
can live anywhere and the CloudEvents can be delivered by anything that can perform an HTTP request.

_Example deploy script with [gcloud CLI][gcloud] to deploy as a Cloud Function._
```shell
# Replace these variables with your own values
PROJECT_ID="itmayziii"
REGION="us-central1"
FUNCTION_SA="app-email-api@itmayziii.iam.gserviceaccount.com"
RUN_SA="app-email-api@itmayziii.iam.gserviceaccount.com"
TRIGGER_SA="eventarc-trigger@itmayziii.iam.gserviceaccount.com"
# Pub/sub topic will need created separately.
TOPIC="send-email"
VERSION="v1_0_0"
```
_Create pub/sub topic_
```shell
gcloud pubsub topics create send-email --labels="managed_by=manual,app=email"
```
_Deploy cloud function_
```shell
gcloud functions deploy SendEmail \
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
  --set-secrets="MG_API_KEY_MG_TOMMYMAY_DEV=mailgun-api-key-mg-tommymay-dev:latest"
  --set-env-vars="PROJECT_ID=$PROJECT_ID"
  --clear-labels
  --update-labels="managed_by=manual,version=$VERSION"
```

[mailgun-api]: https://documentation.mailgun.com/en/latest/api_reference.html
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
