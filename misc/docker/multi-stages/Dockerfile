FROM golang:1.23 AS build

RUN apt-get update
RUN apt-get install -y make

COPY . /src

WORKDIR /src

RUN make build

FROM scratch AS runtime

COPY --from=build /src/bin/libra /usr/local/bin/libra

CMD ["/usr/local/bin/libra"]