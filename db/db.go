package db

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type ConnectionConfig struct {
	Driver string
	DSN    string
}

func NewConnection(connConfig *ConnectionConfig) (*sqlx.DB, error) {
	if connConfig == nil {
		return nil, errors.New("connConfig cannot be null")
	}

	conn, err := sqlx.Connect(connConfig.Driver, connConfig.DSN)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
