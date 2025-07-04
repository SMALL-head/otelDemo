FROM golang:1.23.0 AS builder
WORKDIR /app
COPY . .
RUN go env -w  GOPROXY=https://goproxy.cn,direct  && \
    go mod download && \
    go build -a -o analyzer ./analyzer/.

FROM ubuntu:22.04
COPY --from=builder /app/analyzer /app/analyzer
WORKDIR /app
COPY analyzer/application.yaml application.yaml
CMD ["/app/analyzer"]