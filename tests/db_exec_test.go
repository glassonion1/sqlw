package sqlw_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/sqlw"
)

func TestDBExec(t *testing.T) {

	master := sqlw.Config{
		User:     "root",
		Password: "password",
		Port:     "3306",
		DBName:   "app",
	}
	db, err := sqlw.NewMySQLDB(master)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name     string
		in       sqlw.SQLMutation
		affected int64
	}{
		{
			name:     "Creates one data",
			in:       "INSERT INTO companies(id, name) VALUES('id_9999', 'foo')",
			affected: 1,
		},
		{
			name:     "Updates one data",
			in:       "UPDATE companies SET name='bar' WHERE id='id_0000'",
			affected: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := db.Exec(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
			got, err := res.RowsAffected()
			if err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(got, tt.affected); diff != "" {
				t.Errorf("failed test %s: %v", tt.name, diff)
			}
		})
	}
}
