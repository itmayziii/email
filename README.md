# email Package
email package responds to [CloudEvents][cloud-events] to send emails.

Check out the [official documentation][docs] for end user documentation. This README is dedicated to information
about package development.

[![Static Badge](https://img.shields.io/badge/Documentation-green)][docs]
[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](https://commitizen.github.io/cz-cli/)
[![Go Reference](https://pkg.go.dev/badge/github.com/itmayziii/email.svg)](https://pkg.go.dev/github.com/itmayziii/email)

## Roadmap
* Add support for some sort of database to keep track of emails already sent to ensure exactly once delivery.
* Provide ready to go adapters for other SaaS providers besides Mailgun.
  * Sendgrid
  * Mandrill

## Contributing

### Local Setup
1. Setup environment variables by copying `.env.example` -> `.env` and filling out the missing environment variables.
    ```shell
    cp .env.example .env
    ```

2. There is a standalone HTTP server you can run with:
    ```shell
    go run cmd/standalone/standalone.go
    ```

    This standalone server uses the `NoopSender` so it does not actually send emails which is nice for testing.
    You can send CloudEvents via HTTP by using the public postman collection.

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/135269-e02c0d1c-05d4-4cbe-b3e6-edc2d88a7dd1?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D135269-e02c0d1c-05d4-4cbe-b3e6-edc2d88a7dd1%26entityType%3Dcollection%26workspaceId%3Dfd4b13b1-1b61-4a2a-9a77-f7e2158f0514)

### Releasing
This package uses [release-please][release-please] which will open a "release" pull request anytime something
releasable is merged into the `main` branch. Once the release pull request is merged there is a manual CI step the
package author will need to kick off to create a GitHub release.

[docs]: https://email-package.tommymay.dev
[cloud-events]: https://cloudevents.io/
[release-please]: https://github.com/googleapis/release-please
