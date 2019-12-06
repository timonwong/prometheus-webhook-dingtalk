ARG ARCH="amd64"
ARG OS="linux"

FROM quay.io/prometheus/busybox-${OS}-${ARCH}:latest
LABEL maintainer="Timon Wong <timon86.wang@gmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/prometheus-webhook-dingtalk   /bin/prometheus-webhook-dingtalk
COPY config.example.yml                                 /etc/prometheus-webhook-dingtalk/config.yml
COPY contrib                                            /etc/prometheus-webhook-dingtalk/
COPY template/default.tmpl                              /etc/prometheus-webhook-dingtalk/templates/default.tmpl

RUN mkdir -p /prometheus-webhook-dingtalk && \
    chown -R nobody:nogroup /etc/prometheus-webhook-dingtalk /prometheus-webhook-dingtalk

USER       nobody
EXPOSE     8060
WORKDIR    /prometheus-webhook-dingtalk
ENTRYPOINT [ "/bin/prometheus-webhook-dingtalk" ]
CMD        [ "--config.file=/etc/prometheus-webhook-dingtalk/config.yml" ]
