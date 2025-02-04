# Copyright 2019 Core Services Team.

FROM golang:1.23-alpine as builder

RUN apk add --no-cache ca-certificates git

WORKDIR /account
COPY go.mod .
COPY go.sum .
RUN go mod download
ENV DB_CONNECTION_STRING postgres://postgres:postgres@postgres:5432/account_db?sslmode=disable

COPY . .
RUN CGO_ENABLED=0 go install ./cmd/account

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/bin /bin
USER nobody:nobody
ENTRYPOINT ["/bin/account"]