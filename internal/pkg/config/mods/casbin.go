package mods

type CasbinConfig struct {
	ModelPath        string `mapstructure:"casbin_model_path" yaml:"casbin_model_path" default:"./configs/casbin_model.test.conf"`
	PolicyAdapterUrl string `mapstructure:"casbin_policy_adapter_url" yaml:"casbin_policy_adapter_url" default:"mongodb://localhost:27017/fiber-admin"`
}
