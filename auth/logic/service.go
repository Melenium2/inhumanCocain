package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	Generate(context.Context, *GenerateRequest) (*GenerateResponse, error)
	Validate(context.Context, *ValidateRequest) (*ValidateResponse, error)
}

type service struct {
	config *Config
}

func NewService(c *Config) Service {
	return &service{
		config: c,
	}
}

type Claims struct {
	UserId int32
	Role   string
	jwt.StandardClaims
}

func (s *service) Generate(_ context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
		UserId: req.UserId,
		Role:   req.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 15).Unix(),
		},
	})

	signedJwt, err := jwt.SignedString([]byte(s.config.Jwtsercet))
	if err != nil {
		return &GenerateResponse{
			Err: err.Error(),
		}, nil
	}

	return &GenerateResponse{
		Token: signedJwt,
	}, nil
}

func (s *service) Validate(_ context.Context, req *ValidateRequest) (*ValidateResponse, error) {
	claims := &Claims{}
	if req.Token == "" {
		return &ValidateResponse{
			Err:"Empty token",
		}, nil
	}
	jwtToken, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Jwtsercet), nil
	})
	if jwtToken != nil && !jwtToken.Valid || err != nil {
		println(err.Error())
		return &ValidateResponse{
			Err: err.Error(),
		}, nil
	}

	return &ValidateResponse{
		Token:  req.Token,
		UserId: claims.UserId,
		Role:   claims.Role,
	}, nil
}
