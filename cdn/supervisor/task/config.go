/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package task

import (
	"fmt"
	"time"
)

type Config struct {
	// GCInitialDelay is the delay time from the start to the first GC execution.
	// default: 6s
	GCInitialDelay time.Duration `yaml:"gcInitialDelay" mapstructure:"gcInitialDelay"`

	// GCMetaInterval is the interval time to execute GC meta.
	// default: 2min
	GCMetaInterval time.Duration `yaml:"gcMetaInterval" mapstructure:"gcMetaInterval"`

	// ExpireTime when a task is not accessed within the ExpireTime,
	// and it will be treated to be expired.
	// default: 3min
	ExpireTime time.Duration `yaml:"taskExpireTime" mapstructure:"taskExpireTime"`

	// FailAccessInterval is the interval time after failed to access the URL.
	// unit: minutes
	// default: 30
	FailAccessInterval time.Duration `yaml:"failAccessInterval" mapstructure:"failAccessInterval"`
}

func DefaultConfig() Config {
	config := Config{}
	return config.applyDefaults()
}

func (c Config) applyDefaults() Config {
	if c.GCInitialDelay == 0 {
		c.GCInitialDelay = DefaultGCInitialDelay
	}
	if c.GCMetaInterval == 0 {
		c.GCMetaInterval = DefaultGCMetaInterval
	}
	if c.ExpireTime == 0 {
		c.ExpireTime = DefaultExpireTime
	}
	if c.FailAccessInterval == 0 {
		c.FailAccessInterval = DefaultFailAccessInterval
	}
	return c
}

func (c Config) Validate() []error {
	var errors []error
	if c.GCInitialDelay < 0 {
		errors = append(errors, fmt.Errorf("task GCInitialDelay %d can't be a negative number", c.GCInitialDelay))
	}
	if c.GCMetaInterval <= 0 {
		errors = append(errors, fmt.Errorf("task GCMetaInterval must be greater than 0, but is: %d", c.GCMetaInterval))
	}
	if c.ExpireTime <= 0 {
		errors = append(errors, fmt.Errorf("task ExpireTime must be greater than 0, but is: %d", c.ExpireTime))
	}
	if c.FailAccessInterval <= 0 {
		errors = append(errors, fmt.Errorf("task FailAccessInterval must be greater than 0, but is: %d", c.FailAccessInterval))
	}
	return errors
}

const (
	// DefaultFailAccessInterval is the interval time after failed to access the URL.
	DefaultFailAccessInterval = 3 * time.Minute
)

// gc config
const (
	// DefaultGCInitialDelay is the delay time from the start to the first GC execution.
	DefaultGCInitialDelay = 6 * time.Second

	// DefaultGCMetaInterval is the interval time to execute the GC meta.
	DefaultGCMetaInterval = 2 * time.Minute

	// DefaultExpireTime when a task is not accessed within the ExpireTime,
	// and it will be treated to be expired.
	DefaultExpireTime = 30 * time.Minute
)
