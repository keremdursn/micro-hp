# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# go.mod ve shared module’u birlikte kopyala
COPY hospital-service/go.mod ./go.mod
COPY hospital-service/go.sum ./go.sum
COPY hospital-shared ./hospital-shared

# Modülleri indir
RUN go mod download

# Kaynak kodları kopyala
COPY hospital-service/. .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main cmd/main.go

# Run stage
FROM alpine:latest

WORKDIR /

COPY --from=builder /main .
COPY --from=builder /app/configs ./configs

EXPOSE 8082

ENTRYPOINT ["/main"]
