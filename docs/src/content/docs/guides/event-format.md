---
title: Message Format
description: Learn how to format the message content to send an email.
---

[CloudEvents][cloud-events] gives us a starting point for defining how events should be structured, but it is not
all encompassing, and it is [not trying to be][cloud-event-goals].

> CloudEvents, at its core, defines a set of metadata, called "context attributes", about the event being transferred 
  between systems, and how those pieces of metadata should appear in that message. This metadata is meant to be the 
  minimal set of information needed to route the request to the proper component and to facilitate proper processing 
  of the event by that component. So, while this might mean that some of the application data of the event itself 
  might be duplicated as part of the CloudEvent's set of attributes, this is to be done solely for the purpose of 
  proper delivery, and processing, of the message. Data that is not intended for that purpose should instead be 
  placed within the event (data) itself.

Let's clarify the above statement with an example by looking at
[CloudEvents over HTTP in "structured content mode."][cloud-event-http]

```http
POST / HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/cloudevents+json; charset=utf-8
Content-Length: 681

{
    "id": "1096434104173400",
    "source": "//pubsub.googleapis.com/projects/example-project/topics/email",
    "specversion": "1.0",
    "type": "google.cloud.pubsub.topic.v1.messagePublished",
    "time": "2020-12-20T13:37:33.647Z",
    "data": {
        "sender": "no-reply@example.com",
        "subject": "hello world",
        "body": "some body",
        "to": ["tom@example.com"]
    }
}
```

All the top level attributes are defined by the CloudEvents specification including the `data` attribute. The only
caveat is the content of the `data` attribute is left open for applications to define themselves.

## Application Specific Attributes

| Attribute | Type                          | Description                                                      |
|-----------|-------------------------------|------------------------------------------------------------------|
| sender    | string                        | Who the email is coming from                                     |
| subject   | string                        | What the email is about                                          |
| body      | string (optional w/ template) | HTML body of the email, alternatively provide "template"         |
| to        | []string                      | Who the email should go to                                       |
| template  | string (optional w/ body)     | Go HTML template path                                            |
| data      | map[string][any]              | Arbitrary variables you want to bind to the "body" or "template" |
| cc        | []string (optional)           | Who will be carbon copied on the email                           |
| bcc       | []string (optional)           | Who will be blind carbon copied on the email                     |


## Other Message Formats
Some event producers have a defined way they produce payloads and while it would not be possible for this library
to accommodate every format, we will aim to make it easy to work with the most popular ones.

### GCP Evetnarc + Pub/Sub 
When receiving events produced by GCP pub/sub the event payload will be wrapped in a `message` attribute and 
`message.data` will be base64 encoded data which should contain our application specific attributes. The `data.message`
attribute matches [PubsubMessage structure.][gcp-pub-sub-message] 

The CloudEvent produced by [Eventarc / GCP Pub/Sub][eventarc] is going to look like this:
```http
POST / HTTP/1.1
Host: 127.0.0.1:8080
Content-Type: application/cloudevents+json; charset=utf-8
Content-Length: 681

{
    "id": "1096434104173400",
    "source": "//pubsub.googleapis.com/projects/example-project/topics/email",
    "specversion": "1.0",
    "type": "google.cloud.pubsub.topic.v1.messagePublished",
    "time": "2020-12-20T13:37:33.647Z",
    "data": {
        "message": {
            "attributes": {
                "key": "value"
            },
            "data": "eyJzZW5kZXIiOiJuby1yZXBseUBleGFtcGxlLmNvbSIsInN1YmplY3QiOiJoZWxsbyB3b3JsZCIsImJvZHkiOiJzb21lIGJvZHkiLCJ0byI6WyJzb21lYm9keUBleGFtcGxlLmNvbSJdfQo=",
            "messageId": "2070443601311540",
            "publishTime": "2021-02-26T19:13:55.749Z",
        },
        "subscription": "projects/myproject/subscriptions/mysubscription"
    }
}
```

This package already unwraps the data from `message.data` for you and is compatible with this type of event format.
We detect the CloudEvent type and look for `google.cloud.pubsub.topic.v1.messagePublished` to determine if this 
unwrapping needs done.

## Examples
Check out the public postman collection to see how to send CloudEvents over HTTP in both binary and structured data
mode.

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/135269-e02c0d1c-05d4-4cbe-b3e6-edc2d88a7dd1?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D135269-e02c0d1c-05d4-4cbe-b3e6-edc2d88a7dd1%26entityType%3Dcollection%26workspaceId%3Dfd4b13b1-1b61-4a2a-9a77-f7e2158f0514)

[cloud-events]: https://cloudevents.io/
[cloud-event-goals]: https://github.com/cloudevents/spec/blob/main/cloudevents/primer.md#design-goals
[cloud-event-http]: https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/http-protocol-binding.md#32-structured-content-mode
[gcp-pub-sub-message]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
[eventarc]: https://cloud.google.com/eventarc/docs/overview
