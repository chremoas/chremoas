FROM arm32v6/alpine
LABEL maintainer="maurer.it@gmail.com"

RUN apk update && apk upgrade

ADD ./chremoas /
WORKDIR /

RUN mkdir /etc/chremoas
VOLUME /config

RUN rm -rf /var/cache/apk/*

ENV MICRO_REGISTRY_ADDRESS chremoas-consul:8500

CMD [""]
ENTRYPOINT ["./chremoas", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
