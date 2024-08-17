package testhelpers

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CleanDatabase(t *testing.T, ctx context.Context, connStr string) {
	t.Helper()

	t.Cleanup(func() {
		connpool, err := pgxpool.New(ctx, connStr)
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}
		defer connpool.Close() // Close the connection pool to release the

		stmt := `
    DO $$
      DECLARE
        table_name text;
      BEGIN
        FOR table_name in (SELECT tablename FROM pg_tables WHERE schemaname='public') LOOP
          EXECUTE 'TRUNCATE TABLE ' || table_name || ' RESTART IDENTITY CASCADE;';
        END LOOP;
    END $$;
    `

		_, err = connpool.Exec(ctx, stmt)
		if err != nil {
			t.Fatalf("Error truncating tables: %v", err)
		}
	})
}
