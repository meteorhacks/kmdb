FROM google/golang

RUN mkdir -p /gopath/src/github.com/meteorhacks/kmdb
WORKDIR /gopath/src/github.com/meteorhacks/kmdb
ADD . /gopath/src/github.com/meteorhacks/kmdb
RUN go get github.com/meteorhacks/kmdb
CMD ["kmdb", "-config", "/etc/kmdb.json"]
VOLUME ["/data", "/etc/kmdb.json"]
