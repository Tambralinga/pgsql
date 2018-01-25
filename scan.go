package pgsql

import "database/sql"

// Scan scans *sql.Rows into struct
func Scan(rows *sql.Rows, dest interface{}) error {
	return nil
}

// ScanRow scans *sql.Row into struct
func ScanRow(row *sql.Row, dest interface{}) error {
	return nil
}
