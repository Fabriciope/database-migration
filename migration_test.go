// TODO: study if i must use _test prefix on package name
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
)

const migrationsPath = "./fake_migrations/"

func missingMigationFileFailure(t *testing.T, migrationPath string, err error) {
}

func getCommandFrom(t *testing.T, file string) string {
	content, err := os.ReadFile(file)
	if err != nil {
		t.Logf("error when load file: %s", file)
		t.Logf("received error: %s", err.Error())
		t.FailNow()
	}

	return string(content)
}

type execMock struct {
	mock.Mock
}

func (mock *execMock) Exec(query string, _ ...any) (sql.Result, error) {
	args := mock.Called(query)
	return nil, args.Error(0)
}

// TEST: testar o campo executed depois de rodar o down ou up
// TEST: testar com o coverage

func Test_Should_Run_Migration(t *testing.T) {
	var (
		migrationName = "create_users_table"
		migrationFile = fmt.Sprintf("%s_up.sql", migrationName)
		migrationPath = filepath.Join(migrationsPath, migrationFile)
	)

	command := getCommandFrom(t, migrationPath)

	execMock := &execMock{}
	execMock.On("Exec", command).Return(nil)
	defer execMock.AssertCalled(t, "Exec", command)

	migration := &Migration{
		dbConnection: execMock,
		FileName:     migrationFile,
		Path:         migrationPath,
		Name:         migrationName,
		executed:     false,
		upSqlCommand: command,
	}

	err := migration.Up()
	if err != nil {
		t.Errorf("exepected nil but got an error: %s", err.Error())
	}
}

func Test_Up_Migration_When_File_Is_Empty(t *testing.T) {
	var (
		migrationName = "empty_migration"
		migrationFile = fmt.Sprintf("%s_up.sql", migrationName)
		migrationPath = filepath.Join(migrationsPath, migrationFile)
	)

	command := getCommandFrom(t, migrationPath)

	execMock := &execMock{}
	execMock.On("Exec", command).Return(emptyUpMigrationErr)
	defer execMock.AssertNumberOfCalls(t, "Exec", 0)

	migration := &Migration{
		dbConnection: execMock,
		Name:         migrationName,
		FileName:     migrationFile,
		Path:         migrationPath,
		executed:     false,
		upSqlCommand: command,
	}

	err := migration.Up()
	if err == nil {
		t.Fatal("should return an error but got nil")
	}

	var migrationErr *MigrationError
	if ok := errors.As(err, &migrationErr); !ok {
		t.Errorf("invalid received error: err: %s", err.Error())
	}

	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr.Error() != emptyUpMigrationErr.Error() {
		t.Errorf("expected error *%s* but got: %s", emptyUpMigrationErr.Error(), unwrappedErr.Error())
	}

	if migrationErr.FileName != migrationFile {
		t.Errorf("exepected file %s but got %s", migrationName, migrationErr.FileName)
	}

	if migrationErr.SqlCommand != command {
		t.Errorf("expeced command %s but got %s", command, migrationErr.SqlCommand)
	}

    if migration.executed {
        t.Error("executed field must be false")
    }
}

func Test_Up_Migration_When_Up_Command_Already_Run(t *testing.T) {
	var (
		migrationName = "create_address_table"
		migrationFile = fmt.Sprintf("%s_up.sql", migrationName)
		migrationPath = filepath.Join(migrationsPath, migrationFile)
	)

	command := getCommandFrom(t, migrationPath)

	execMock := &execMock{}
	execMock.On("Exec", command).Return(nil)
	defer execMock.AssertNumberOfCalls(t, "Exec", 0)

	migration := &Migration{
		dbConnection: execMock,
		Name:         migrationName,
		FileName:     migrationFile,
		Path:         migrationPath,
		executed:     true,
		upSqlCommand: command,
	}

	err := migration.Up()
	if err == nil {
		t.Fatal("should return an error but got nil")
	}

	var migrationErr *MigrationError
	if ok := errors.As(err, &migrationErr); !ok {
		t.Errorf("invalid received error: %s", err.Error())
	}

	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr.Error() != cannotUpMigrationErr.Error() {
		t.Errorf("expected error *%s* but got: %s", cannotUpMigrationErr.Error(), unwrappedErr.Error())
	}

	if migrationErr.FileName != migrationFile {
		t.Errorf("exepected file %s but got %s", migrationName, migrationErr.FileName)
	}

	if migrationErr.SqlCommand != command {
		t.Errorf("expeced command %s but got %s", command, migrationErr.SqlCommand)
	}

    if !migration.executed {
        t.Error("executed field must be true")
    }
}

