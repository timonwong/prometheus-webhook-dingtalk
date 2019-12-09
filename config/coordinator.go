// Copyright 2019 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Coordinator coordinates configurations beyond the lifetime of a
// single configuration.
type Coordinator struct {
	configFilePath string
	logger         log.Logger

	// Protects config and subscribers
	mutex        sync.Mutex
	config       *Config
	frozenConfig *Config
	subscribers  []func(*Config) error
}

// NewCoordinator returns a new coordinator with the given configuration file
// path. It does not yet load the configuration from file. This is done in
// `Reload()`.
func NewCoordinator(configFilePath string, frozenConfig *Config, l log.Logger) *Coordinator {
	c := &Coordinator{
		configFilePath: configFilePath,
		logger:         l,
		frozenConfig:   frozenConfig,
	}

	return c
}

// Subscribe subscribes the given Subscribers to configuration changes.
func (c *Coordinator) Subscribe(ss ...func(*Config) error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.subscribers = append(c.subscribers, ss...)
}

func (c *Coordinator) notifySubscribers() error {
	for _, s := range c.subscribers {
		if err := s(c.config); err != nil {
			return err
		}
	}

	return nil
}

// loadFromFile triggers a configuration load, discarding the old configuration.
func (c *Coordinator) loadFromFile() error {
	conf, err := LoadFile(c.configFilePath)
	if err != nil {
		return err
	}

	c.config = conf
	return nil
}

// Reload triggers a configuration reload from file and notifies all
// configuration change subscribers.
func (c *Coordinator) Reload() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	logger := log.With(c.logger, "file", c.configFilePath)
	if c.frozenConfig != nil {
		logger = c.logger
		c.config = c.frozenConfig
	} else {
		level.Info(logger).Log("msg", "Loading configuration file")
		if err := c.loadFromFile(); err != nil {
			level.Error(logger).Log(
				"msg", "Loading configuration file failed",
				"err", err,
			)
			return err
		}
		level.Info(logger).Log("msg", "Completed loading of configuration file")
	}

	if err := c.notifySubscribers(); err != nil {
		logger.Log("msg", "one or more config change subscribers failed to apply new config", "err", err)
		return err
	}

	return nil
}
