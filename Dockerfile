FROM golang:1.16.5-alpine3.13

RUN mkdir /go/src/app

RUN apk update \ && apk add zip

WORKDIR /go/src/app

ADD . /go/src/app