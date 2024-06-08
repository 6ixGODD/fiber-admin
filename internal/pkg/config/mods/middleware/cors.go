package middleware

type CorsConfig struct {
	AllowOrigins     string `mapstructure:"allow_origins" yaml:"allow_origins" default:"*"`
	AllowMethods     string `mapstructure:"allow_methods" yaml:"allow_methods" default:"GET,POST,HEAD,PUT,DELETE,PATCH"`
	AllowHeaders     string `mapstructure:"allow_headers" yaml:"allow_headers" default:""`
	AllowCredentials bool   `mapstructure:"allow_credentials" yaml:"allow_credentials" default:"false"`
	ExposeHeaders    string `mapstructure:"expose_headers" yaml:"expose_headers" default:""`
	MaxAge           int    `mapstructure:"max_age" yaml:"max_age" default:"0"`
}
