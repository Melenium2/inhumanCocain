package notifications

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateNotificationShouldReturnNoError(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	n := &Notification{
		Message: "Supper message",
		CreatedAt: time.Now().Unix(),
		Checked: false,
		For: 1,
		Status: Info,
	}

	mock.ExpectQuery(`^insert into notifications \(message, created_at, noti_status, for_user, checked\) values\(\$1, \$2, \$3, \$4, \$5\) returning id$`).
			WithArgs(n.Message, n.CreatedAt, n.Status, n.For, n.Checked).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	assert.NoError(t, NewRepository(db, log.NewNopLogger()).CreateNotification(context.Background(), n))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateNotificationShouldReturnError(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	n := &Notification{
		Message: "Supper message",
		CreatedAt: time.Now().Unix(),
		Checked: false,
		For: 1,
		Status: "random",
	}

	mock.ExpectQuery(`^insert into notifications \(message, created_at, noti_status, for_user, checked\) values\(\$1, \$2, \$3, \$4, \$5\) returning id$`).
			WithArgs(n.Message, n.CreatedAt, n.Status, n.For, n.Checked).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	assert.Error(t, NewRepository(db, log.NewNopLogger()).CreateNotification(context.Background(), n))
	assert.Error(t, mock.ExpectationsWereMet())
}

func TestFindNotificationByIdShouldReturnListOfNotifiaction(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "message", "created_at", "noti_status", "for_user", "checked"}).
			AddRow(1, "New Message", time.Now().Unix(), Info, 1, false).
			AddRow(2, "New Message 2", time.Now().Unix(), Info, 1, false)

	mock.ExpectQuery(`select \* from notifications where for_user = \$1 and checked = false`).
			WithArgs(1).
			WillReturnRows(rows)
	
	res, err := NewRepository(db, log.NewNopLogger()).FindNotificationsById(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(res))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindNotificationByIdShouldReturnErrFromDb(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	rows := sqlmock.NewRows([]string{"id", "message", "created_at", "noti_status", "for_user", "checked"}).
			AddRow(1, "New Message", time.Now().Unix(), Info, 1, false).
			AddRow(2, "New Message 2", time.Now().Unix(), Info, 1, false)

	mock.ExpectQuery(`select \* from notifications where for_user = \$1 and checked = false`).
			WithArgs(1).
			WillReturnRows(rows)
	
	res, err := NewRepository(db, log.NewNopLogger()).FindNotificationsById(context.Background(), 2)
	assert.Error(t, err)
	assert.Nil(t, res)
	assert.Error(t, mock.ExpectationsWereMet())
}

func TestCheckNotificationShouldReturnNoErr(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	indexes := []int{1, 2}
	sIndexes := []string{"1", "2"}
	userId := 1
	mock.ExpectExec(fmt.Sprintf("^update notifications set checked = true where id in \\(%s\\) and for_user = \\$1$", strings.Join(sIndexes, ","))).
			WithArgs(userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

	assert.NoError(t, NewRepository(db, log.NewNopLogger()).CheckNotification(context.Background(), indexes, userId))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCheckNotificationShouldReturnErrWhenUserNotFound(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	indexes := []int{1, 2}
	sIndexes := []string{"1", "2"}
	userId := -381
	mock.ExpectExec(fmt.Sprintf("^update notifications set checked = true where id in \\(%s\\) and for_user = \\$1$", strings.Join(sIndexes, ","))).
			WithArgs(userId)

	assert.Error(t, NewRepository(db, log.NewNopLogger()).CheckNotification(context.Background(), indexes, userId))
}


func TestCheckNotificationShouldReturnErrWhenNotificationsNotFound(t *testing.T) {
	db, mock, err := setUpTestDatabase()
	assert.NoError(t, err)

	indexes := []int{3, 4}
	sIndexes := []string{"3", "4"}
	userId := 1
	mock.ExpectExec(fmt.Sprintf("^update notifications set checked = true where id in \\(%s\\) and for_user = \\$1$", strings.Join(sIndexes, ","))).
			WithArgs(userId)

	assert.Error(t, NewRepository(db, log.NewNopLogger()).CheckNotification(context.Background(), indexes, userId))
}