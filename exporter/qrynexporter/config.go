// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qrynexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	defaultDSN = "tcp://127.0.0.1:9000/cloki"
)

// Config defines configuration for logging exporter.
type Config struct {
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`
	// QueueSettings is a subset of exporterhelper.QueueSettings,
	// because only QueueSize is user-settable.
	QueueSettings QueueSettings `mapstructure:"sending_queue"`

	ClusteredClickhouse bool `mapstructure:"clustered_clickhouse"`

	// DSN is the ClickHouse server Data Source Name.
	// For tcp protocol reference: [ClickHouse/clickhouse-go#dsn](https://github.com/ClickHouse/clickhouse-go#dsn).
	// For http protocol reference: [mailru/go-clickhouse/#dsn](https://github.com/mailru/go-clickhouse/#dsn).
	DSN string `mapstructure:"dsn"`

	// Logs is used to configure the log data.
	Logs LogsConfig `mapstructure:"logs"`
	// Metrics is used to configure the metric data.
	Metrics MetricsConfig `mapstructure:"metrics"`
}

// QueueSettings is a subset of exporterhelper.QueueSettings.
type QueueSettings struct {
	// QueueSize set the length of the sending queue
	QueueSize int `mapstructure:"queue_size"`
}

// LogsConfig holds the configuration for log data.
type LogsConfig struct {
	// AttritubeLabels is the string representing attribute labels.
	AttritubeLabels string `mapstructure:"attritube_labels"`
	// ResourceLabels is the string representing resource labels.
	ResourceLabels string `mapstructure:"resource_labels"`
	// Format is the string representing the format.
	Format string `mapstructure:"format"`
}

// MetricsConfig holds the configuration for metric data.
type MetricsConfig struct {
	// Namespace is the prefix attached to each exported metric name.
	// See: https://prometheus.io/docs/practices/naming/#metric-names
	Namespace string `mapstructure:"namespace"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	return nil
}

func (cfg *Config) enforcedQueueSettings() exporterhelper.QueueSettings {
	return exporterhelper.QueueSettings{
		Enabled:      true,
		NumConsumers: 1,
		QueueSize:    cfg.QueueSettings.QueueSize,
	}
}
