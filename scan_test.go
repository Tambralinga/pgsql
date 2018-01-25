package pgsql_test

import (
	"testing"
	"time"

	"github.com/acoshift/pgsql"
)

func TestScan(t *testing.T) {
	db := open(t)
	defer db.Close()

	_, err := db.Exec(`
		drop table if exists test_pgsql_scan;
		create table test_pgsql_scan (
			id int primary key,
			name varchar,
			value int,
			created_at timestamp not null default now()
		);
		insert into test_pgsql_scan (
			id, name, value
		) values
			(1, 'name 1', 0),
			(2, 'name 2', 1),
			(3, 'name 3', 123);
	`)
	if err != nil {
		t.Fatalf("prepare table error; %v", err)
	}
	defer db.Exec(`drop table test_pgsql_scan`)

	type scanStruct struct {
		ID        int
		Name      string
		Value     int
		CreatedAt time.Time
	}

	row := db.QueryRow(`select * from test_pgsql_scan order by id limit 1`)
	var s scanStruct
	err = pgsql.ScanRow(row, &s)
	if err != nil {
		t.Error("ScanRow error")
	}
	if s.ID != 1 {
		t.Errorf("expected scanned id is 1, got %d", s.ID)
	}
	if s.Name != "name 1" {
		t.Errorf("expected scanned name is 'name 1', got %s", s.Name)
	}
	if s.Value != 0 {
		t.Errorf("expected scanned value is 0, got %d", s.Value)
	}
	if s.CreatedAt.IsZero() {
		t.Errorf("expected scanned created_at is not zero, got %v", s.CreatedAt)
	}
}
