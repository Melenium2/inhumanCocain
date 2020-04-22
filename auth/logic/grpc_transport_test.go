package auth

import (
	"context"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/go-kit/kit/log"
	pbs "github.com/inhumanLightBackend/auth/pb"
	"github.com/stretchr/testify/assert"
)

var (
	newConfig = func() (*Config, error) {
		var configPath string = "_config.toml"
		config := NewConfig()
		_, err := toml.DecodeFile(configPath, config)
		if err != nil {
			return nil, err
		}

		return config, nil
	}
)

func TestGRPCGenerateShouldReturnNewTokenFromGivenRequest(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))
	
	tr := NewGRPCServer(ep, log.NewNopLogger())
	req := &pbs.GenerateRequest{
		UserId: 1,
		Role: "admin",
	}
	r, err := tr.Generate(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Token)
}

func TestGRPCGenerateShouldReturnNewTokenFromEmptyRequest(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))
	
	tr := NewGRPCServer(ep, log.NewNopLogger())
	req := &pbs.GenerateRequest{}
	r, err := tr.Generate(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, r.Token)
}

func TestGRPCValidateShouldReturnClaimsFromGivenToken(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))
	
	tr := NewGRPCServer(ep, log.NewNopLogger())
	req := &pbs.GenerateRequest{
		UserId: 1,
		Role: "admin",
	}
	r, err := tr.Generate(context.Background(), req)
	assert.NoError(t, err)
	token := r.Token
	claims, err := tr.Validate(context.Background(), &pbs.ValidateRequest{Token: token})
	assert.NoError(t, err)
	assert.Equal(t, claims.UserId, int32(1))
	assert.Equal(t, claims.Role, "admin")
}

func TestGRPCValidateShouldReturnErrorFromEmptyRequest(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))
	
	tr := NewGRPCServer(ep, log.NewNopLogger())
	claims, err := tr.Validate(context.Background(), &pbs.ValidateRequest{})
	assert.NoError(t, err)
	assert.Equal(t, "Empty token", claims.Err)
}
