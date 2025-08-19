package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type basePgRepository struct {
	pool *pgxpool.Pool
}

func NewBasePgRepository(pool *pgxpool.Pool) basePgRepository {
	return basePgRepository{pool: pool}
}
