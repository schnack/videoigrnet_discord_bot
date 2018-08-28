FROM golang:1.11

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /go/src/app/db

ENTRYPOINT ["/go/bin/app"]