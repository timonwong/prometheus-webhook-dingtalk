## 1.4.0 / 2019-12-11

- [FEATURE/ENHANCEMENT] Allow override global default message template in config file. #76
- [CHANGE] UI: Use CommonMark instead of GFM.
- [ENHANCEMENT] Filter hidden flags in flags UI. #73
- [ENHANCEMENT] Add goroutine info in runtime UI. #74

## 1.3.0 / 2019-12-09

- [FEATURE/ENHANCEMENT] Improved compatibility: Now the following v0.3.0 command line flags are supported as well: `--ding.profile`, `--ding.timeout` and `--template.file`. #72
- [FEATURE/ENHANCEMENT] Add support to reload through API `/-/reload` (Disabled by default, can be enabled via the `--web.enable-lifecycle` flag). #70
- [FEATURE/ENHANCEMENT] Add ready and health check API endpoint: `/-/healthy` and `/-/ready`. #71
- [CHANGE] (Backward compatible) Rename `--web.ui-enabled` to `--web.enable-ui`.

## 1.2.2 / 2019-12-09

- [FIX] Fix excessive rendering requests while in web UI preview. #65

## 1.2.1 / 2019-12-08

- [FIX] Fix default template (which misleads users).

## 1.2.0 / 2019-12-08

**NOTE** For security reason, the Web UI is disabled by default. In order to enable it, pass the `--web.ui-enabled` flag
when program starts.

- [FEATURE] Add web UI for playground (preview & validate your templates, etc). #62, #63
- [ENHANCEMENT] Add support to configuration reload through SIGHUP signal. #60
- [ENHANCEMENT] Validate target incoming webhook url in case of typos. #61

## 1.1.0 / 2019-12-06

- [FEATURE] Allow template customization for target individually. #58
- [ENHANCEMENT] Change default template.

## 1.0.0 / 2019-12-05

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

## Targets, previously was known as "profiles"
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
- [FEATURE] Add user mention support (all or specific mobiles) to dingtalk notification #54
