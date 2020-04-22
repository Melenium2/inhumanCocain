package user


// import (
// 	"context"
// 	"errors"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	test_user = &User{
// 		Email: "user@gmail.com",
// 		Login: "Someone",
// 		CreatedAt: time.Now().Unix(),
		
// 	}
// )

// type mockRepo struct {}

// func (r *mockRepo) CreateUser(ctx context.Context, user *User) error {
// 	return nil
// }
// func (r *mockRepo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
// 	return test_user, nil
// }
// func (r *mockRepo) FindUserById(ctx context.Context, id int) (*User, error) {
// 	return test_user, nil
// }
// func (r *mockRepo) UpdateUser(ctx context.Context, user *User) error {
// 	return nil
// }

// var test_repo = &mockRepo{}

// type mockSusccessService struct {}

// func (s *mockSusccessService) CreateUser(ctx context.Context, user *User) (string, error) {
// 	return "Success", test_repo.CreateUser(ctx, user)
// }

// func (s *mockSusccessService)	FindUserByEmail(ctx context.Context, email string) (*User, error) {
// 	return test_repo.FindUserByEmail(ctx, email)
// }

// func (s *mockSusccessService)	FindUserById(ctx context.Context, id int) (*User, error) {
// 	return test_repo.FindUserById(ctx, id)
// }

// func (s *mockSusccessService)	UpdateUser(ctx context.Context, user *User) (string, error) {
// 	return "Updated", test_repo.UpdateUser(ctx, user)
// }

// var test_success_service = &mockSusccessService {}

// type mockFailService struct {}

// func (s *mockFailService) CreateUser(ctx context.Context, user *User) (string, error) {
// 	return "", errors.New("Fail")
// }

// func (s *mockFailService)	FindUserByEmail(ctx context.Context, email string) (*User, error) {
// 	return nil, errors.New("Fail")
// }

// func (s *mockFailService)	FindUserById(ctx context.Context, id int) (*User, error) {
// 	return nil, errors.New("Fail")
// }

// func (s *mockFailService)	UpdateUser(ctx context.Context, user *User) (string, error) {
// 	return "", errors.New("Fail")
// }

// var test_fail_service = &mockFailService {}

// func TestCreateUserEndpointShouldReturnFuncThatReturnsMessage(t *testing.T) {
// 	ep := NewEndpoints(test_success_service)
// 	resp, err := ep.CreateEndpoint(context.Background(), CreateUserRequest{})
// 	assert.Nil(t, err)
// 	mess, ok := resp.(CreateUserResponse)
// 	assert.True(t, ok)
// 	assert.Equal(t, mess.Ok, "Success")
// }

// func TestFindUserByEmailShoudReturnFuncThatReturnsUser(t *testing.T) {
// 	ep := NewEndpoints(test_success_service)
// 	resp, err := ep.FindByEmailEndpoint(context.Background(), FindUserByEmailRequest{})
// 	assert.NoError(t, err)
// 	user, ok := resp.(FindUserByEmailResponse)
// 	assert.True(t, ok)
// 	assert.Equal(t, user.User.Email, test_user.Email)
// }

// func TestFindUserByIdShoudReturnFuncThatReturnsUser(t *testing.T) {
// 	ep := NewEndpoints(test_success_service)
// 	resp, err := ep.FindByIdEndpoint(context.Background(), FindUserByIdRequest{})
// 	assert.NoError(t, err)
// 	user, ok := resp.(FindUserByIdResponse)
// 	assert.True(t, ok)
// 	assert.Equal(t, user.User.Email, test_user.Email)
// }

// func TestUpdateUserShoudReturnFuncThatReturnsMess(t *testing.T) {
// 	ep := NewEndpoints(test_success_service)
// 	resp, err := ep.UpdateEndpoint(context.Background(), UpdateUserRequest{})
// 	assert.Nil(t, err)
// 	mess, ok := resp.(UpdateUserResponse)
// 	assert.True(t, ok)
// 	assert.Equal(t, mess.Ok, "Updated")	
// }

// func TestCreateUserEndpointShouldReturnFuncThatReturnsError(t *testing.T) {
// 	ep := NewEndpoints(test_fail_service)
// 	resp, err := ep.CreateEndpoint(context.Background(), CreateUserRequest{})
// 	assert.Nil(t, err)
// 	mess, ok := resp.(CreateUserResponse)
// 	assert.True(t, ok)
// 	assert.Equal(t, mess.Err, "Fail")
// }

// func TestFindUserByEmailShoudReturnFuncThatReturnsError(t *testing.T) {
// 	ep := NewEndpoints(test_fail_service)
// 	resp, err := ep.FindByEmailEndpoint(context.Background(), FindUserByEmailRequest{})
// 	assert.NoError(t, err)
// 	user, ok := resp.(FindUserByEmailResponse)
// 	assert.True(t, ok)
// 	assert.Nil(t, user.User)
// }

// func TestFindUserByIdShoudReturnFuncThatReturnsError(t *testing.T) {
// 	ep := NewEndpoints(test_fail_service)
// 	resp, err := ep.FindByIdEndpoint(context.Background(), FindUserByIdRequest{})
// 	assert.NoError(t, err)
// 	user, ok := resp.(FindUserByIdResponse)
// 	assert.True(t, ok)
// 	assert.Nil(t, user.User)
// }

// func TestUpdateUserShoudReturnFuncThatReturnsError(t *testing.T) {
// 	ep := NewEndpoints(test_fail_service)
// 	resp, err := ep.UpdateEndpoint(context.Background(), UpdateUserRequest{})
// 	assert.Nil(t, err)
// 	mess, ok := resp.(UpdateUserResponse)
// 	assert.True(t, ok)
// 	assert.Equal(t, mess.Err, "Fail")	
// }