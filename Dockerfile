FROM golang
MAINTAINER Octoblu, Inc. <docker@octoblu.com>
EXPOSE 80

ADD https://raw.githubusercontent.com/pote/gpm/v1.3.2/bin/gpm /go/bin/
RUN chmod +x /go/bin/gpm

COPY Godeps /go/src/github.com/octoblu/go-meshblu-http-server/
WORKDIR /go/src/github.com/octoblu/go-meshblu-http-server
RUN gpm install

COPY . /go/src/github.com/octoblu/go-meshblu-http-server

RUN env CGO_ENABLED=0 go build -a -ldflags '-s' .

CMD ["./go-meshblu-http-server"]
