## 1.0.0 / unreleased

**BREAKING CHANGES**

Now instead of configuration via command line arguments, the YAML configuration file is used for better flexibility.

The example configuration file looks like:

```yaml
## Request timeout
# timeout: 5s

## Customizable template file path
## In docker, by default the current working directory is set to /prometheus-webhook-dingtalk
## However it's recommended to use absolute path whenever possible
# template: template/default.tmpl
targets:
  webhook1:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
    # secret for signature
    secret: SEC000000000000000000000
  webhook2:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
  webhook_mention_all:
    url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxxxxx
    mention:
      all: true
  webhook_mention_users:
    mention:
      mobiles: ['156xxxx8827', '189xxxx8325']
```

- [ENHANCEMENT] Add various template functions from [sprig](http://masterminds.github.io/sprig/) #47
- [ENHANCEMENT] Add signature support due to the new security enforcement requirement of dingtalk #49
