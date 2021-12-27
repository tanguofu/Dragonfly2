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

package config

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mitchellh/mapstructure"
	testifyassert "github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestSchedulerConfig_Load(t *testing.T) {
	assert := testifyassert.New(t)

	config := &Config{
		DynConfig: &DynConfig{
			RefreshInterval: 5 * time.Minute,
			CDNDir:          "tmp",
		},
		Scheduler: &SchedulerConfig{
			Algorithm:       "default",
			WorkerNum:       10,
			BackSourceCount: 3,
			GC: &GCConfig{
				PeerGCInterval: 1 * time.Minute,
				PeerTTL:        10 * time.Minute,
				TaskGCInterval: 1 * time.Minute,
				TaskTTL:        10 * time.Minute,
			},
		},
		Server: &ServerConfig{
			IP:       "127.0.0.1",
			Host:     "foo",
			Port:     8002,
			CacheDir: "foo",
			LogDir:   "foo",
		},
		Manager: &ManagerConfig{
			Addr:               "127.0.0.1:65003",
			SchedulerClusterID: 1,
			KeepAlive: KeepAliveConfig{
				Interval: 1 * time.Second,
			},
		},
		Host: &HostConfig{
			IDC:         "foo",
			NetTopology: "baz",
			Location:    "bar",
		},
		Job: &JobConfig{
			GlobalWorkerNum:    1,
			SchedulerWorkerNum: 1,
			LocalWorkerNum:     5,
			Redis: &RedisConfig{
				Host:      "127.0.0.1",
				Port:      6379,
				Password:  "password",
				BrokerDB:  1,
				BackendDB: 2,
			},
		},
		Metrics: &MetricsConfig{
			Addr:           ":8000",
			EnablePeerHost: true,
		},
	}

	schedulerConfigYAML := &Config{}
	contentYAML, _ := os.ReadFile("./testdata/scheduler.yaml")
	var dataYAML map[string]interface{}
	if err := yaml.Unmarshal(contentYAML, &dataYAML); err != nil {
		t.Fatal(err)
	}

	if err := mapstructure.Decode(dataYAML, &schedulerConfigYAML); err != nil {
		t.Fatal(err)
	}
	assert.True(reflect.DeepEqual(config, schedulerConfigYAML))
}

func TestConvert(t *testing.T) {
	tests := []struct {
		name   string
		value  *Config
		expect func(t *testing.T, cfg *Config, err error)
	}{
		{
			name: "convert common config",
			value: &Config{
				Manager: &ManagerConfig{
					Addr: "127.0.0.1:65003",
				},
				Job: &JobConfig{
					Redis: &RedisConfig{
						Host: "",
					},
				},
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert := testifyassert.New(t)
				assert.Equal("127.0.0.1", cfg.Job.Redis.Host)
			},
		},
		{
			name: "convert config when host not empty",
			value: &Config{
				Manager: &ManagerConfig{
					Addr: "127.0.0.1:65003",
				},
				Job: &JobConfig{
					Redis: &RedisConfig{
						Host: "111.111.11.1",
					},
				},
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert := testifyassert.New(t)
				assert.Equal("111.111.11.1", cfg.Job.Redis.Host)
			},
		},
		{
			name: "convert config when manager addr is empty",
			value: &Config{
				Manager: &ManagerConfig{
					Addr: "",
				},
				Job: &JobConfig{
					Redis: &RedisConfig{
						Host: "111.111.11.1",
					},
				},
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert := testifyassert.New(t)
				assert.Equal("111.111.11.1", cfg.Job.Redis.Host)
			},
		},
		{
			name: "convert config when manager host is empty",
			value: &Config{
				Manager: &ManagerConfig{
					Addr: ":65003",
				},
				Job: &JobConfig{
					Redis: &RedisConfig{
						Host: "",
					},
				},
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert := testifyassert.New(t)
				assert.Equal("", cfg.Job.Redis.Host)
			},
		},
		{
			name: "convert config when manager host is localhost",
			value: &Config{
				Manager: &ManagerConfig{
					Addr: "localhost:65003",
				},
				Job: &JobConfig{
					Redis: &RedisConfig{
						Host: "",
					},
				},
			},
			expect: func(t *testing.T, cfg *Config, err error) {
				assert := testifyassert.New(t)
				assert.Equal("localhost", cfg.Job.Redis.Host)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.value.Convert()
			tc.expect(t, tc.value, err)
		})
	}
}
