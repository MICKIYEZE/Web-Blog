FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY templates ./templates

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o blog .

FROM alpine:3.19

RUN addgroup -S blog && adduser -S blog -G blog
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/blog /app/blog
COPY --from=builder /app/templates /app/templates

RUN chown -R blog:blog /app

USER blog

EXPOSE 8080
CMD ["./blog"]
