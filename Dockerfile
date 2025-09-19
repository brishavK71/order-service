# build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /order-service

# final
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /order-service /usr/local/bin/order-service
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/order-service"]
