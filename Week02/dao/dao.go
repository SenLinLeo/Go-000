package dao

import (
	"database/sql"

	"github.com/pkg/errors"
)

func Dao() error {
	// if no error return nil
	// has error return sqlerror
	// return sql.ErrNoRows
	return errors.Wrap(sql.ErrNoRows, "Dao return error.")
}
