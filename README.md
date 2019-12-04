# prometheus-webhook-dingtalk

Generating [DingTalk] notification from [Prometheus] [AlertManager] WebHooks.

## Building and running

### Build

```bash
make
```

### Running

```bash
./prometheus-webhook-dingtalk <flags>
```

## Usage

```
usage: prometheus-webhook-dingtalk [<flags>]

Flags:
  -h, --help               Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=:8060
                           The address to listen on for web interface.
      --config.file=config.yml
                           Path to the configuration file.
      --log.level=info     Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt  Output format of log messages. One of: [logfmt, json]
      --version            Show application version.

```

## Configuration

Example configuration

```yaml
# timeout: 5s
# template: template/default.tmpl
targets:
  webhook1:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
  webhook2:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
```

## Using Docker

You can deploy this tool using the Docker image from following registry:

* [DockerHub]\: [timonwong/prometheus-webhook-dingtalk](https://registry.hub.docker.com/u/timonwong/prometheus-webhook-dingtalk/)
* [Quay.io]\: [timonwong/prometheus-webhook-dingtalk](https://quay.io/repository/timonwong/prometheus-webhook-dingtalk)

[Prometheus]: https://prometheus.io
[AlertManager]: https://github.com/prometheus/alertmanager
[DingTalk]: https://www.dingtalk.com
[DockerHub]: https://hub.docker.com
[Quay.io]: https://quay.io
