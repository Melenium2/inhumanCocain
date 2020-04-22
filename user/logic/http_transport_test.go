package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

//Expample (No working)
func TestCreateUserSholdReturnSuccessMess(t *testing.T) {
	db, mock, err := createDb()
	assert.NoError(t, err)
	defer db.Close()

	user := &User{
		Contacts: "Content",
		Email: "john@email.com",
		Password: "1234567",
		Login: "john",
	}
	_ = user.BeforeCreate()

	mock.ExpectQuery(`insert into users 
		\(username, email, encrypted_password, created_at, token, contacts, role, is_active\)
			values \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\) returning id$`).
		WithArgs(user.Login, user.Email, user.EncryptedPassword, user.CreatedAt, user.Token, user.Contacts, user.Role, user.IsActive).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	repo := NewRepo(db)
	ep := NewEndpoints(NewService(repo))

	b, err := json.Marshal(user)
	assert.NoError(t, err)

	r := httptest.NewRequest("POST", "/api/v1/create", bytes.NewReader(b))
	w := httptest.NewRecorder()
	NewHTTPTransport(ep, log.NewNopLogger()).ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	
	body, err := ioutil.ReadAll(w.Body)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(body), "Success"))
}

func createDb() (*sqlx.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dbx := sqlx.NewDb(db, "sqlmock")
	return dbx, mock, nil 
}