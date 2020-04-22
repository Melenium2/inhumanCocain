package notifications

import (
	"context"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	setUpTestDatabase = func() (*sqlx.DB, sqlmock.Sqlmock, error) {
		db, mock, err := sqlmock.New()
		if err != nil {
			return nil, nil, err
		}

		xdb := sqlx.NewDb(db, "sqlmock")
		return xdb, mock, nil
	}
)

func TestCreateShouldReturnCreatedMessage(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)
	defer db.Close()

	p := &Notification{
		Message: "Supper message",
		CreatedAt: time.Now().Unix(),
		Checked: false,
		For: 1,
		Status: Info,
	}

	mock.ExpectQuery(`^insert into notifications \(message, created_at, noti_status, for_user, checked\) values\(\$1, \$2, \$3, \$4, \$5\) returning id$`).
			WithArgs(p.Message, p.CreatedAt, p.Status, p.For, p.Checked).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))


	repo := NewRepository(db, log.NewNopLogger())
	res, err := NewService(repo, log.NewNopLogger()).Create(context.Background(), p)
	assert.NoError(t, err)
	assert.Equal(t, "Created", res)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateShouldReturnErrorWithStatusNotFound(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)
	defer db.Close()

	p := &Notification{
		Message: "Supper message",
		CreatedAt: time.Now().Unix(),
		Checked: false,
		For: 1,
		Status: "random",
	}

	mock.ExpectQuery(`^insert into notifications \(message, created_at, noti_status, for_user, checked\) values\(\$1, \$2, \$3, \$4, \$5\) returning id$`).
			WithArgs(p.Message, p.CreatedAt, p.Status, p.For, p.Checked).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))


	repo := NewRepository(db, log.NewNopLogger())
	_, err = NewService(repo, log.NewNopLogger()).Create(context.Background(), p)
	assert.Error(t, err)
	assert.Error(t, mock.ExpectationsWereMet())
}

