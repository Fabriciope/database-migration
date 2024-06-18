package main

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type execMock struct {
	mock.Mock
}

func (mock *execMock) Exec(query string, _ ...any) (sql.Result, error) {
	args := mock.Called(query)
	return nil, args.Error(0)
}

