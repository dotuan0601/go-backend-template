package mysql

import (
	"database/sql"
	"fmt"

	"github.com/dotuan0601/go-backend-template/internal/config"
)

type MySQL struct {
	DB *sql.DB
}

func New(
	cfg config.MySQLConfig,
) (*MySQL, error) {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
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

	return &MySQL{
		DB: db,
	}, nil
}

func (m *MySQL) Close() error {

	if m.DB == nil {
		return nil
	}

	return m.DB.Close()
}
