FROM golang:bookworm AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main .

FROM debian:bookworm
ENV TZ=Asia/Shanghai
RUN apt-get update && apt-get install -y ca-certificates tzdata && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    rm -rf /var/cache/apt/*
COPY --from=builder /app/main /app/main
CMD /app/main