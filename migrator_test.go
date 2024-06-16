package main

import (
	"errors"
	"testing"
)

const (
	// TODO: make more migrations file for testing
	fakeMigrationsDir1 string = "./fake_migrations1_for_tests/"
	fakeMigrationsDir2 string = "./fake_migrations2_for_tests/"
	fakeMigrationsDir3 string = "./fake_migrations3_for_tests/"
	fakeMigrationsDir4 string = "./fake_migrations4_for_tests/"
)

// TEST: add test suite and make more tests with diferents cases

func Test_NewMigrator(t *testing.T) {
	execMock := new(execMock)
	migrator, err := NewMigrator(fakeMigrationsDir2, execMock)
	if err != nil {
		t.Errorf("expected nil but got error: %s", err.Error())
	}

	if migrator == nil {
		t.Errorf("expected migrator but got nil")
	}
}

func Test_NewMigrator_With_Empty_Dir(t *testing.T) {
	execMock := new(execMock)
	migrator, err := NewMigrator(fakeMigrationsDir3, execMock)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if migrator != nil {
		t.Errorf("expected nil but got migrator")
	}

	if ok := errors.Is(err, EmptyMigrationsDirErr); !ok {
		t.Errorf("expected error: %s, but got: %s", EmptyMigrationsDirErr.Error(), err.Error())
	}
}

func Test_NewMigrator_With_Empty_File(t *testing.T) {
	execMock := new(execMock)
	migrator, err := NewMigrator(fakeMigrationsDir1, execMock)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if migrator != nil {
		t.Errorf("expected nil but got migrator")
	}

	if !errors.Is(err, EmptyMigrationFileErr) {
		t.Errorf("expected error: %s, but got: %s", EmptyMigrationFileErr.Error(), err.Error())
	}
}

func Test_NewMigrator_Without_Both_Migrations_File(t *testing.T) {
	execMock := new(execMock)
	migrator, err := NewMigrator(fakeMigrationsDir4, execMock)
	if err == nil {
		t.Errorf("expected error but got nil")
	}

	if migrator != nil {
		t.Errorf("expected nil but got migrator")
	}

	if !errors.Is(err, InvalidMigrationErr) {
		t.Errorf("expected error: %s, but got: %s", InvalidMigrationErr.Error(), err.Error())
	}
}
