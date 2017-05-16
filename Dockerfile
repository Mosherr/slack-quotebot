FROM alpine:latest

MAINTAINER Daniel Scherr <danscherr@gmail.com>

WORKDIR "/opt"

ADD .docker_build/quotebot /opt/bin/quotebot
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/quotebot"]

