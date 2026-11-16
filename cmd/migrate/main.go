package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/dotuan0601/go-backend-template/internal/config"
)

func main() {

	env := flag.String(
		"env",
		"dev",
		"application environment",
	)

	direction := flag.String(
		"direction",
		"up",
		"migration direction",
	)

	flag.Parse()

	// =========================
	// Config
	// =========================

	configPath := fmt.Sprintf(
		"configs/config-%s.yml",
		*env,
	)

	cfg, err := config.Load(
		configPath,
	)

	if err != nil {

		log.Fatal(err)
	}

	// =========================
	// Migration DSN
	// =========================

	dsn := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%d)/%s",
		cfg.MySQL.Username,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.Database,
	)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)

	if err != nil {

		log.Fatal(err)
	}

	// =========================
	// Execute
	// =========================

	switch *direction {

	case "up":

		err = m.Up()

	case "down":

		err = m.Steps(
			-1,
		)

	default:

		log.Fatalf(
			"unsupported direction: %s",
			*direction,
		)
	}

	if err != nil &&
		err != migrate.ErrNoChange {

		log.Fatal(err)
	}

	fmt.Println(
		"migration completed",
	)
}
