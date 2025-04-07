FROM golang:1.23 AS builder
COPY . /app
WORKDIR /app
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o main ./receiver


FROM ubuntu:latest
COPY --from=builder /app/main /app/main
WORKDIR /app
CMD ["/app/main"]