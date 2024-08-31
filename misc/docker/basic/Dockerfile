FROM golang:1.23

RUN apt-get update
RUN apt-get install -y make

COPY . /src

WORKDIR /src

RUN make build

RUN cp ./bin/libra /usr/local/bin/libra

CMD ["/usr/local/bin/libra"]