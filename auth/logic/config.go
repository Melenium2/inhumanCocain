package auth

const (
	CtxUserKey ctxKey = iota
	CtxErrorKey ctxKey = iota
)

type ctxKey int32

type Config struct {
	Jwtsercet     string `toml:"jwt_secret"`
	Jwtexpiration int64    `toml:"jwt_expiration"`
	GrpcPort      string `toml:"grpc_port"`
}

func NewConfig() *Config {
	return &Config{}
}