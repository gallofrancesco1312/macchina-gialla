FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o macchina-gialla .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/macchina-gialla .
COPY counter.json .
CMD ["./macchina-gialla"]
