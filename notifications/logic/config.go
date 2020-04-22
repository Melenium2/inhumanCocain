package notifications

type Config struct {
	DatabaseUrl string `toml:"database_url"`
	GrpcPort    string `toml:"grpc_port"`
}

func NewConfig() *Config {
	return &Config{}
}
