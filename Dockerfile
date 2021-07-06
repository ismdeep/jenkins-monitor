FROM golang as builder
RUN mkdir -p /jenkins-monitor
ADD . /jenkins-monitor
WORKDIR /jenkins-monitor
RUN go build -o jenkins-monitor


FROM frolvlad/alpine-glibc
MAINTAINER L. Jiang <l.jiang.1024@gmail.com>

ADD requirements.txt k8s-restart.py k8s-restart.sh /
RUN chmod +x /k8s-restart.sh
COPY --from=builder /jenkins-monitor/jenkins-monitor /jenkins-monitor
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk add bash
RUN apk add python3
RUN apk add py3-pip
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone
RUN pip3.8 install -r /requirements.txt

CMD ["/jenkins-monitor","-c", "/config.json"]