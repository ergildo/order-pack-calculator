FROM golang:1.24-alpine AS builder

WORKDIR  /build

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /build/order-pack-calculator-api ./cmd/api...

FROM alpine:3.18

WORKDIR  /app

COPY --from=builder /build/order-pack-calculator-api /app/order-pack-calculator-api 
COPY --from=builder /build/migrations /app/migrations

EXPOSE 8080

CMD ["/app/order-pack-calculator-api"]