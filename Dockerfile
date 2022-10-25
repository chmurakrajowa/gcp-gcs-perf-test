#
# usage-collector 
# Author: Damian Janiszewski
#
# Multistage Dockerfile definition
#
# Stage 0: golang compiler and build container
FROM golang:latest as builder-go
# FROM golang:1.15.15 as builder-go

WORKDIR /go/src/
RUN env

# Download dependencies
COPY go.mod go.sum /go/src/
RUN go mod download -json
RUN go list -m all

# Copy sources and compile
COPY gcp-gcs-perf-test.go config.go init.go /go/src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -tags 'static netgo' -ldflags '-w' gcp-gcs-perf-test.go config.go init.go

# Stage 1: running container 
FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN update-ca-certificates

LABEL version="0.1.2"
LABEL author "Damian Janiszewski"

# Copy binaries from stage 0 builder container
COPY --from=builder-go /go/src/gcp-gcs-perf-test /usr/local/bin/

CMD ["/usr/local/bin/gcp-gcs-perf-test"]
