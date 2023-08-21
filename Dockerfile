FROM ubuntu:latest

WORKDIR /usr/local/bin

COPY ./blog-backend-service .

ENV VIPER_CONFIG=/usr/local/conf/config.yaml LOG_ROOT_DIR=/data/log DB_HOST=172.17.0.2

EXPOSE 8080

ENTRYPOINT ["blog-backend-service"]
