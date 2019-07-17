FROM golang:1.12-alpine

RUN apk add git
RUN apk --update add ca-certificates

RUN adduser -D app

USER app
WORKDIR /home/app

COPY --chown=app:app . .

RUN go build -o twitter-serverless-serving .

ENTRYPOINT ["/home/app/twitter-serverless-serving"]
