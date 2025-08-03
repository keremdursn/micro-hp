# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Shared modül dahil olsun diye bu şekilde yapıyoruz
COPY auth-service/go.mod ./go.mod
COPY auth-service/go.sum ./go.sum
COPY hospital-shared ./hospital-shared

RUN go mod download

COPY auth-service/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

# Run stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /main .
COPY --from=builder /app/configs ./configs

EXPOSE 8081

ENTRYPOINT ["/main"]
