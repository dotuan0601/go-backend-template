package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dotuan0601/go-backend-template/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

type DataSource struct {
	db *sql.DB
}

func New(cfg config.MySQLConfig) (*DataSource, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open(
		"mysql",
		dsn,
	)

	if err != nil {
		return nil, err
	}

	// Connection pool config
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Verify database connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DataSource{
		db: db,
	}, nil
}

func (d *DataSource) DB() *sql.DB {
	return d.db
}

func (d *DataSource) Close() error {
	if d.db == nil {
		return nil
	}

	return d.db.Close()
}
