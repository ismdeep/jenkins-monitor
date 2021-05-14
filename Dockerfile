FROM golang as builder
RUN mkdir -p /jenkins-monitor
ADD . /jenkins-monitor
WORKDIR /jenkins-monitor
RUN go build -o jenkins-monitor


FROM alpine
MAINTAINER L. Jiang <l.jiang.1024@gmail.com>

COPY requirements.txt /
COPY k8s-restart.py /
COPY k8s-restart.sh /
RUN chmod +x /k8s-restart.sh
COPY --from=builder /jenkins-monitor/jenkins-monitor /jenkins-monitor
RUN wget -q -O /etc/apk/keys/sgerrand.rsa.pub https://alpine-pkgs.sgerrand.com/sgerrand.rsa.pub
RUN wget https://github.com/sgerrand/alpine-pkg-glibc/releases/download/2.29-r0/glibc-2.29-r0.apk
RUN apk add glibc-2.29-r0.apk
RUN apk add python3
RUN apk add py3-pip
RUN pip3.8 install -r /requirements.txt

CMD ["/jenkins-monitor","-c", "/config.json"]