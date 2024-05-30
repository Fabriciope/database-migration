package main

import (
	"database/sql"
)

type Up interface {
	Up() error
}

type Down interface {
	Down() error
}

type Exec interface {
	Exec(string, ...any) (sql.Result, error)
}

// TODO: lower case
// TEST: testar novamente e comitar
type Migration struct {
	dbConnection   Exec
	Name           string
	upSqlCommand   string
	downSqlCommand string
	executed       bool
}

func (migration *Migration) Up() error {
	if migration.upSqlCommand == "" {
		return migration.newUpMigrationError(emptyUpMigrationErr)
	}

	if migration.executed {
		return migration.newUpMigrationError(cannotUpMigrationErr)
	}

	if _, err := migration.dbConnection.Exec(migration.upSqlCommand); err != nil {
		return migration.newUpMigrationError(err)
	}

	migration.executed = true

	return nil
}

func (migration *Migration) Down() error {
	if migration.downSqlCommand == "" {
		return migration.newDownMigrationError(emptyDownMigrationErr)
	}

	if !migration.executed {
		return migration.newDownMigrationError(cannotDownMigrationErr)
	}

	if _, err := migration.dbConnection.Exec(migration.downSqlCommand); err != nil {
		return migration.newDownMigrationError(err)
	}

	migration.executed = false

	return nil
}

func (migration *Migration) newUpMigrationError(targedErr error) error {
	return &MigrationError{
		FileName:   migration.Name + upSuffix,
		SqlCommand: migration.upSqlCommand,
		Err:        targedErr,
	}
}

func (migration *Migration) newDownMigrationError(targedErr error) error {
	return &MigrationError{
		FileName:   migration.Name + downSuffix,
		SqlCommand: migration.downSqlCommand,
		Err:        targedErr,
	}
}
