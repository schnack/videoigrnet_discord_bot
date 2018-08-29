FROM golang:1.11-alpine

RUN apk add --update git gcc libc-dev && \
    ln -s /usr/share/easy-rsa/easyrsa /usr/local/bin && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /go/src/app/db

ENTRYPOINT ["/go/bin/app"]