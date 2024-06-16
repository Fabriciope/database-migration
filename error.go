package main

import (
	"errors"
	"fmt"
)

// NOTE: Migrator

var (
	CannotReadFileErr     error = newError("error reading file")
	EmptyMigrationsDirErr error = newError("migrations directory must have at leat two files")
	EmptyMigrationFileErr error = newError("empty migration")
	InvalidMigrationErr   error = fmt.Errorf("each migrations must have at least two file, %s and %s", upSuffix, downSuffix)
)

type MigratorError struct {
	File string
	Err  error
}

func (err *MigratorError) Error() string {
	return fmt.Sprintf("error: %s - at file: %s", err.Err.Error(), err.File)
}

func (err *MigratorError) Is(target error) bool {
	_, ok := target.(*MigratorError)
	return ok
}

func (err *MigratorError) Unwrap() error {
	switch t := err.Err.(type) {
	case interface{ Unwrap() []error }:
		return t.Unwrap()[0]

	}
	return err.Err
}

// NOTE: Migration

var (
	// TODO: upercase
	emptyUpMigrationErr         error = newError("empty up migration")
	emptyDownMigrationErr       error = newError("empty down migration")
	MigrationAlreadyExecutedErr error = newError("this migration is already up")
	MigrationNotExecutedYetErr  error = newError("this migration has not run yet")
)

type MigrationError struct {
	FileName   string
	SqlCommand string
	Err        error
}

func (err *MigrationError) Error() string {
	return fmt.Sprintf("error: %s - at file: %s", err.Err.Error(), err.FileName)
}

func (err *MigrationError) Is(target error) bool {
	_, ok := target.(*MigrationError)
	return ok
}

func (err *MigrationError) Unwrap() error {
	return err.Err
}

func newError(message string) error {
	return errors.New(message)
}
