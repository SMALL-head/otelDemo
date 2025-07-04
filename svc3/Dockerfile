FROM golang:1.23.0 AS builder
WORKDIR /app
COPY . .
RUN go env -w  GOPROXY=https://goproxy.cn,direct  && \
    go mod download && \
    go build -a -o svc3_bin ./svc3/.

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/svc3_bin /app/svc3
CMD ["/app/svc3"]