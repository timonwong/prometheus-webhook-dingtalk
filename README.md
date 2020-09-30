# prometheus-webhook-dingtalk

[![Build Status](https://img.shields.io/circleci/build/github/timonwong/prometheus-webhook-dingtalk)](https://circleci.com/gh/timonwong/prometheus-webhook-dingtalk)
[![Go Report Card](https://goreportcard.com/badge/github.com/timonwong/prometheus-webhook-dingtalk)](https://goreportcard.com/report/github.com/timonwong/prometheus-webhook-dingtalk)
[![Docker Pulls](https://img.shields.io/docker/pulls/timonwong/prometheus-webhook-dingtalk)](https://hub.docker.com/r/timonwong/prometheus-webhook-dingtalk)

Generating [DingTalk] notification from [Prometheus] [AlertManager] WebHooks.

## Install

### Precompiled binaries

Precompiled binaries for released versions are available in [release page](https://github.com/timonwong/prometheus-webhook-dingtalk/releases):
It's always recommended to use latest stable version available.

### Docker

You can deploy this tool using the Docker image from following registry:

* DockerHub: [timonwong/prometheus-webhook-dingtalk](https://hub.docker.com/r/timonwong/prometheus-webhook-dingtalk)

### Compiling the binary

#### Prerequisites

1. [Go](https://golang.org/doc/install) (1.13 or greater is required)
2. [Nodejs](https://nodejs.org/)
3. [Yarn](https://yarnpkg.com/)

#### Build

Clone the repository and build manually:

```bash
make build
```

## Usage

```
usage: prometheus-webhook-dingtalk [<flags>]

Flags:
  -h, --help                    Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=:8060
                                The address to listen on for web interface.
      --web.enable-ui           Enable Web UI mounted on /ui path
      --web.enable-lifecycle    Enable reload via HTTP request.
      --config.file=config.yml  Path to the configuration file.
      --log.level=info          Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt       Output format of log messages. One of: [logfmt, json]
      --version                 Show application version.
```

For Kubernetes users, check out [./contrib/k8s](./contrib/k8s).

## Configuration

常见问题可以看看 [FAQ](./docs/FAQ_zh.md)

```yaml
## Request timeout
# timeout: 5s

## Customizable templates path
# templates:
#   - contrib/templates/legacy/template.tmpl

## You can also override default template using `default_message`
## The following example to use the 'legacy' template from v0.3.0
# default_message:
#   title: '{{ template "legacy.title" . }}'
#   text: '{{ template "legacy.content" . }}'

## Targets, previously was known as "profiles"
targets:
  webhook1:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
    # secret for signature
    secret: SEC000000000000000000000
  webhook2:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
  webhook_legacy:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
    # Customize template content
    message:
      # Use legacy template
      title: '{{ template "legacy.title" . }}'
      text: '{{ template "legacy.content" . }}'
  webhook_mention_all:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
    mention:
      all: true
  webhook_mention_users:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
    mention:
      mobiles: ['156xxxx8827', '189xxxx8325']
```

[Prometheus]: https://prometheus.io
[AlertManager]: https://github.com/prometheus/alertmanager
[DingTalk]: https://www.dingtalk.com
