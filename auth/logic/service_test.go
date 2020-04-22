package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShouldReturnResponseWithToken(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	resp, err := NewService(config).Generate(context.Background(), &GenerateRequest{
		UserId: 1,
		Role: "user",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestGenerateWithEmptyRequestShouldReturnResponseWithToken(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	resp, err := NewService(config).Generate(context.Background(), &GenerateRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestValidateShouldReturnUserClaims(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	service := NewService(config)
	resp, err := service.Generate(context.Background(), &GenerateRequest{
		UserId: 1,
		Role: "user",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	
	claims, err := service.Validate(context.Background(), &ValidateRequest{
		Token: resp.Token,
	})
	assert.NoError(t, err)
	assert.Equal(t, int32(1), claims.UserId)
	assert.Equal(t, "user", claims.Role)
}

func TestValidateShouldReturnResponseWithErrorEmptyToken(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	service := NewService(config)
	claims, err := service.Validate(context.Background(), &ValidateRequest{})
	assert.NoError(t, err)
	assert.Equal(t, "Empty token", claims.Err)
}

func TestValidateShouldReturnResponseWithError(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	service := NewService(config)
	claims, err := service.Validate(context.Background(), &ValidateRequest{
		Token: "e123131232",
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, claims.Err)
}

