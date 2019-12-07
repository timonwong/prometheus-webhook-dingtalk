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

### 如何恢复为老版的样式 (v1.0.0 之前)

一、准备好默认的文件, 参见 `contrib/templates/legacy/default.tmpl`

二、编辑配置文件, 修改 `templates` 一节, 如下所示

```yaml
templates:
  # Docker 的话可以直接这样写
  - /etc/prometheus-webhook-dingtalk/templates/legacy/template.tmpl
```
