FROM ubuntu:latest

# необходимые пакеты
RUN apt-get update && apt-get install -y \
    clang \
    gcc \
    g++ \
    python3 \
    && rm -rf /var/lib/apt/lists/*


WORKDIR /usr/src/app


COPY run_code.sh ./


RUN chmod +x run_code.sh

CMD ["./run_code.sh"]