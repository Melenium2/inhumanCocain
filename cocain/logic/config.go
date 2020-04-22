package gates

type Config struct {
	DatabaseUrl string `toml:"database_url"`
	HttPort     string `toml:"http_port"`
	ConsulPort  string `toml:"consul_port"`
	MaxRetry    int    `toml:"max_retry"`
	MaxTimeout  int    `toml:"max_timeout"`
}

func NewConfig() *Config {
	return &Config{}
}
