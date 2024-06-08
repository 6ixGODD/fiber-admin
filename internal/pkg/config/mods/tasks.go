package mods

type TasksConfig struct {
	SyncLogsSpec  string `mapstructure:"sync_logs_spec" yaml:"sync_logs_spec" default:"@hourly"`
	UpdateKeySpec string `mapstructure:"update_key_spec" yaml:"update_key_spec" default:"@weekly"`
}
