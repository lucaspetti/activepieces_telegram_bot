FROM golang:1.24.5-alpine AS builder

ADD . /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -a -o ./activepieces_telegram_bot .

FROM alpine:latest
COPY --from=builder /app/activepieces_telegram_bot ./
RUN chmod +x ./activepieces_telegram_bot

ENTRYPOINT [ "./activepieces_telegram_bot" ]
