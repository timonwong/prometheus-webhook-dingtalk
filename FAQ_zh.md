## 常见问题

### GeneratorURL 不对 / 模板中链接的问题

相关 Issue:

- https://github.com/timonwong/prometheus-webhook-dingtalk/issues/27
- https://github.com/timonwong/prometheus-webhook-dingtalk/issues/20

请配置 prometheus 的 `--web.external-url`:

```
      --web.external-url=<URL>   The URL under which Prometheus is externally reachable (for example, if Prometheus is served via a reverse proxy). Used for generating relative and absolute links
                                 back to Prometheus itself. If the URL has a path portion, it will be used to prefix all HTTP endpoints served by Prometheus. If omitted, relevant URL components
                                 will be derived automatically.
```
