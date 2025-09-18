# Build stage
FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./cmd/itpl_server/main.go

# Runtime stage
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Seoul

WORKDIR /app   

COPY --from=builder /app/main /usr/local/bin/main
COPY config.yaml /app/config.yaml

EXPOSE 8080

ENTRYPOINT ["main"]
