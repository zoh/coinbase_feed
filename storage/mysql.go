package storage

import (
	"coinbase_feed/models"
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLClient ...
type MySQLClient struct {
	db *sql.DB
}

// NewMySQLClient create mysql client
func NewMySQLClient(uri string) (*MySQLClient, error) {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &MySQLClient{
		db: db,
	}, nil
}

// Close disconnect
func (m *MySQLClient) Close() error {
	return m.db.Close()
}

// CreateSchema create tables, indexes, ect
func (m *MySQLClient) CreateSchema() error {
	var err error
	row := m.db.QueryRow("SELECT DATABASE();")
	var res interface{}
	err = row.Scan(&res)
	if err != nil {
		return err
	}
	log.Println("selected db:", res)

	_, err = m.db.Exec(`Create Table IF NOT EXISTS coinbase_feeder (
		timestamp BIGINT, 
		symbol varchar(20),
		bid float(20),
		ask float(20)
	);`)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQLClient) Write(ctx context.Context, t models.Ticker) error {
	_, err := m.db.ExecContext(ctx, "Insert into coinbase_feeder (timestamp, symbol, bid, ask) values (?, ?, ?, ?)",
		t.Timestamp, t.Symbol, t.Bid, t.Ask)
	return err
}

func (m *MySQLClient) ReadLast(ctx context.Context, symbol string) (*models.Ticker, error) {
	var result models.Ticker
	row := m.db.QueryRowContext(ctx, "select timestamp, symbol, bid, ask from coinbase_feeder where symbol = ? order by timestamp desc limit 1", symbol)
	err := row.Scan(&result.Timestamp, &result.Symbol, &result.Bid, &result.Ask)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
