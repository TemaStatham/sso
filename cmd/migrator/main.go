package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migraionPath, migrationTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to starage")
	flag.StringVar(&migraionPath, "migration-path", "", "path to migration")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "name of migration table")
	flag.Parse()

	if storagePath == "" {
		panic("storage path is empty")
	}

	if migraionPath == "" {
		panic("migration path is empty")
	}

	m, err := migrate.New(
		"file://"+migraionPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationTable),
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}

		panic(err)
	}

	fmt.Println("migrations applyed successfuly")
}
