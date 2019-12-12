package compat030

import (
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/timonwong/prometheus-webhook-dingtalk/config"
)

type dingTalkProfilesValue map[string]*config.SecretURL

func asDingTalkProfiles(s kingpin.Settings) (target *dingTalkProfilesValue) {
	target = &dingTalkProfilesValue{}
	s.SetValue(target)
	return
}

func (s *dingTalkProfilesValue) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	profile, webhookURL := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	// Validate profile name
	if !config.TargetValidNameRE.MatchString(profile) {
		return fmt.Errorf("invalid profile name: %q", profile)
	}

	// Validate webhook url
	url, err := config.ParseURL(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook url: %s", err)
	}

	targetURL := config.SecretURL(*url)
	(*s)[profile] = &targetURL
	return nil
}

func (s *dingTalkProfilesValue) Get() interface{} {
	return *s
}

func (s *dingTalkProfilesValue) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *dingTalkProfilesValue) IsCumulative() bool {
	return true
}
