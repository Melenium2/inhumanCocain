package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	errRecordNotFound = errors.New("Record not found")
)

type Repositiry interface {
	CreateUser(context.Context, *User) error
	FindUserByEmail(context.Context, string) (*User, error)
	FindUserById(context.Context, int) (*User, error)
	UpdateUser(context.Context, *User) error
}

type repo struct {
	db     *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repositiry {
	return &repo{
		db:     db,
	}
}

func (r *repo) CreateUser(ctx context.Context, user *User) error {
	return r.db.QueryRowxContext(
		ctx, 
		"insert into users (username, email, encrypted_password, created_at, token, contacts, role, is_active) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id",
		user.Login,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
		user.Token,
		user.Contacts,
		user.Role,
		user.IsActive,
	).Scan(&user.ID)
}

func (r *repo) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	if err := r.db.GetContext(
		ctx, 
		user,
		"select * from users where email = $1",
		email,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) FindUserById(ctx context.Context, id int) (*User, error) {
	user := &User{}
	if err := r.db.GetContext(
		ctx, 
		user,
		"select * from users where id = $1",
		id,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (r *repo) UpdateUser(ctx context.Context, user *User) error {
	_, err := r.db.ExecContext(
		ctx,
		`update users set 
		username = $2, email = $3, encrypted_password = $4, created_at = $5, 
		token = $6, contacts = $7, role = $8, is_active = $9 
		where id = $1`,
		user.ID,
		user.Login,
		user.Email,
		user.EncryptedPassword,
		user.CreatedAt,
		user.Token,
		user.Contacts,
		user.Role,
		user.IsActive,
	)

	if err != nil {
		return err
	}

	return nil
}
