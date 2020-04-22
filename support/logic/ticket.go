package support

import "time"

const (
	Opened    = "opened"
	InProcess = "in process"
	Closed    = "closed"
)

// Ticket model
type Ticket struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Section     string `json:"section" db:"section"`
	From        int    `json:"from" db:"from_user"`
	Helper      int    `json:"helper" db:"helper"`
	CreatedAt   int64  `json:"created_at" db:"created_at"`
	Status      string `json:"status" db:"status"`
}

// Fill fields before ticket create
func (t *Ticket) BeforeCreate() {
	t.Helper = -1
	t.CreatedAt = time.Now().UTC().Unix()
	t.Status = Opened
}

// Ticket message model
type TicketMessage struct {
	ID       int   `json:"id" db:"id"`
	Who      int   `json:"who" db:"who"`
	TicketId int   `json:"ticket_id" db:"ticket_id"`
	Message  string `json:"message" db:"message_text"`
	SendedAt int64  `json:"sended_at" db:"sended_at"`
}

// Fill fields before message create
func (tm *TicketMessage) BeforeCreate() {
	tm.SendedAt = time.Now().UTC().Unix()
}
