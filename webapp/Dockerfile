FROM golang:1.23.7-alpine AS builder

WORKDIR /app

COPY . .

RUN GOPROXY="https://goproxy.cn" CGO_ENABLED=0 go build -o /app/webapp .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/webapp .
COPY static ./static/
COPY html ./html/
COPY entrypoint.sh ./entrypoint.sh

RUN chmod +x ./entrypoint.sh
RUN chmod +x ./webapp
ENTRYPOINT [ "/app/entrypoint.sh" ]

EXPOSE 8080