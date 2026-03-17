FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o rider-service ./cmd/api

FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app
WORKDIR /home/app
COPY --from=builder /app/rider-service ./rider-service
COPY docs ./docs
USER app
EXPOSE 8084
ENTRYPOINT ["/home/app/rider-service"]
