FROM golang:1.23-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init ddos
RUN go mod tidy
RUN go build -ldflags="-s -w" -o app

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/app /app/app

EXPOSE 80
ENTRYPOINT ["/app/app"]
