FROM ubuntu:17.10

RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    git-core \
    gcc

RUN curl -s https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz | tar -v -C /usr/local -xz
RUN mkdir -p /go

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH /usr/local/go/bin:/go/bin:/usr/local/bin:$PATH
WORKDIR $GOPATH

RUN go get -u github.com/EUDAT-GEF/BridgIt; exit 0
WORKDIR $GOPATH/src/github.com/EUDAT-GEF/BridgIt
RUN mkdir -p tmp \
    mkdir -p build
RUN go build
RUN cp -r BridgIt tmp/

CMD ["cp", "-r", "tmp/.", "build"]
