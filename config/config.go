package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	DefaultConfig = Config{
		Timeout: 5 * time.Second,
	}

	targetValidNameRe = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*$`)
)

func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	// If the entire config body is empty the UnmarshalYAML method is
	// never called. We thus have to set the DefaultConfig at the entry
	// point as well.
	*cfg = DefaultConfig
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

type Config struct {
	Template string            `yaml:"template"` // Customized template file (see template/default.tmpl for example)
	Timeout  time.Duration     `yaml:"timeout"`  // Timeout for invoking DingTalk webhook
	Targets  map[string]Target `yaml:"targets"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultConfig
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Config
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	for name := range c.Targets {
		if !targetValidNameRe.MatchString(name) {
			return fmt.Errorf("invalid target name: %s", name)
		}
	}

	return nil
}

type Target struct {
	URL     string         `yaml:"url"`
	Secret  string         `yaml:"secret"`
	Mention *TargetMention `yaml:"mention"`
}

type TargetMention struct {
	All     bool     `yaml:"all"`
	Mobiles []string `yaml:"mobiles"`
}
