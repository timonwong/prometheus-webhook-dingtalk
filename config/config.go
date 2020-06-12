package config

import (
	"errors"
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
	DefaultTarget        = Target{}
	DefaultTargetMessage = TargetMessage{
		Title: `{{ template "ding.link.title" . }}`,
		Text:  `{{ template "ding.link.content" . }}`,
	}

	TargetValidNameRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9\-_]*$`)
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
	NoBuiltinTemplate bool              `yaml:"no_builtin_template"`
	Template          string            `yaml:"template,omitempty"`
	Templates         []string          `yaml:"templates,omitempty"`
	DefaultMessage    *TargetMessage    `yaml:"default_message,omitempty"`
	Timeout           time.Duration     `yaml:"timeout"`
	Targets           map[string]Target `yaml:"targets"`
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
		if !TargetValidNameRE.MatchString(name) {
			return fmt.Errorf("invalid target name: %q", name)
		}
	}

	if c.Template != "" {
		c.Templates = append(c.Templates, c.Template)
	}

	return nil
}

func (c *Config) String() string {
	b, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("<error creating config string: %s>", err)
	}
	return string(b)
}

func (c *Config) GetDefaultMessage() TargetMessage {
	if c.DefaultMessage != nil {
		return *c.DefaultMessage
	}
	return DefaultTargetMessage
}

type Target struct {
	URL     *SecretURL     `yaml:"url,omitempty"`
	Secret  Secret         `yaml:"secret,omitempty"`
	Mention *TargetMention `yaml:"mention,omitempty"`
	Message *TargetMessage `yaml:"message,omitempty"`
}

func (c *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultTarget
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain Target
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	if c.URL == nil {
		return errors.New("url cannot be empty")
	}

	return nil
}

type TargetMention struct {
	All     bool     `yaml:"all,omitempty"`
	Mobiles []string `yaml:"mobiles,omitempty"`
}

type TargetMessage struct {
	Title string `yaml:"title"`
	Text  string `yaml:"text"`
}

func (c *TargetMessage) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*c = DefaultTargetMessage
	// We want to set c to the defaults and then overwrite it with the input.
	// To make unmarshal fill the plain data struct rather than calling UnmarshalYAML
	// again, we have to hide it using a type indirection.
	type plain TargetMessage
	if err := unmarshal((*plain)(c)); err != nil {
		return err
	}

	return nil
}
