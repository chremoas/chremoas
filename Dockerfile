FROM scratch
MAINTAINER Brian Hechinger <wonko@4amlunch.net>

ADD chremoas-linux-amd64 chremoas
ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
VOLUME /etc/chremoas

ENV MICRO_REGISTRY_ADDRESS chremoas-consul:8500

ENTRYPOINT ["/chremoas", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
