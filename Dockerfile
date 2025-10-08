FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o fintech

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/fintech .
EXPOSE 8080
CMD ["./fintech"]
