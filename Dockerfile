FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

#============= 2 Минимальный образ

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app 

COPY --from=builder /app/server .

COPY migrations ./migrations

EXPOSE 3001

CMD ["./server"]