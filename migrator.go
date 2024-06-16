package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	upSuffix   string = "_up.sql"
	downSuffix string = "_down.sql"
)

type Migrator struct {
	migrations    []*Migration
	migrationsDir string
}

func NewMigrator(migrationDir string, conn Exec) (*Migrator, error) {
	dir, err := os.Open(migrationDir)
	if err != nil {
		return nil, fmt.Errorf("error when load directory: %s - err: %s", migrationDir, err.Error())
	}

	migrationsName, err := dir.Readdirnames(-1)
	if err != nil || len(migrationsName) == 0 {
		return nil, EmptyMigrationsDirErr
	}
	dir.Close()

	upMigrations, downMigrations, err := splitUpAndDownMigrations(migrationsName)
	if err != nil {
		return nil, err
	}

	if len(upMigrations) != len(downMigrations) {
		return nil, InvalidMigrationErr
	}

	migrator := new(Migrator)

	for _, upMigrationName := range upMigrations {
		upMigrationPath := filepath.Join(migrationDir, upMigrationName)
		upSqlCommand, err := getSqlCommandFrom(upMigrationPath)
		if err != nil {
			return nil, err
		}

		downMigrationName := strings.Replace(upMigrationName, upSuffix, downSuffix, 1)
		downMigrationIndex := slices.Index(downMigrations, downMigrationName)
		if downMigrationIndex < 0 {
			return nil, &MigratorError{
				File: upMigrationName,
				Err:  errors.New("migrations does not have down migration"),
			}
		}

		downMigrationPath := filepath.Join(migrationDir, downMigrationName)
		downSqlCommand, err := getSqlCommandFrom(downMigrationPath)
		if err != nil {
			return nil, err
		}

		migration := &Migration{
			dbConnection:   conn,
			Name:           strings.TrimSuffix(upMigrationName, upSuffix),
			upSqlCommand:   upSqlCommand,
			downSqlCommand: downSqlCommand,
			executed:       false,
		}

		migrator.migrations = append(migrator.migrations, migration)
		downMigrations = append(downMigrations[0:downMigrationIndex], downMigrations[downMigrationIndex+1:]...)
	}

	return migrator, nil
}

func splitUpAndDownMigrations(allMigrationsName []string) ([]string, []string, error) {
	capacity := len(allMigrationsName) / 2
	upMigrations := make([]string, 0, capacity)
	downMigrations := make([]string, 0, capacity)

	for _, fileName := range allMigrationsName {
		if strings.HasSuffix(fileName, upSuffix) {
			upMigrations = append(upMigrations, fileName)
		} else if strings.HasSuffix(fileName, downSuffix) {
			downMigrations = append(downMigrations, fileName)
		} else {
			return nil, nil, &MigratorError{
				File: fileName,
				Err:  fmt.Errorf("file must have %s or %s suffix", upSuffix, downSuffix),
			}
		}
	}
	return upMigrations, downMigrations, nil
}

func getSqlCommandFrom(migrationFile string) (string, error) {
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return "", &MigratorError{
			File: migrationFile,
			Err:  fmt.Errorf("%w - error caused by: %w", CannotReadFileErr, err),
		}
	}

	command := string(content)
	if command == "" {
		return "", &MigratorError{
			File: migrationFile,
			Err:  EmptyMigrationFileErr,
		}
	}
	return command, nil
}

// TEST: estudar como implementar a ordem de prioridade das migrations (up and down)
func (migrator *Migrator) UpAll() error {
	return nil
}

func (migrator *Migrator) DownAll() error {
	return nil
}

func (migrtor *Migrator) WalkUpMigrations(func(string, Up), error) error {
	return nil
}

func (migrtor *Migrator) WalkDownMigrations(func(string, Down) error) error {
	return nil
}
