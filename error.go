package main

import (
	"errors"
	"fmt"
)

var (
	emptyUpMigrationErr      = errors.New("empty up migration")
	emptyDownMigrationErr    = errors.New("empty down migration")
	MigrationAlreadyExecutedErr = errors.New("this migration is already up")
	MigrationNotExecutedYetErr  = errors.New("this migration has not run yet")
)

type MigrationError struct {
	FileName   string
	SqlCommand string
	Err        error
}

func (err *MigrationError) Error() string {
	return fmt.Sprintf("error: %s - at file: %s", err.Err.Error(), err.FileName)
}

func (err *MigrationError) Is(targed error) bool {
	_, ok := targed.(*MigrationError)
	return ok
}

func (err *MigrationError) Unwrap() error {
	return err.Err
}
