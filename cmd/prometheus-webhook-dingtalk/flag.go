package main

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type dingTalkProfilesValue map[string]string

func DingTalkProfiles(s kingpin.Settings) (target *dingTalkProfilesValue) {
	target = &dingTalkProfilesValue{}
	s.SetValue(target)
	return
}

func (s *dingTalkProfilesValue) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	profile, webhookURL := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	if profile == "" {
		return errors.New("profile part cannot be empty")
	}
	if webhookURL == "" {
		return errors.New("webhook-url part cannot be emtpy")
	}

	(*s)[profile] = webhookURL
	return nil
}

func (s *dingTalkProfilesValue) Get() interface{} {
	return (map[string]string)(*s)
}

func (s *dingTalkProfilesValue) String() string {
	return fmt.Sprintf("%s", map[string]string(*s))
}

func (s *dingTalkProfilesValue) IsCumulative() bool {
	return true
}
