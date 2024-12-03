FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server

FROM golang:1.23

WORKDIR /

COPY --from=builder /app/server .

CMD ["./server"]
