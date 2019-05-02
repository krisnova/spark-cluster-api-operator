FROM golang:latest
MAINTAINER Kris Nova "kris@nivenly.com"
ADD . /go/src/github.com/kris-nova/spark-cluster-api-operator
WORKDIR /go/src/github.com/kris-nova/spark-cluster-api-operator
RUN make compile
CMD ["./bin/spark-cluster-api-operator"]