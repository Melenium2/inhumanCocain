package user

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER  = "USER"
	ADMIN = "ADMIN"
)

// User model
type User struct {
	ID                int    `json:"id,omitempty"`
	Login             string `json:"login" db:"username"`
	Email             string `json:"email" db:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-" db:"encrypted_password"`
	CreatedAt         int64  `json:"registration_date" db:"created_at"`
	Token             string `json:"api_token" db:"token"`
	Contacts          string `json:"contacts" db:"contacts"`
	Role              string `json:"user_role" db:"role"`
	IsActive          bool   `json:"-" db:"is_active"`
}

func (user *User) Validate() error {
	return validation.ValidateStruct(
		user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.By(requiredIf(user.EncryptedPassword == "")), validation.Length(6, 100)),
	)
}

// Fill fields before user create
func (user *User) BeforeCreate() error {
	if len(user.Password) > 0 {
		enc, err := encryptString(user.Password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = enc
	}

	user.CreatedAt = time.Now().UTC().Unix()
	user.Token = generateToken(user)
	user.Role = USER
	user.IsActive = true

	return nil
}

// Set new password to the User
func (user *User) SetPassword(password string) error {
	if len(password) > 0 {
		enc, err := encryptString(password)
		if err != nil {
			return err
		}
		user.EncryptedPassword = enc
		return nil
	}

	return errors.New("Empty param: 'password'")
}

// Set new random api_token
func (user *User) GenerateNewToken() {
	user.Token = generateToken(user)
}

// Change account status
func (user *User) ChangeActiveStatus(newStatus bool) {
	user.IsActive = newStatus
}

// Compare password of user and request
func (user *User) ComparePassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(pwd)) == nil
}

// Generate new api_token
func generateToken(user *User) string {
	hash := md5.New()
	hash.Write([]byte(time.Now().String() + user.Email))
	return hex.EncodeToString(hash.Sum(nil))
}

// Condition for validates module
func requiredIf(cond bool) validation.RuleFunc {
	return func(value interface{}) error {
		if cond {
			return validation.Validate(value, validation.Required)
		}

		return nil
	}
}

// Encrypt password
func encryptString(s string) (string, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(enc), nil
}
