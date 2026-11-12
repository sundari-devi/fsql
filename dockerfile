FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/fsql ./cmd/fsql

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fsql /app/fsql
CMD ["/bin/sh"]
