# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY personnel-service/go.mod ./go.mod
COPY personnel-service/go.sum ./go.sum
COPY hospital-shared ./hospital-shared

RUN go mod download

COPY personnel-service/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

# Run stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /main .
COPY --from=builder /app/configs ./configs

EXPOSE 8083

ENTRYPOINT ["/main"]
