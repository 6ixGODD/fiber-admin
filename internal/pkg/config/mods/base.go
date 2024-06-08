package mods

type BaseConfig struct {
	AppName          string `mapstructure:"app_name" yaml:"app_name" default:"fiber-admin"`
	AppPort          string `mapstructure:"app_port" yaml:"app_port" default:"3000"`
	AppHost          string `mapstructure:"app_host" yaml:"app_host" default:"localhost"`
	AppVersion       string `mapstructure:"app_version" yaml:"app_version" default:"1.0.0"`
	EnableTls        bool   `mapstructure:"enable_tls" yaml:"enable_tls" default:"false"`
	TlsCertFile      string `mapstructure:"cert_file" yaml:"cert_file" default:""`
	TlsKeyFile       string `mapstructure:"key_file" yaml:"key_file" default:""`
	EnableCors       bool   `mapstructure:"enable_cors" yaml:"enable_cors" default:"true"`
	EnablePrometheus bool   `mapstructure:"enable_prometheus" yaml:"enable_prometheus" default:"true"`
}
