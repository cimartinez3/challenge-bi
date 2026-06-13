FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /novobanco ./cmd/api

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
RUN adduser -D -u 1001 appuser
USER appuser
COPY --from=builder /novobanco /novobanco
EXPOSE 8080
ENTRYPOINT ["/novobanco"]
