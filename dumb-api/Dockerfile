# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.23 AS builder

RUN mkdir -p /src/dumb-api
WORKDIR /src/dumb-api

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
RUN go mod vendor

ADD . .
RUN wget https://github.com/gobuffalo/cli/releases/download/v0.18.14/buffalo_0.18.14_Linux_x86_64.tar.gz
RUN tar -xvzf buffalo_0.18.14_Linux_x86_64.tar.gz
RUN mv buffalo /usr/local/bin/
RUN buffalo build --static -o /bin/app -mod=vendor
# Install the soda cli
RUN go install github.com/gobuffalo/pop/v6/soda@latest

FROM golang:1.23-alpine as final
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .
COPY --from=builder src/dumb-api/migrations ./migrations
COPY --from=builder src/dumb-api/database.yml .
COPY --from=builder /go/bin/soda /usr/local/bin/