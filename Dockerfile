# Borrowed heavily from https://webdevstation.com/posts/one-of-thee-easiest-ways-to-host-go-web-apps/
# This will *specifically* build json-prettifier-web into a Docker container

FROM golang:alpine AS builder
RUN apk add --no-cache --update \
        git \
        ca-certificates
EXPOSE 8080
ADD . /app
WORKDIR /app
COPY go.mod ./
RUN  go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main cmd/json-prettifier-web/main.go

FROM alpine
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
