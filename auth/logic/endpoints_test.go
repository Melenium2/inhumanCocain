package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	generateResponse = &GenerateResponse{
		Token: "Super token",
	}

	validateResponse = &ValidateResponse{
		UserId: 1,
		Role: "user",
		Token: "123123",
	}
)

type mockSuccessService struct {}

func (s *mockSuccessService) Generate(_ context.Context, _ *GenerateRequest) (*GenerateResponse, error) {
	return generateResponse, nil
}

func (s *mockSuccessService) Validate(_ context.Context, _ *ValidateRequest) (*ValidateResponse, error) {
	return validateResponse, nil
}

var successService = &mockSuccessService{}

type mockFailService struct {}

func (s *mockFailService) Generate(_ context.Context, _ *GenerateRequest) (*GenerateResponse, error) {
	return nil, errors.New("Fail")
}

func (s *mockFailService) Validate(_ context.Context, _ *ValidateRequest) (*ValidateResponse, error) {
	return nil, errors.New("Fail")
}

var failService = &mockFailService{}

func TestGenerateEndpointShouldReturnFuncThatReturnsNewToken(t *testing.T) {
	ep := NewEndpoints(successService)
	resp, err := ep.GenerateEndpoint(context.Background(), &GenerateRequest{})
	assert.NoError(t, err)
	genResp, ok := resp.(*GenerateResponse)
	assert.True(t, ok)
	assert.Equal(t, "Super token", genResp.Token)
}

func TestValidateEndpointShouldReturnFuncThatReturnsClaimsFromToken(t *testing.T) {
	ep := NewEndpoints(successService)
	r, err := ep.ValidateEndpoint(context.Background(), &ValidateRequest{})
	assert.NoError(t, err)
	resp, ok := r.(*ValidateResponse)
	assert.True(t, ok)
	assert.Equal(t, 1, int(resp.UserId))
}

func TestGenerateEndpointShouldReturnFuncThatReturnsErrorMessage(t *testing.T) {
	ep := NewEndpoints(failService)
	r, err := ep.GenerateEndpoint(context.Background(), &GenerateRequest{})
	assert.NoError(t, err)
	resp, ok := r.(*GenerateResponse)
	assert.True(t, ok)
	assert.Equal(t, "Fail", resp.Err)
}

func TestValidateEndpointShouldReturnFuncThatReturnsErrorMessage(t *testing.T) {
	ep := NewEndpoints(failService)
	r, err := ep.ValidateEndpoint(context.Background(), &ValidateRequest{})
	assert.NoError(t, err)
	resp, ok := r.(*ValidateResponse)
	assert.True(t, ok)
	assert.Equal(t, "Fail", resp.Err)
}

