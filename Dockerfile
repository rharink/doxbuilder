FROM ubuntu:14.04

MAINTAINER Robert den Harink <robert@rauwekost.com>

ENV PATH /usr/local/go/bin:/go/bin:/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin:/bin:/sbin
ENV GOPATH /go
ENV GOROOT /usr/local/go

RUN apt-get update
RUN apt-get install -y curl git

# Install Golang
RUN curl -O https://storage.googleapis.com/golang/go1.4.1.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.4.1.linux-amd64.tar.gz
RUN rm -f go1.4.1.linux-amd64.tar.gz

# Instal doxbuilder
RUN go get github.com/rauwekost/doxbuilder
RUN go build github.com/rauwekost/doxbuilder

# Install libreoffice
RUN apt-get install -y libreoffice

# Make required dirs
RUN mkdir -p /tmp

ADD configuration.yml .

EXPOSE 3000

CMD doxbuilder -c configuration.yml -p 3000
