package pgstmt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/acoshift/pgsql/pgstmt"
)

func TestInsert(t *testing.T) {
	t.Parallel()

	t.Run("insert", func(t *testing.T) {
		q, args := pgstmt.Insert(func(b pgstmt.InsertStatement) {
			b.Into("users")
			b.Columns("username", "name", "created_at")
			b.Value("tester1", "Tester 1", pgstmt.Default)
			b.Value("tester2", "Tester 2", "now()")
			b.OnConflict("username").DoNothing()
			b.Returning("id", "name")
		}).SQL()

		assert.Equal(t,
			"insert into users (username, name, created_at) values ($1, $2, default), ($3, $4, $5) on conflict (username) do nothing returning id, name",
			q,
		)
		assert.EqualValues(t,
			[]interface{}{
				"tester1", "Tester 1",
				"tester2", "Tester 2", "now()",
			},
			args,
		)
	})

	t.Run("insert select", func(t *testing.T) {
		q, args := pgstmt.Insert(func(b pgstmt.InsertStatement) {
			b.Into("films")
			b.Select(func(b pgstmt.SelectStatement) {
				b.Columns("*")
				b.From("tmp_films")
				b.Where(func(b pgstmt.Cond) {
					b.LtRaw("date_prod", "2004-05-07")
				})
			})
		}).SQL()

		assert.Equal(t,
			"insert into films select * from tmp_films where (date_prod < 2004-05-07)",
			q,
		)
		assert.Empty(t, args)
	})

	t.Run("insert on conflict do update", func(t *testing.T) {
		q, args := pgstmt.Insert(func(b pgstmt.InsertStatement) {
			b.Into("users")
			b.Columns("username", "email")
			b.Value("tester1", "tester1@localhost")
			b.OnConflict("username").DoUpdate(func(b pgstmt.UpdateStatement) {
				b.Set("email").ToRaw("excluded.email")
				b.Set("updated_at").ToRaw("now()")
			})
			b.Returning("id")
		}).SQL()

		assert.Equal(t,
			stripSpace(`
				insert into users (username, email)
				values ($1, $2)
				on conflict (username) do update
				set email = excluded.email,
					updated_at = now()
				returning id
			`),
			q,
		)
		assert.EqualValues(t,
			[]interface{}{
				"tester1", "tester1@localhost",
			},
			args,
		)
	})
}
