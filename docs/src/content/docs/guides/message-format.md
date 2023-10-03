---
title: Message Format
description: Learn how to format the message content sent in CloudEvents.
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

All the top level attributes are defined by the CloudEvents specification including the `data` attribute. The only
caveat is the content of the `data` attribute is left open for applications to define themselves. This package has made
the decision to follow another existing specification which is [GCP pub/sub message format.][gcp-pub-sub-message]

This means that the `data.message` attribute matches the [PubsubMessage structure][gcp-pub-sub-message] and the
`data.message.data` attribute is a base64 encoded string containing the actual information needed to send an email
such as the sender and recipients.

## Application Specific Attributes
These are the attributes that should be base64 encoded and put into `message.data`.

| Attribute | Type             | Description                                                      |
|-----------|------------------|------------------------------------------------------------------|
| sender    | string           | Who the email is coming from                                     |
| subject   | string           | What the email is about                                          |
| body      | string           | HTML body of the email, alternatively provide "template"         |
| to        | []string         | Who the email should go to                                       |
| template  | string           | Go HTML template path                                            |
| data      | map[string][any] | Arbitrary variables you want to bind to the "body" or "template" |


[cloud-events]: https://cloudevents.io/
[cloud-event-goals]: https://github.com/cloudevents/spec/blob/main/cloudevents/primer.md#design-goals
[cloud-event-http]: https://github.com/cloudevents/spec/blob/main/cloudevents/bindings/http-protocol-binding.md#32-structured-content-mode
[gcp-pub-sub-message]: https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