func Test_Up_Migration_When_Database_Return_An_Error(t *testing.T) {
	var (
		migrationName = "create_users_table"
		migrationFile = fmt.Sprintf("%s_up.sql", migrationName)
		migrationPath = filepath.Join(migrationsPath, migrationFile)
	)

	command := getCommandFrom(t, migrationPath)

	fakeErrorFromDatabase := errors.New("error by running sql command")

	execMock := &execMock{}
	execMock.On("Exec", command).Return(fakeErrorFromDatabase)
	defer execMock.AssertCalled(t, "Exec", command)

	migration := &Migration{
		dbConnection: execMock,
		FileName:     migrationFile,
		Path:         migrationPath,
		Name:         migrationName,
		executed:     false,
		upSqlCommand: command,
	}

	err := migration.Up()
	if err == nil {
		t.Fatal("should return an error but got nil")
	}

	var migrationErr *MigrationError
	if ok := errors.As(err, &migrationErr); !ok {
		t.Errorf("invalid received error: %s", err.Error())
	}

	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr.Error() != fakeErrorFromDatabase.Error() {
		t.Errorf("expected error *%s* but got: %s", fakeErrorFromDatabase.Error(), unwrappedErr.Error())
	}

	if migrationErr.FileName != migrationFile {
		t.Errorf("exepected file %s but got %s", migrationName, migrationErr.FileName)
	}

	if migrationErr.SqlCommand != command {
		t.Errorf("expeced command %s but got %s", command, migrationErr.SqlCommand)
	}
    
    if migration.executed {
        t.Error("executed field must be false")
    }
}

func Test_Down_Migration_When_Up_Command_Has_Not_Run_Yet(t *testing.T) {
	var (
		migrationName = "create_address_table"
		migrationFile = fmt.Sprintf("%s_down.sql", migrationName)
		migrationPath = filepath.Join(migrationsPath, migrationFile)
	)

	command := getCommandFrom(t, migrationPath)

	execMock := &execMock{}
	execMock.On("Exec", command).Return(nil)
	defer execMock.AssertNumberOfCalls(t, "Exec", 0)

	migration := &Migration{
		dbConnection:   execMock,
		Name:           migrationName,
		FileName:       migrationFile,
		Path:           migrationPath,
		executed:       false,
		downSqlCommand: command,
	}

	err := migration.Down()
	if err == nil {
		t.Fatal("should return an error but got nil")
	}

	var migrationErr *MigrationError
	if ok := errors.As(err, &migrationErr); !ok {
		t.Errorf("invalid received error: %s", err.Error())
	}

	unwrappedErr := errors.Unwrap(err)
	if unwrappedErr.Error() != cannotDownMigrationErr.Error() {
		t.Errorf("expected error *%s* but got: %s", cannotDownMigrationErr.Error(), unwrappedErr.Error())
	}

	if migrationErr.FileName != migrationFile {
		t.Errorf("exepected file %s but got %s", migrationName, migrationErr.FileName)
	}

	if migrationErr.SqlCommand != command {
		t.Errorf("expeced command %s but got %s", command, migrationErr.SqlCommand)
	}

    if migration.executed {
        t.Error("executed field must be false")
    }
}