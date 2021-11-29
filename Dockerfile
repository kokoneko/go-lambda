FROM golang:1.17.3-alpine3.14

RUN mkdir /go/src/app

RUN apk update \ && apk add zip

RUN apk add --no-cache imagemagick bash pngcrush optipng=0.7.7-r0 \
    gcc \
    imagemagick-dev \
    libc-dev

WORKDIR /go/src/app

ADD . /go/src/app