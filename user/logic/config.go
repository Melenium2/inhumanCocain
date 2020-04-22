package user

import (
	"context"

	auth "github.com/inhumanLightBackend/auth/logic"
	"github.com/spf13/cast"
)

type Config struct {
	HTTPPort    string `toml:"http_port"`
	GRPCPort    string `toml:"gRPC_port"`
	DatabaseURL string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{
		HTTPPort: ":7070",
		GRPCPort: ":7071",
	}
}

func contextToClaims(ctx context.Context) map[string]interface{} {
	return cast.ToStringMap(ctx.Value(auth.CtxUserKey))
}
