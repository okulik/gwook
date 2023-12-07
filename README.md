# Hello, Gwook!

To allow the delivery of event notifications to custom endpoints, a new service called Gwook is introduced here. This service exposes an HTTP endpoint `/notifications`, responsible for forwarding all incoming events - received, for instance, from a message such as GCP Pub/Sub or AWS Kinesis - to [Svix](https://www.svix.com), a third-party webhooks as a service provider.

Ideally, Gwook would directly receive events from those message busses by subscribing to multiple topics using gRPC streaming, but by implementing its own HTTP endpoint allows it to stay flexible and independent.

# Gwook Architecture

Gwook hosts a web server that processes incoming events at the `/notifications` endpoint, secured by basic authentication. Additionally, it features an authentication-free health check endpoint located at `/health`. While basic authentication is currently employed, a more flexible token-based authentication method would be preferable. However, for the sake of simplicity, basic authentication is used in the current setup. The inclusion of the health check endpoint serves the purpose of monitoring the service's health and uptime, which is standard for ECS or k8s liveness and readiness probes, particularly in containerized environments.

The service works using Golang's standard library HTTP server, which spawns a separate goroutine for each connection. While this design enables the handling of multiple concurrent requests, it creates limitations on managing a high volume of concurrent connections. To accommodate a larger number of concurrent connections, an alternative HTTP server implementation, such as [github.com/valyala/fasthttp](https://github.com/valyala/fasthttp), would be better.

## Events validation

Upon receiving the event, Gwook initiates the validation process. Presently, we use the `model.Notification` struct, manually created from some example events in the test/events folder, specifically from OpenAPI's `event-spec.json` file. This struct unmarshals the event from the POST request's body. While all events are intended to reflect same common structure, variations may exist within the `Data` field, depending on the event type. For a robust production-level service, employing a tool to generate Golang structs for all event types and utilizing reflection to differentiate between received event types would be necessary.

An even more efficient approach would involve replacing JSON with protobuf definitions to serialize messages received from the /notifications endpoint. Apart from checking the event's structure and field types, we should also examine field contents and add further validations if necessary. For JSON-formatted events, we could use the [github.com/go-playground/validator](https://github.com/go-playground/validator) validator, while for protobuf, the [github.com/bufbuild/protoc-gen-validate](https://github.com/bufbuild/protoc-gen-validate) validator.

## Events filtering

It's worth mentioning that perhaps not all events should be forwarded to Svix, as we might want to keep certain events internal. Additionally, we might not be interested in all types of events. Our current implementation forwards all events to Svix. However, for a production-level service, this functionality should be configurable for every type of use separately and typically stored in a database.

## Sending events to Svix

After receiving, validating, and filtering the event, the subsequent step involves sending the event to Svix via the [Create Message](https://api.svix.com/docs#tag/Message/operation/v1.message.create) API. However, prior to this, we must establish a mapping between a fixtional user ID and the corresponding Svix application ID. This mapping can be stored in a database and managed via an API.

In our current implementation, as a simplification, we're using a single Svix application ID for all users that we're reading from the `SVIX_APPLICATION_ID` environment variable.

### Svix applications and endpoints

The Svix application serves as an essential entity required for each user intending to receive notifications. This application can be automatically created upon user signup or enabled on demand via an API endpoint implemented elsewhere. Once the Svix application object is created, multiple endpoint objects can be attached to it. Consequently, when a message is generated in Svix for a particular application, a singular notification is dispatched to all associated endpoints. However, in our current implementation, we're disregarding the Svix application and endpoint objects, directing all event messages to Svix.

### Elimination of duplicate messages

To prevent unnecessary triggering of webhooks, it's crucial to avoid sending duplicate messages.

Svix API provides an option to specify an idempotency key, uniquely identifying each event. Svix's idempotency function records the resulting status code and body from the initial request made with any given idempotency key for any successful request. Subsequent requests using the same key yield the same result. Currently, we employ the event ID as the key, ensuring uniqueness for each event.

However, if we prefer to establish idempotency based on a combination of different event properties with an expiration value, an alternative strategy can be used. For instance, by generating the idempotency key using a formula like `event.type + event.project + event.source`, we could utilize a distributed cache such as Redis with an expiration time, say 1 minute. Before sending an event to Svix, we would check the cache to determine if the key exists. If it does, we would cancel the sending process. If the key is absent, we would add it to the cache and proceed to send the event to Svix. This method ensures that identical events are not dispatched to Svix multiple times within a short timeframe.

# API errrors

Follows a list of possible errors that Gwook API can return.

- `400 Bad Request` if we're sending invalid event or the request body is not a valid JSON
- `401 Unauthorised` if we fail to provide correct basic authentication credentials
- `409 Conflict` if we're sending the same message multiple times
- `429 Too Many Requests` if Svix is rate-limiting us
- `500 Internal Server` if Svix API is down, even after retying

It's important to note that while some errors are directly returned by the Gwook service, others are reflected by the Svix API. For instance, the `409 Conflict` error and `429 Too Many Requests` (we don't limit the Gwook service as it's an internal API) are returned by the Svix API. Additionally, a `500 Internal Server` error may be returned when there's an issue with the Svix API. Detailed error information is included as JSON in the response body to provide more context about the source of the error.

Received errors that can be categorized as non-recoverable (such as 400, 401, 409) should not be retried but logged and forwarded to external error tracking services like Honeybadger or Sentry. Even after retry attempts, if we continue to encounter the 429 error, it should also be logged and sent to the error tracking service.

# A wish list

There's a number of very important features that are currently missing in the current implementation. For instance, we are missing a number of useful APIs. Additionally, we're not utilizing any external error tracking or monitoring services and we're not employing any tracing or metrics collection tools. Follows a list of features that would make Gwook more useful, maintainable, and production-ready.

- **Application & Webhook Endpoint Management API**: Establish an API enabling users to manage applications and their webhook endpoints, managing mappings within a relational database. Basically creating a public wrapper around certain Svix API aspects to protect our internal implementation.
- **Event Filtering API**: Develop an API enabling users to filter events based on type or other properties, utilizing a relational database for mappings.
- **Event Metrics API**: Create an API enabling users some notifications metrics, such as the number of events forwarded to webhooks.
- **Event Struct Management**: Utilize tooling to generate Golang structs for Pub/Sub events or transition from JSON to protobuf definitions for serialization of Pub/Sub messages.
- **Token-based Authentication**: Replace basic authentication with token-based authentication, storing user API tokens in a database and offering an API to manage them.
- **Enhanced Testing Suite**: Introduce integration tests, employing tools like VCR to record and replay HTTP interactions with Svix (sandbox/production) for validating code behavior with actual API responses.
- **Error Tracking Service Integration**: Incorporate an external error tracking service like Sentry or Honeybadger for reporting service errors.
- **Telemetry Implementation**: Integrate tracing using tools like New Relic, Datadog, Jaeger, or Tempo to trace requests. Also incorporate metrics and monitoring utilizing New Relic, Datadog, Prometheus, VictoriaMetrics, or Grafana and add logging, utilizing a structured logging library like [zap](https://github.com/uber-go/zap).
- **Dead-letter Topic Configuration**: Set up a dead-letter topic in Pub/Sub for handling messages failing delivery to Svix even after retries (e.g., 429 or 500 errors), enabling later retries and monitoring of failed messages.
- **Containerized Deployment**: Run the service in a containerized environment (AWS ECS, k8s) and behind a load balancer, for scaling purpose. Deploy the container image to a private registry like AWS ECR or GCP Artifact Registry rather than GitHub Container Registry.
- **Message ID Uniqueness & Distributed Cache**: Optionally ensure message ID uniqueness for short intervals using a distributed cache like Redis.
- **HTTP Server Optimization**: Optionally consider using alternative HTTP server implementations like [fasthttp](https://github.com/valyala/fasthttp) to handle a high number of concurrent connections.
- **Webhook as a Service Alternative**: Optionally explore alternative webhook-as-a-service providers like [webhookrelay](https://webhookrelay.com/) or [hook0](https://www.hook0.com). If Svix remains to be the the only reasonable option, perhaps implementing our own webhook mechanism to avoid long-term risk makes sense.
- **Svix library HTTP Implementation**: Since the current design shares a single HTTP client instance across all requests, HTTP client used by the Svix library might become a bottleneck because of connection pooling. Consider using a separate HTTP client instance for each request or replace Svix library and call Svix REST API directly.

# Development notes

You can start the Gwook service in development in two ways: by executing the Makefile's target run or by running service's Docker container.

## Using the Makefile

To start the service locally, run the following command:
```bash
$ make run
```

## Using Docker

For running the containerized service, use the following command to build the container image:
```bash
$ docker build -t gwook .
```

Alternatively, run `make docker-build` to build the container image locally.

Then, run the container with the following command:
```bash
$ docker run -p 4000:4000 -e AUTH_USERNAME=admin AUTH_PASSWORD=admin gwook
```

Alternatively, run `make docker-run` to build and run the container image locally.

## Testing Gwook API from the shell

To sent a notification to Svix, run the following command:
```bash
$ curl -u admin:admin \
  -d @test/events/event-01.json \
  "http://localhost:4000/notifications"
```

## GitHub actions

The service is configured to run the tests and build the Docker image on every push to the `main` branch. Docker image is stored in the GitHub Container Registry and versioned with the git commit SHA.

## Customizing running the service

The Gwook service can be additionally customised by setting the following environment variables:
* AUTH_USERNAME - Sets basic authentication user name (required).
* AUTH_PASSWORD - Sets basic authentication password (required).
* AUTH_REALM - Sets authentication realm (required).

* HTTP_SERVER_GRACEFUL_SHUTDOWN_TIMEOUT - Sets the HTTP server's graceful shutdown timeout (optional, defaults to 10s).
* HTTP_SERVER_PORT - Sets the HTTP server's port (optional, defaults to 4000).
* HTTP_SERVER_IDLE_TIMEOUT - Sets the HTTP server's idle timeout (optional, defaults to 60s).
* HTTP_SERVER_READ_TIMEOUT - Sets the HTTP server's read timeout (optional, defaults to 10s).
* HTTP_SERVER_WRITE_TIMEOUT - Sets the HTTP server's write timeout (optional, defaults to 20s).

* SVIX_AUTH_TOKEN - Sets the Svix API authentication token (required).
* SVIX_DEBUG - Enables Svix API debug mode (optional, defaults to false).
* SVIX_SERVER_URL - Sets the Svix API URL (optional).
