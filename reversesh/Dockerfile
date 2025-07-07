FROM golang:1.23.0 AS builder
WORKDIR /app
COPY . .
RUN go env -w  GOPROXY=https://goproxy.cn,direct  && \
    go mod download && \
    go build -a -o r_bin ./reversesh/.

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/r_bin /app/reversesh
CMD ["/app/reversesh"]