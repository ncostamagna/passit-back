FROM golang:1.25.0-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ./bin/passit-back ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/bin/passit-back /app/passit-back

RUN chmod +x /app/passit-back

ENTRYPOINT ["/app/passit-back"]
