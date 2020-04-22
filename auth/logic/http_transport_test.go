package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

func TestHTTPGenerateShouldReturnNewTokenFromGivenRequest(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))
	
	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(&GenerateRequest {
		UserId: 1,
		Role: "user",
	})
	req := httptest.NewRequest("POST", "/generate", body)
	resp := httptest.NewRecorder()
	NewHTTPTransport(ep, log.NewNopLogger()).ServeHTTP(resp, req)
	r, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, resp.Code, 200)
	assert.True(t, strings.Contains(string(r), "token"))
}

func TestHTTPGenerateShouldReturnResponseWithToken(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))

	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(&GenerateRequest{})
	req := httptest.NewRequest("POST", "/generate", body)
	resp := httptest.NewRecorder()
	NewHTTPTransport(ep, log.NewNopLogger()).ServeHTTP(resp, req)
	r, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, resp.Code, 200)
	assert.True(t, strings.Contains(string(r), "token"))
}

func TestHTTPValidateShouldReturnResponseWithClaims(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))

	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(&GenerateRequest {
		UserId: 1,
		Role: "user",
	})
	req := httptest.NewRequest("POST", "/generate", body)
	resp := httptest.NewRecorder()
	http := NewHTTPTransport(ep, log.NewNopLogger())
	http.ServeHTTP(resp, req)
	req = httptest.NewRequest("POST", "/validate", resp.Body)
	http.ServeHTTP(resp, req)
	r, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, resp.Code, 200)
	assert.True(t, strings.Contains(string(r), "role"))
}

func TestHTTPValidateShouldReturnResponseWith500Error(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))

	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(&ValidateRequest {})
	req := httptest.NewRequest("POST", "/validate", body)
	resp := httptest.NewRecorder()
	http := NewHTTPTransport(ep, log.NewNopLogger())
	http.ServeHTTP(resp, req)
	assert.Equal(t, 500, resp.Code)
}

func TestHTTPValidateShouldReturnResponseWithErrorField(t *testing.T) {
	config, err := newConfig()
	assert.NoError(t, err)
	ep := NewEndpoints(NewService(config))

	body := &bytes.Buffer{}
	json.NewEncoder(body).Encode(&ValidateRequest {
		Token: "e1313131kfskfdkm",
	})
	req := httptest.NewRequest("POST", "/validate", body)
	resp := httptest.NewRecorder()
	http := NewHTTPTransport(ep, log.NewNopLogger())
	http.ServeHTTP(resp, req)
	r, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.Code)
	assert.True(t, strings.Contains(string(r), "error"))
}
