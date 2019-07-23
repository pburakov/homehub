FROM openjdk:11-jre

MAINTAINER Paul Burakov <pburakov@gmail.com>

COPY target/dist .

EXPOSE 31321/tcp

ENTRYPOINT sh bin/server
