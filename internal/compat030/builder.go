package compat030

import (
	"time"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
)

type Builder struct {
	dingTalkProfiles *dingTalkProfilesValue
	requestTimeout   *time.Duration
	templateFile     *string
	isCompat         bool
}

func NewBuilder(a *kingpin.Application) *Builder {
	b := &Builder{}
	action := func(ctx *kingpin.ParseContext) error {
		b.isCompat = true
		return nil
	}
	b.dingTalkProfiles = dingTalkProfiles(a.Flag("ding.profile", "").Action(action).Hidden())
	b.requestTimeout = a.Flag("ding.timeout", "").Hidden().Action(action).Duration()
	b.templateFile = a.Flag("template.file", "").Hidden().Action(action).String()
	return b
}

func (b *Builder) IsCompatibleMode() bool {
	return b.isCompat
}

func (b *Builder) BuildConfig() (*config.Config, error) {
	conf := config.DefaultConfig
	if *b.requestTimeout != 0 {
		conf.Timeout = *b.requestTimeout
	}

	if *b.templateFile != "" {
		conf.Templates = []string{*b.templateFile}
	}

	conf.Targets = make(map[string]config.Target)
	for name, targetURL := range *b.dingTalkProfiles {
		url, err := config.ParseURL(targetURL)
		if err != nil {
			return nil, err
		}

		secretURL := config.SecretURL(*url)
		targetConfig := config.DefaultTarget
		targetConfig.URL = &secretURL
		conf.Targets[name] = targetConfig
	}

	return &conf, nil
}
