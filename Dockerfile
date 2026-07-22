FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY server/go.mod server/go.sum ./server/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /signaling-server ./server/cmd/server

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

COPY --from=builder /signaling-server /usr/local/bin/signaling-server

EXPOSE 9321

ENTRYPOINT ["signaling-server"]
CMD ["--port", "9321"]
