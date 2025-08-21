package mysql

import "database/sql"

type baseMysqlRepository struct {
	db *sql.DB
}

func NewBaseMysqlRepository(db *sql.DB) baseMysqlRepository {
	return baseMysqlRepository{db: db}
}
