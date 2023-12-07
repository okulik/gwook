FROM golang:1.21.4-alpine as builder

WORKDIR /app
ADD . .

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  app

RUN apk add build-base curl
RUN make build-prod

FROM golang:1.21.4-alpine

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=0:0 /app/bin /

ENV HTTP_SERVER_PORT 4000
EXPOSE 4000

USER app:app

ENTRYPOINT ["/svixer"]
