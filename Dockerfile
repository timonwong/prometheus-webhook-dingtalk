ARG ARCH="amd64"
FROM gcr.io/distroless/static:debug-${ARCH} as etc

COPY config.example.yml                 /etc/prometheus-webhook-dingtalk/config.yml
COPY contrib                            /etc/prometheus-webhook-dingtalk/
COPY template/default.tmpl              /etc/prometheus-webhook-dingtalk/templates/default.tmpl

RUN ["/busybox/sh", "-c", "mkdir -p /prometheus-webhook-dingtalk"]
RUN ["/busybox/sh", "-c", "chown -R nobody:nobody /etc/prometheus-webhook-dingtalk /prometheus-webhook-dingtalk"]

ARG ARCH="amd64"
FROM gcr.io/distroless/static:nonroot-${ARCH}
LABEL maintainer="Timon Wong <timon86.wang@gmail.com>"

ARG ARCH="amd64"
ARG OS="linux"
COPY .build/${OS}-${ARCH}/prometheus-webhook-dingtalk   /bin/prometheus-webhook-dingtalk
COPY --from=etc /etc/prometheus-webhook-dingtalk       /etc/prometheus-webhook-dingtalk
COPY --from=etc /prometheus-webhook-dingtalk           /prometheus-webhook-dingtalk

USER       nobody
EXPOSE     8060
VOLUME     [ "/prometheus-webhook-dingtalk" ]
WORKDIR    /prometheus-webhook-dingtalk
ENTRYPOINT [ "/bin/prometheus-webhook-dingtalk" ]
CMD        [ "--config.file=/etc/prometheus-webhook-dingtalk/config.yml" ]
