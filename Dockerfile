FROM golang:1.16.4-alpine AS builder

WORKDIR /build
ENV CGO_ENABLED=1 \
  GOOS=linux \
  GOARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o medusa-api ./cmd/medusa-api/main.go
RUN go build -o medusa-engine ./cmd/medusa-engine/main.go
RUN go build -o medusa-daemon ./cmd/medusa-daemon/main.go

FROM alpine:3.13.6

RUN apk add ca-certificates
WORKDIR /app

COPY --from=builder /build/config/avatar.png ./config/avatar.png
COPY --from=builder /build/config/mailer.yaml ./config/mailer.yaml
COPY --from=builder /build/config/mailer ./config/mailer
COPY --from=builder /build/medusa-api ./
COPY --from=builder /build/medusa-engine ./
COPY --from=builder /build/medusa-daemon ./
