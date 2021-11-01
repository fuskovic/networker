FROM golang:alpine3.10

LABEL maintainer="fuskovic"

RUN mkdir -p /go/src/github.com/fuskovic/networker
RUN mkdir -p /go/bin/

ADD . /go/src/github.com/fuskovic/networker

WORKDIR /go/src/github.com/fuskovic/networker

RUN apk update && \
    apk add git && \
    apk add curl && \
    apk add make && \
    apk add gcc && \
    apk add libc-dev && \
    apk add libpcap-dev && \
    apk add openssl

RUN go build -o /usr/local/bin/networker main.go

CMD [ "/bin/sh" ]