FROM alpine:latest

MAINTAINER Daniel Scherr <danscherr@gmail.com>

COPY . /app/src/github.com/mosherr/quotebot
WORKDIR /app/src/github.com/mosherr/quotebot

ENV HOME /app
ENV GOVERSION=1.8
ENV GOROOT $HOME/.go/$GOVERSION/go
ENV GOPATH $HOME
ENV PATH $PATH:$HOME/bin:$GOROOT/bin:$GOPATH/bin

RUN mkdir -p $HOME/.go/$GOVERSION
RUN cd $HOME/.go/$GOVERSION; curl -s https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz | tar zxf -
RUN go install -v github.com/mosherr/quotebot