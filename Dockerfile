FROM        quay.io/prometheus/busybox:latest
MAINTAINER  Timon Wong <timon86.wang@gmail.com>

COPY prometheus-webhook-dingtalk  /bin/prometheus-webhook-dingtalk
COPY template/default.tmpl        /usr/share/prometheus-webhook-dingtalk/template/default.tmpl

EXPOSE      8060
ENTRYPOINT  [ "/bin/prometheus-webhook-dingtalk" ]
