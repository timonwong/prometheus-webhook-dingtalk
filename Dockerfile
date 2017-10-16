FROM        quay.io/prometheus/busybox:latest
MAINTAINER  Timon Wong <timon86.wang@gmail.com>

COPY prometheus-webhook-dingtalk /bin/prometheus-webhook-dingtalk

EXPOSE      9117
ENTRYPOINT  [ "/bin/prometheus-webhook-dingtalk" ]
