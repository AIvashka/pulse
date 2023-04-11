package database

import (
	"context"
	"database/sql"
	"fmt"
	"pulse/structs"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

/*

CREATE TABLE spreads (
    id SERIAL PRIMARY KEY,
    symbol TEXT NOT NULL,
    quantity FLOAT NOT NULL,
    spread FLOAT NOT NULL,
    save_time TIMESTAMP NOT NULL
);


*/

func NewDatabase(config *structs.Config) (*Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	//fmt.Println(connectionString)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (d *Database) SaveSpreadRecords(records []*structs.SpreadRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO spreads(symbol, quantity, spread, save_time) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		_, err := stmt.ExecContext(ctx, record.Symbol, record.Quantity, record.Spread, record.SaveTime)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) SaveSpreadRecord(record *structs.SpreadRecord) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := d.db.ExecContext(ctx, "INSERT INTO spreads (symbol, quantity, spread, save_time) VALUES ($1, $2, $3, $4)", record.Symbol, record.Quantity, record.Spread, record.SaveTime)
	if err != nil {
		return err
	}

	return nil
}
