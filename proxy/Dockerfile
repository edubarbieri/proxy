FROM golang:1.16-alpine3.13 AS build
WORKDIR /src
ADD . /src
RUN go build -o proxy


FROM alpine:3.13
WORKDIR /app

ENV REDIS_ADDRESS=localhost:6379
ENV REDIS_PASSWORD=
ENV DEBUG=true
ENV CONFIG_PATH=config.json

COPY --from=build /src/proxy /app/
ENTRYPOINT ./proxy