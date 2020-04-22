package notifications

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateNotification(context.Context, *Notification) error
	FindNotificationsById(context.Context, int) ([]*Notification, error) 
	CheckNotification(context.Context, []int, int) error
}

type repositroy struct {
	db     *sqlx.DB
	logger log.Logger
}

func NewRepository(db *sqlx.DB, logger log.Logger) Repository {
	return &repositroy{
		db: db,
		logger: logger,
	}
}

func (r *repositroy) CreateNotification(ctx context.Context, n *Notification) error {
	if err := n.Validate(); err != nil {
		return err
	}
	n.BeforeCreate()

	return r.db.QueryRowContext(
		ctx,
		"insert into notifications (message, created_at, noti_status, for_user, checked) values($1, $2, $3, $4, $5) returning id",
		n.Message,
		n.CreatedAt,
		n.Status,
		n.For,
		n.Checked,
	).Scan(&n.ID)
}

func (r *repositroy) FindNotificationsById(ctx context.Context, userId int) ([]*Notification, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"select * from notifications where for_user = $1 and checked = false",
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := make([]*Notification, 0)
	for rows.Next() {
		notification := &Notification{}
		if err := rows.Scan(&notification.ID, &notification.Message, &notification.CreatedAt,
		&notification.Status, &notification.For, &notification.Checked); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *repositroy) CheckNotification(ctx context.Context, indexes []int, userId int) error {
	sIndexes := make([]string, 0)
	for _, n := range indexes {
		str := strconv.Itoa(n)
		sIndexes = append(sIndexes, str)
	}

	_, err := r.db.ExecContext(
		ctx,
		fmt.Sprintf("update notifications set checked = true where id in (%s) and for_user = $1", strings.Join(sIndexes, ",")),
		userId,
	)
	if err != nil {
		return err
	}

	return nil
}