package pgsql

import "database/sql"

// Scan scans *sql.Rows into struct
func Scan(s *sql.Rows, dest interface{}) error {
	return nil
}
