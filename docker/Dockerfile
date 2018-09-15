FROM debian:latest

RUN mkdir -p /tmp/data

COPY ./data/full/data.zip /tmp/data/data.zip
COPY ./data/full/options.txt /tmp/data/options.txt

COPY hlcup /tmp/hlcup

EXPOSE 80

CMD ["/tmp/hlcup"]