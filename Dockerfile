FROM golang:1.21

#作者信息
MAINTAINER "chris"


#工作目录
WORKDIR /opt
ADD .  /opt


# 安装ffmpeg
RUN apt-get update && \
    apt-get install -y ffmpeg

#在Docker工作目录下执行命令
RUN go build -o main ./src/main.go


#暴露端口
EXPOSE 8080

#执行项目的命令
CMD ["/opt/main"]






