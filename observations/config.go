package observations

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/ipfs/ipfs-cluster/config"
	ma "github.com/multiformats/go-multiaddr"
)

const metricsConfigKey = "metrics"
const tracingConfigKey = "tracing"
const envConfigKey = "cluster_observations"

// Default values for this Config.
const (
	DefaultEnableStats            = false
	DefaultPrometheusEndpoint     = "/ip4/0.0.0.0/tcp/8888"
	DefaultStatsReportingInterval = 2 * time.Second

	DefaultEnableTracing       = false
	DefaultJaegerAgentEndpoint = "/ip4/0.0.0.0/udp/6831"
	DefaultTracingSamplingProb = 0.3
	DefaultTracingServiceName  = "cluster-daemon"
)

// MetricsConfig configures metrics collection.
type MetricsConfig struct {
	config.Saver

	EnableStats            bool
	PrometheusEndpoint     ma.Multiaddr
	StatsReportingInterval time.Duration
}

type jsonMetricsConfig struct {
	EnableStats            bool   `json:"enable_stats"`
	PrometheusEndpoint     string `json:"prometheus_endpoint"`
	StatsReportingInterval string `json:"reporting_interval"`
}

// ConfigKey provides a human-friendly identifier for this type of Config.
func (cfg *MetricsConfig) ConfigKey() string {
	return metricsConfigKey
}

// Default sets the fields of this Config to sensible values.
func (cfg *MetricsConfig) Default() error {
	cfg.EnableStats = DefaultEnableStats
	endpointAddr, _ := ma.NewMultiaddr(DefaultPrometheusEndpoint)
	cfg.PrometheusEndpoint = endpointAddr
	cfg.StatsReportingInterval = DefaultStatsReportingInterval

	return nil
}

// Validate checks that the fields of this Config have working values,
// at least in appearance.
func (cfg *MetricsConfig) Validate() error {
	if cfg.EnableStats {
		if cfg.PrometheusEndpoint == nil {
			return errors.New("metrics.prometheus_endpoint is undefined")
		}
		if cfg.StatsReportingInterval < 0 {
			return errors.New("metrics.reporting_interval is invalid")
		}
	}
	return nil
}

// LoadJSON sets the fields of this Config to the values defined by the JSON
// representation of it, as generated by ToJSON.
func (cfg *MetricsConfig) LoadJSON(raw []byte) error {
	jcfg := &jsonMetricsConfig{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		logger.Error("Error unmarshaling observations config")
		return err
	}

	cfg.Default()

	// override json config with env var
	err = envconfig.Process(envConfigKey, jcfg)
	if err != nil {
		return err
	}

	err = cfg.loadMetricsOptions(jcfg)
	if err != nil {
		return err
	}

	return cfg.Validate()
}

func (cfg *MetricsConfig) loadMetricsOptions(jcfg *jsonMetricsConfig) error {
	cfg.EnableStats = jcfg.EnableStats
	endpointAddr, err := ma.NewMultiaddr(jcfg.PrometheusEndpoint)
	if err != nil {
		return fmt.Errorf("loadMetricsOptions: PrometheusEndpoint multiaddr: %v", err)
	}
	cfg.PrometheusEndpoint = endpointAddr

	return config.ParseDurations(
		metricsConfigKey,
		&config.DurationOpt{
			Duration: jcfg.StatsReportingInterval,
			Dst:      &cfg.StatsReportingInterval,
			Name:     "metrics.reporting_interval",
		},
	)
}

// ToJSON generates a human-friendly JSON representation of this Config.
func (cfg *MetricsConfig) ToJSON() ([]byte, error) {
	jcfg := &jsonMetricsConfig{
		EnableStats:            cfg.EnableStats,
		PrometheusEndpoint:     cfg.PrometheusEndpoint.String(),
		StatsReportingInterval: cfg.StatsReportingInterval.String(),
	}

	return config.DefaultJSONMarshal(jcfg)
}

// TracingConfig configures tracing.
type TracingConfig struct {
	config.Saver

	EnableTracing       bool
	JaegerAgentEndpoint ma.Multiaddr
	TracingSamplingProb float64
	TracingServiceName  string
}

type jsonTracingConfig struct {
	EnableTracing       bool    `json:"enable_tracing"`
	JaegerAgentEndpoint string  `json:"jaeger_agent_endpoint"`
	TracingSamplingProb float64 `json:"sampling_prob"`
	TracingServiceName  string  `json:"service_name"`
}

// ConfigKey provides a human-friendly identifier for this type of Config.
func (cfg *TracingConfig) ConfigKey() string {
	return tracingConfigKey
}

// Default sets the fields of this Config to sensible values.
func (cfg *TracingConfig) Default() error {
	cfg.EnableTracing = DefaultEnableTracing
	agentAddr, _ := ma.NewMultiaddr(DefaultJaegerAgentEndpoint)
	cfg.JaegerAgentEndpoint = agentAddr
	cfg.TracingSamplingProb = DefaultTracingSamplingProb
	cfg.TracingServiceName = DefaultTracingServiceName
	return nil
}

// Validate checks that the fields of this Config have working values,
// at least in appearance.
func (cfg *TracingConfig) Validate() error {
	if cfg.EnableTracing {
		if cfg.JaegerAgentEndpoint == nil {
			return errors.New("tracing.jaeger_agent_endpoint is undefined")
		}
		if cfg.TracingSamplingProb < 0 {
			return errors.New("tracing.sampling_prob is invalid")
		}
	}
	return nil
}

// LoadJSON sets the fields of this Config to the values defined by the JSON
// representation of it, as generated by ToJSON.
func (cfg *TracingConfig) LoadJSON(raw []byte) error {
	jcfg := &jsonTracingConfig{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		logger.Error("Error unmarshaling observations config")
		return err
	}

	cfg.Default()

	// override json config with env var
	err = envconfig.Process(envConfigKey, jcfg)
	if err != nil {
		return err
	}

	err = cfg.loadTracingOptions(jcfg)
	if err != nil {
		return err
	}

	return cfg.Validate()
}

func (cfg *TracingConfig) loadTracingOptions(jcfg *jsonTracingConfig) error {
	cfg.EnableTracing = jcfg.EnableTracing
	agentAddr, err := ma.NewMultiaddr(jcfg.JaegerAgentEndpoint)
	if err != nil {
		return fmt.Errorf("loadTracingOptions: JaegerAgentEndpoint multiaddr: %v", err)
	}
	cfg.JaegerAgentEndpoint = agentAddr
	cfg.TracingSamplingProb = jcfg.TracingSamplingProb
	cfg.TracingServiceName = jcfg.TracingServiceName

	return nil
}

// ToJSON generates a human-friendly JSON representation of this Config.
func (cfg *TracingConfig) ToJSON() ([]byte, error) {
	jcfg := &jsonTracingConfig{
		EnableTracing:       cfg.EnableTracing,
		JaegerAgentEndpoint: cfg.JaegerAgentEndpoint.String(),
		TracingSamplingProb: cfg.TracingSamplingProb,
		TracingServiceName:  cfg.TracingServiceName,
	}

	return config.DefaultJSONMarshal(jcfg)
}
