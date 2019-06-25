FROM golang:1.12.6-alpine AS builder
WORKDIR /go/src/github.com/takashabe/edge-exporter/example/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/takashabe/edge-exporter/example/app .
CMD ["./app"]
