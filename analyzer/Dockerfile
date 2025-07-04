FROM golang:1.23.0 AS builder
WORKDIR /app
COPY . .
RUN go env -w  GOPROXY=https://goproxy.cn,direct  && \
    go mod download && \
    go build -a -o analyzer ./analyzer/.

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/analyzer/analyzer /app/analyzer
COPY analyzer/application.yaml application.yaml
CMD ["/app/analyzer"]