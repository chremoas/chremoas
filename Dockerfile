FROM alpine
MAINTAINER Brian Hechinger <wonko@4amlunch.net>

ADD chremoas-linux-amd64 chremoas
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
VOLUME /etc/chremoas

RUN rm -rf /var/cache/apk/*

ENTRYPOINT ["/chremoas", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
