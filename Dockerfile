#build gor go object file
FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o main ./cmd/itpl_server/main.go

# exec stage 
FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main /usr/local/bin/main

EXPOSE 8080

CMD ["main"]
