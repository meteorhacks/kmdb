FROM google/golang

RUN mkdir -p /gopath/src/github.com/meteorhacks/kmdb
ADD . /gopath/src/github.com/meteorhacks/kmdb
RUN cd /gopath/src/github.com/meteorhacks/kmdb && go get ./...
CMD ["kmdb", "-config", "/etc/kmdb.json"]
VOLUME ["/data", "/etc/kmdb.json"]
