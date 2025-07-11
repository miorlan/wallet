FROM golang:1.23 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet-app ./cmd/

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /wallet-app .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./wallet-app"]