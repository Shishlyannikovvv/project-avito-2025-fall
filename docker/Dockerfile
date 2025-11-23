FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:latest
COPY --from=builder /api /api
EXPOSE 8080
CMD ["/api"]
