FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /receipt-processor cmd/main.go

FROM alpine:latest
WORKDIR /
COPY --from=builder /receipt-processor /receipt-processor
EXPOSE 8080
CMD ["/receipt-processor"]
