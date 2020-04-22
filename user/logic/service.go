package user

import (
	"context"
	"reflect"

	auth "github.com/inhumanLightBackend/auth/logic"
)

type Service interface {
	Authenticate(context.Context, string, string) (*User, error)
	CreateUser(context.Context, *User) (string, error)
	FindUserByEmail(context.Context, string) (*User, error)
	FindUserById(context.Context, int) (*User, error)
	UpdateUser(context.Context, *User) (string, error)
}

type service struct {
	repository Repositiry
}

func NewService(repo Repositiry) Service {
	return &service{
		repository: repo,
	}
}

// Check user in db. If user exist and password right, then return user model.
func (s *service) Authenticate(ctx context.Context, login string, pass string) (*User, error) {
	user, err := s.repository.FindUserByEmail(ctx, login)
	if err != nil {
		return nil, err
	}
	if !user.ComparePassword(pass) {
		return nil, err
	}

	return user, nil
}

// Validate then create user frm given data
func (s *service) CreateUser(ctx context.Context, user *User) (string, error) {
	if err := user.Validate(); err != nil {
		return "", err
	}

	if err := user.BeforeCreate(); err != nil {
		return "", err
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return "", err
	}

	return "Success", nil
}

// Find user by email and return model if query executed without errors
func (s *service) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repository.FindUserByEmail(ctx, email)
}

func (s *service) FindUserById(ctx context.Context, id int) (*User, error) {
	return s.repository.FindUserById(ctx, id)
}

// Update user by given fields
func (s *service) UpdateUser(ctx context.Context, user *User) (string, error) {
	var response *auth.ValidateResponse
	{
		userRawId := ctx.Value(auth.CtxUserKey).(string)

		serv, _, err := AuthGRPCService(AuthserverPort, nil)
		if err != nil {
			return "", err
		}
		response, err = serv.Validate(ctx, &auth.ValidateRequest{
			Token: userRawId,
		})
		if err != nil {
			return "", err
		}
	}

	if response.Role == USER && user.Role != "" {
		return "", errPermissionsDeni
	}
	if user.Token != "" || user.EncryptedPassword != "" || user.Password != "" {
		return "", errPermissionsDeni
	}

	var (
		newUser *User
		err error
	)
	{
		isZeroValue := func(x interface{}) bool {
			return x == reflect.Zero(reflect.TypeOf(x)).Interface()
		}
		newUser, err = s.FindUserById(ctx, int(response.UserId))
		if err != nil {
			return "", err
		}
		m := reflect.ValueOf(user)
		if m.Kind() == reflect.Ptr {
			m = m.Elem()
			for i := 0; i < m.NumField(); i++ {
				field := m.Field(i)
				if !isZeroValue(field.Interface()) {
					reflect.ValueOf(newUser).Elem().Field(i).Set(field)
				}
			}
		}
	}

	if err := newUser.Validate(); err != nil {
		return "", err
	}
	if err := s.repository.UpdateUser(ctx, newUser); err != nil {
		return "", err
	}

	return "Updated", nil
}
