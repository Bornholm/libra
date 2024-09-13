FROM golang:1.23 AS build

COPY go.mod go.sum /src/

RUN cd /src && go mod download

RUN apt-get update
RUN apt-get install -y make

COPY . /src

WORKDIR /src

RUN make build

FROM scratch AS runtime

COPY --from=build /src/bin/libra /usr/local/bin/libra

CMD ["/usr/local/bin/libra"]