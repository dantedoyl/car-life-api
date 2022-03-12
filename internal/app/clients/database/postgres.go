package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
)

type Postgres struct {
	postgresDatabase *sql.DB
}

func NewPostgres(dataSourceName string) (*Postgres, error) {
	sqlConn, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := sqlConn.Ping(); err != nil {
		return nil, err
	}

	return &Postgres{
		postgresDatabase: sqlConn,
	}, nil
}

func (p *Postgres) GetDatabase() *sql.DB {
	return p.postgresDatabase
}

func (p *Postgres) Close() {
	p.postgresDatabase.Close()
}
