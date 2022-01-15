FROM golang:1.17.6 AS stage

WORKDIR /go/src/

RUN apt-get -y install git

RUN git clone https://github.com/tee09/httpserver.git

RUN cd httpserver;\
    go build -o .

FROM alpine

COPY --from=stage /go/src/httpserver/httpserver /httpserver

CMD ["./httpserver"]
