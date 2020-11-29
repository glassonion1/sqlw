package sqlw_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/sqlw"
)

func TestDBPrepareQuery(t *testing.T) {

	master := sqlw.Config{
		User:     "root",
		Password: "password",
		Port:     "3306",
		DBName:   "app",
	}
	rep1 := master
	rep1.Port = "3307"
	rep2 := master
	rep2.Port = "3308"
	db, err := sqlw.NewMySQLDB(master, rep1, rep2)
	if err != nil {
		t.Error(err)
	}

	type user struct {
		ID   string
		Name string
	}
	tests := []struct {
		name string
		in   sqlw.SQLQuery
		arg  string
		want []user
	}{
		{
			name: "Finds a data from the table",
			in:   "SELECT * FROM users WHERE ID = ?",
			arg:  "id_0000",
			want: []user{
				{
					ID:   "id_0000",
					Name: "hoge",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// No context
			{
				stmt, err := db.PrepareQuery(tt.in)
				if err != nil {
					t.Error(err)
				}
				defer stmt.Close()
				rows, err := stmt.Query(tt.arg)
				if err != nil {
					t.Error(err)
				}
				defer rows.Close()
				got := []user{}
				for rows.Next() {
					user := user{}
					if err := rows.Scan(&user.ID, &user.Name); err != nil {
						t.Error(err)
					}
					got = append(got, user)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("failed test %s: %v", tt.name, diff)
				}
			}

			// With context
			{
				stmt, err := db.PrepareQueryContext(context.Background(), tt.in)
				if err != nil {
					t.Error(err)
				}
				defer stmt.Close()
				rows, err := stmt.Query(tt.arg)
				if err != nil {
					t.Error(err)
				}
				defer rows.Close()
				got := []user{}
				for rows.Next() {
					user := user{}
					if err := rows.Scan(&user.ID, &user.Name); err != nil {
						t.Error(err)
					}
					got = append(got, user)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("failed test %s: %v", tt.name, diff)
				}
			}
		})
	}
}

func TestDBPrepareQueryForMaster(t *testing.T) {

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

	type user struct {
		ID   string
		Name string
	}
	tests := []struct {
		name string
		in   sqlw.SQLQuery
		arg  string
		want []user
	}{
		{
			name: "Finds a data from the table",
			in:   "SELECT * FROM users WHERE ID = ?",
			arg:  "id_0000",
			want: []user{
				{
					ID:   "id_0000",
					Name: "hoge",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// No context
			{
				stmt, err := db.PrepareQueryForMaster(tt.in)
				if err != nil {
					t.Error(err)
				}
				defer stmt.Close()
				rows, err := stmt.Query(tt.arg)
				if err != nil {
					t.Error(err)
				}
				defer rows.Close()
				got := []user{}
				for rows.Next() {
					user := user{}
					if err := rows.Scan(&user.ID, &user.Name); err != nil {
						t.Error(err)
					}
					got = append(got, user)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("failed test %s: %v", tt.name, diff)
				}
			}

			// With context
			{
				stmt, err := db.PrepareQueryContextForMaster(context.Background(), tt.in)
				if err != nil {
					t.Error(err)
				}
				defer stmt.Close()
				rows, err := stmt.Query(tt.arg)
				if err != nil {
					t.Error(err)
				}
				defer rows.Close()
				got := []user{}
				for rows.Next() {
					user := user{}
					if err := rows.Scan(&user.ID, &user.Name); err != nil {
						t.Error(err)
					}
					got = append(got, user)
				}
				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("failed test %s: %v", tt.name, diff)
				}
			}
		})
	}
}

func TestDBPrepareMutation(t *testing.T) {

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
		args     []interface{}
		affected int64
	}{
		{
			name:     "Creates one data",
			in:       "INSERT INTO companies(id, name) VALUES(?, ?)",
			args:     []interface{}{"id_0010", "foo"},
			affected: 1,
		},
		{
			name:     "Updates one data",
			in:       "UPDATE companies SET name = ? WHERE id = ?",
			args:     []interface{}{"XXXXXX", "id_0000"},
			affected: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stmt, err := db.PrepareMutation(tt.in)
			if err != nil {
				t.Error(err)
			}
			defer stmt.Close()
			res, err := stmt.Exec(tt.args...)
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

func TestDBPrepareMutationContext(t *testing.T) {

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
		args     []interface{}
		affected int64
	}{
		{
			name:     "Creates one data",
			in:       "INSERT INTO companies(id, name) VALUES(?, ?)",
			args:     []interface{}{"id_0020", "foo"},
			affected: 1,
		},
		{
			name:     "Updates one data",
			in:       "UPDATE companies SET name = ? WHERE id = ?",
			args:     []interface{}{"bar", "id_0000"},
			affected: 1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			stmt, err := db.PrepareMutationContext(context.Background(), tt.in)
			if err != nil {
				t.Error(err)
			}
			defer stmt.Close()
			res, err := stmt.Exec(tt.args...)
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
