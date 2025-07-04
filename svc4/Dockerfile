FROM golang:1.23.0 AS builder
WORKDIR /app
COPY . .
RUN go env -w  GOPROXY=https://goproxy.cn,direct  && \
    go mod download && \
    go build -a -o svc4_bin ./svc4/.

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/svc4_bin /app/svc4
CMD ["/app/svc4"]