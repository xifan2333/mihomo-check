FROM debian:bookworm

ENV TZ=Asia/Shanghai

RUN apt-get update && apt-get install -y ca-certificates tzdata curl && \
    ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    rm -rf /var/cache/apt/*  && \
    mkdir -p /app

RUN ARCH=$(uname -m) && \
    if [ "$ARCH" = "x86_64" ]; then \
        FILE="BestSub_Linux_x86_64.tar.gz"; \
    elif [ "$ARCH" = "aarch64" ]; then \
        FILE="BestSub_Linux_aarch64.tar.gz"; \
    elif [ "$ARCH" = "armv7l" ]; then \
        FILE="BestSub_Linux_armv7.tar.gz"; \
    elif [ "$ARCH" = "i386" ]; then \
        FILE="BestSub_Linux_i386.tar.gz"; \
    else \
        echo "Unsupported architecture: $ARCH"; exit 1; \
    fi && \
    curl -L -o /tmp/${FILE} https://github.com/bestruirui/BestSub/releases/latest/download/${FILE} && \
    tar -xzf /tmp/${FILE} -C /app && \
    rm /tmp/${FILE}

CMD /app/BestSub