FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . ./

RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/server .

EXPOSE 3000

ENV PORT=3000

STOPSIGNAL SIGINT

CMD ["./server"]
