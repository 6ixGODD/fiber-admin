package mods

type PrometheusConfig struct {
	Namespace  string `mapstructure:"prometheus_namespace" yaml:"prometheus_namespace" default:"data_collection_hub"`
	Subsystem  string `mapstructure:"prometheus_subsystem" yaml:"prometheus_subsystem" default:""`
	MetricPath string `mapstructure:"prometheus_metric_path" yaml:"prometheus_metric_path" default:"/metrics"`
}
