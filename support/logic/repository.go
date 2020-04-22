package support

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(context.Context, *Ticket) error
	Accept(context.Context, int, int) error
	Find(context.Context, int) (*Ticket, error)
	FindAll(context.Context, int) ([]*Ticket, error)
	ChangeStatus(context.Context, int, string) error
	AddMessage(context.Context, *TicketMessage) error
	TakeMessages(context.Context, int) ([]*TicketMessage, error)
}

type repository struct {
	db     *sqlx.DB
	logger log.Logger
}

func NewRepository(db *sqlx.DB, logger log.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, t *Ticket) error {
	t.BeforeCreate()

	return r.db.QueryRowContext(
		ctx,
		`insert into tickets (title, description, section, from_user, helper, created_at, status) 
		values ($1, $2, $3, $4, $5, $6, $7) returning id`,
		t.Title,
		t.Description,
		t.Section,
		t.From,
		t.Helper,
		t.CreatedAt,
		t.Status,
	).Scan(&t.ID)
}

func (r *repository) Accept(ctx context.Context, ticketId int, helper int) error {
	_, err := r.db.ExecContext(
		ctx,
		"update tickets set helper = $1 where id = $2",
		helper,
		ticketId,
	)

	return err
}

func (r *repository) Find(ctx context.Context, ticketId int) (*Ticket, error) {
	ticket := &Ticket{}

	if err := r.db.GetContext(ctx, ticket, "select * from tickets where id = $1", ticketId); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *repository) FindAll(ctx context.Context, userId int) ([]*Ticket, error) {
	var tickets []*Ticket
	if err := r.db.SelectContext(ctx, &tickets, "select * from tickets where from_user = $1", userId); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *repository) ChangeStatus(ctx context.Context, ticketId int, status string) error {
	_, err := r.db.ExecContext(
		ctx,
		"update tickets set status = $1 where id = $2",
		status,
		ticketId,
	)

	return err
}

func (r *repository) AddMessage(ctx context.Context, tm *TicketMessage) error {
	tm.BeforeCreate()

	return r.db.QueryRowContext(
		ctx,
		"insert into ticket_messages (who, ticket_id, message_text, sended_at) values ($1, $2, $3, $4) returning id",
		tm.Who,
		tm.TicketId,
		tm.Message,
		tm.SendedAt,
	).Scan(&tm.ID)
}

func (r *repository) TakeMessages(ctx context.Context, ticketId int) ([]*TicketMessage, error) {
	var messages []*TicketMessage
	if err := r.db.SelectContext(ctx, &messages, "select * from ticket_messages where ticket_id = $1", ticketId); err != nil {
		return nil, err
	}

	return messages, nil
}
