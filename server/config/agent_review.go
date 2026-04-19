package config

type AgentReview struct {
	Enabled               bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Endpoint              string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	PublicBaseURL         string `mapstructure:"public-base-url" json:"publicBaseUrl" yaml:"public-base-url"`
	CallbackPath          string `mapstructure:"callback-path" json:"callbackPath" yaml:"callback-path"`
	RequestTimeoutSeconds int    `mapstructure:"request-timeout-seconds" json:"requestTimeoutSeconds" yaml:"request-timeout-seconds"`
	RequestToken          string `mapstructure:"request-token" json:"requestToken" yaml:"request-token"`
	CallbackToken         string `mapstructure:"callback-token" json:"callbackToken" yaml:"callback-token"`
}
