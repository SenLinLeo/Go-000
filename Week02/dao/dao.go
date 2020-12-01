package dao

import "database/sql"

func Dao() error {
	// if no error return nil
	// has error return sqlerror
	return sql.ErrNoRows
}
