FROM golang:1.14-alpine AS build

#ARG LDFLAG
ARG BRANCH
ARG COMMIT
ARG VERSION
ARG BINARY

RUN mkdir /app
ADD . /app/
WORKDIR /app
#RUN CGO_ENABLED=0 go build ${LDFLAGS} .
RUN CGO_ENABLED=0 go build -ldflags "-w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}" .


FROM scratch
MAINTAINER Brian Hechinger <wonko@4amlunch.net>
VOLUME /etc/chremoas
COPY --from=build /app/${BINARY} /${BINARY}

ENTRYPOINT ["/${BINARY}", "--configuration_file", "/etc/chremoas/chremoas.yaml"]
