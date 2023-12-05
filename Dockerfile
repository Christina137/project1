FROM golang:1.21


MAINTAINER "chris"

WORKDIR /opt
ADD .  /opt

RUN apt-get update && \
    apt-get install -y ffmpeg

RUN go build -o main ./src/main.go

EXPOSE 8080

CMD ["/opt/main"]






