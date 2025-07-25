# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

# Run stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /main .
COPY --from=builder /app/configs ./configs

EXPOSE 8081

ENTRYPOINT ["/main"]
