package notifications

import (
	"errors"
	"time"
)

var (
	ErrTicketStatusNotFound = errors.New("Status not found")
)

const (
	Info    = "info"
	Warnign = "warning"
	Error   = "error"
)

type Notification struct {
	ID        int    `json:"id" db:"id"`
	Message   string `json:"message" db:"message"`
	CreatedAt int64  `json:"date" db:"created_at"`
	Status    string `json:"status" db:"noti_status"`
	For       int    `json:"for" db:"for_user"`
	Checked   bool   `json:"checked" db:"checked"`
}

func (n *Notification) Validate() error {
	switch n.Status {
	case Info:
	case Warnign:
	case Error:
	default:
		return ErrTicketStatusNotFound
	}

	return nil
}

func (n *Notification) BeforeCreate() {
	n.CreatedAt = time.Now().UTC().Unix()
	n.Checked = false
}
