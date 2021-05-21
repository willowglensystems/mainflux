# SentientQ Messages

The `sq-messaging` directory contains the SentientQ Message proto file and the message converter to convert SentientQ messages to Mainflux messages and vice versa.

SentientQ messages are of the following format:
| SentientQ   | Mainflux             | Note                                             |
|-------------|----------------------|--------------------------------------------------|
| Source      | Publisher            | Name of entity that is the source of the message |
| Timestamp   | Created              | Created timestamp                                |
| Body        | Payload              | Body of the message                              |
| Type        | Metadata.Type        | This a system type needed so that the receiver knows how to parse the body.   |
| Version     | Metadata.Version     | This is how Kubernetes defines its APIs e.g. ATC.v2 |
| ContentType | Metadata.ContentType | This is either json or cbor (binary) |
| Token       | Metadata.Token       | Used for our zero trust model to prevent access through the messaging system. Mainflux authenticates at the API request level. |
| Priority    | Metadata.Priority    | This can be moved to configuration since it is statically configured. If some use cases require this field to be changed at runtime then it could be added to the metadata of the message. |
| Initiator   | Metadata.Initiator   | This is parent message/service/entity of the current message. This field can be used to create a stack trace across messages. This field can be part of the metadata that gets appended as it goes to each service. |
| Metadata    | Metadata             | Additional metadata as required by implementations and deployments of SentientQ.  |