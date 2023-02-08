FROM ubuntu:latest

WORKDIR /usr/local/bin

COPY ./blog-backend-service .

COPY ./config.yaml .

ENV VIPER_CONFIG=/usr/local/bin/config.yaml LOG_ROOT_DIR=/data/log DB_HOST=172.17.0.2

EXPOSE 8080

ENTRYPOINT ["blog-backend-service"]
