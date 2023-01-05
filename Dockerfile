FROM golang:1.19.4 AS build-env
ENV GO111MODULE=on
WORKDIR /app/ipsec_checker
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ipsec_checker

FROM ubuntu:22.04 AS dev-build
ENV TZ="Asia/Almaty"
ENV DEBIAN_FRONTEND=noninteractive

RUN apt -y update
RUN apt -y install strongswan vim ca-certificates && update-ca-certificates

# BACKUP script rotate.
WORKDIR /usr/local/bin
COPY --from=build-env /app/ipsec_checker/ipsec_checker .
USER root
RUN chmod +x ipsec_checker
COPY --from=build-env /app/ipsec_checker/.env.example /usr/local/etc/ipsec_checker/.env
RUN ls -lah /usr/local/etc/
RUN ls -lah /usr/local/etc/ipsec_checker/

# Add files
ADD docker/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT /entrypoint.sh