
FROM golang:1.22-alpine AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /broker ./cmd/broker

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

COPY --from=builder /broker .
COPY config.json .

EXPOSE 8080

ENTRYPOINT ["/app/broker"]