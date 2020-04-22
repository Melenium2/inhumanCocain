package gates

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type OpaqueRepository interface {
	FindToken(context.Context, string) (*Opaque, error)
	SaveToken(context.Context, int, string) (*Opaque, error)
	RemoveToken(context.Context, string) error
}

type GateRepo struct {
	Db *sqlx.DB
	Logger log.Logger
}

func NewGateRepo(db *sqlx.DB, logger log.Logger) OpaqueRepository {
	return &GateRepo{
		Db: db,
		Logger: logger,
	}
}

// Give opaque, get jwt
func (gr *GateRepo) FindToken(ctx context.Context, opaque string) (*Opaque, error) {
	o := &Opaque{}
	if err := gr.Db.GetContext(
		ctx,
		o,
		"select * from opaque_store where opaque = $1",
		opaque,
	); err != nil {
		return nil, err
	}

	return o, nil
}

// Give jwt, get opaque
func (gr *GateRepo) SaveToken(ctx context.Context, userId int, jwt string) (*Opaque, error) {
	o := NewOpaque(userId, jwt)
	if err := gr.Db.QueryRowContext(ctx, 
		"insert into opaque_store (user_id, jwt, opaque, created_at) values ($1, $2, $3, $4) returning id",
		o.UserId,
		o.Jwt,
		o.Opaque,
		o.CreatedAt,
	).Scan(&o.Id); err != nil {
		return nil, err
	}
	
	return o, nil
}

func (gr *GateRepo) RemoveToken(ctx context.Context, opaque string) error {
	_, err := gr.Db.ExecContext(
		ctx,
		"delete from opaque_store where opaque = $1",
		opaque,
	)
	
	return err
}

