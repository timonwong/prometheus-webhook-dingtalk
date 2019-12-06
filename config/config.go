package config

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	defaultConfig = Config{
		Timeout: 5 * time.Second,
	}
	defaultTarget = Target{
		Message: defaultTargetMessage,
	}
	defaultTargetMessage = TargetMessage{
		Title: `{{ template "ding.link.title" . }}`,
		Text:  `{{ template "ding.link.content" . }}`,
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
	// never called. We thus have to set the defaultConfig at the entry
	// point as well.
	*cfg = defaultConfig
	err = yaml.UnmarshalStrict(content, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

type Config struct {
	Template  string            `yaml:"template"`
	Templates []string          `yaml:"templates"`
	Timeout   time.Duration     `yaml:"timeout"`
	Targets   map[string]Target `yaml:"targets"`
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = defaultConfig
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

	if c.Template != "" {
		c.Templates = append(c.Templates, c.Template)
	}

	return nil
}

type Target struct {
	URL     string         `yaml:"url"`
	Secret  string         `yaml:"secret"`
	Mention *TargetMention `yaml:"mention"`
	Message TargetMessage  `yaml:"message"`
}

func (c *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = defaultTarget
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Target
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}

type TargetMention struct {
	All     bool     `yaml:"all"`
	Mobiles []string `yaml:"mobiles"`
}

type TargetMessage struct {
	Title string `yaml:"title"`
	Text  string `yaml:"text"`
}

func (c *TargetMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = defaultTargetMessage
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain TargetMessage
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}
