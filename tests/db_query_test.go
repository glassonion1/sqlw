package sqlw_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/sqlw"
)

func TestDBQuery(t *testing.T) {

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
		want []user
	}{
		{
			name: "Find a data from the table",
			in:   "SELECT * FROM users",
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

			rows, err := db.Query(context.Background(), tt.in)
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
		})
	}
}

func TestDBQueryForMaster(t *testing.T) {

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
		want []user
	}{
		{
			name: "Find any data from the table",
			in:   "SELECT * FROM users",
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

			rows, err := db.QueryForMaster(context.Background(), tt.in)
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
		})
	}
}

func TestDBQueryRow(t *testing.T) {

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
		want user
	}{
		{
			name: "Find any data from the table",
			in:   "SELECT * FROM users",
			want: user{
				ID:   "id_0000",
				Name: "hoge",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			row := db.QueryRow(context.Background(), tt.in)
			if err := row.Err(); err != nil {
				t.Error(err)
			}
			got := user{}
			if err := row.Scan(&got.ID, &got.Name); err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("failed test %s: %v", tt.name, diff)
			}
		})
	}
}

func TestDBQueryRowForMaster(t *testing.T) {

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
		want user
	}{
		{
			name: "Find any data from the table",
			in:   "SELECT * FROM users",
			want: user{
				ID:   "id_0000",
				Name: "hoge",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			row := db.QueryRowForMaster(context.Background(), tt.in)
			if err := row.Err(); err != nil {
				t.Error(err)
			}
			got := user{}
			if err := row.Scan(&got.ID, &got.Name); err != nil {
				t.Error(err)
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("failed test %s: %v", tt.name, diff)
			}
		})
	}
}
