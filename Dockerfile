FROM golang:1.12-alpine

WORKDIR /go/src/app
COPY . .

RUN go install -v ./...

# Env vars:
# Set DESTINATION to ip:port
# Set LISTEN_PORT to :port
# Set IS_RECEIVER to true on start/end box

CMD ["around-the-world-in-80"]