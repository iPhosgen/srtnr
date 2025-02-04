package storage

import (
	"context"
	"fmt"

	"github.com/iphosgen/srtnr/config"
	"github.com/jackc/pgx/v4"
)

type Storage interface {
	Save(url string, userId string, shorted string) error
	Lookup(shorted string) (url string, err error)
}

type PostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresStorage(config *config.DatabaseConfig) (*PostgresStorage, error) {
	dsn := config.BuildDSN()
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return &PostgresStorage{conn: conn}, nil
}

func (s *PostgresStorage) Save(url string, userId string, shorted string) (err error) {
	var id string
	query := `INSERT INTO urls (url, user_id, shorted) VALUES ($1, $2, $3)`
	err = s.conn.QueryRow(context.Background(), query, url, userId).Scan(&id)
	if err != nil {
		err = fmt.Errorf("failed to save url: %v", err)
	}
	return
}

func (s *PostgresStorage) Lookup(shorted string) (url string, err error) {
	query := `SELECT url FROM urls WHERE shorted = $1`
	err = s.conn.QueryRow(context.Background(), query, shorted).Scan(&url)
	if err != nil {
		err = fmt.Errorf("failed to lookup url: %v", err)
	}
	return
}
