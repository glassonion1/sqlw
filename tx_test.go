package sqlw_test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/sqlw"
)

func TestTxQuery(t *testing.T) {

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
			name: "Finds a data from the table in the transaction",
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

			{
				got := []user{}
				fn := func(ctx context.Context, tx *sqlw.Tx) error {
					rows, err := tx.Query(ctx, tt.in)
					if err != nil {
						return err
					}
					defer rows.Close()
					for rows.Next() {
						user := user{}
						if err := rows.Scan(&user.ID, &user.Name); err != nil {
							t.Error(err)
						}
						got = append(got, user)
					}
					return nil
				}
				if err := db.Transaction(context.Background(), fn); err != nil {
					t.Error(err)
				}

				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("failed test %s: %v", tt.name, diff)
				}
			}
		})
	}
}

func TestTxQueryRow(t *testing.T) {

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
			name: "Finds any data from the table in the transaction",
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

			{
				got := user{}
				fn := func(ctx context.Context, tx *sqlw.Tx) error {
					row := tx.QueryRow(ctx, tt.in)
					if err := row.Err(); err != nil {
						t.Error(err)
					}
					if err := row.Scan(&got.ID, &got.Name); err != nil {
						t.Error(err)
					}
					return nil
				}
				if err := db.Transaction(context.Background(), fn); err != nil {
					t.Error(err)
				}

				if diff := cmp.Diff(got, tt.want); diff != "" {
					t.Errorf("failed test %s: %v", tt.name, diff)
				}
			}
		})
	}
}

func TestTxExec(t *testing.T) {

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

	type product struct {
		ID   string
		Name string
	}
	tests := []struct {
		name string
		in   []sqlw.SQLMutation
		want []product
		err  error
	}{
		{
			name: "Executes mutation queries that suceeds",
			in: []sqlw.SQLMutation{
				"INSERT INTO products(id, name) VALUES('id_0000', 'hoge')",
				"INSERT INTO products(id, name) VALUES('id_0001', 'hoge')",
				"UPDATE products SET name='fuga' WHERE id='id_0000'",
			},
			want: []product{
				{
					ID:   "id_0000",
					Name: "fuga",
				},
				{
					ID:   "id_0001",
					Name: "hoge",
				},
			},
			err: nil,
		},
		{
			name: "Executes mutation queries that failes",
			in: []sqlw.SQLMutation{
				"INSERT INTO products(id, name) VALUES('id_0002', 'piyo')",
				"UPDATE products SET name='piyo' WHERE id='id_0000'",
				"UPDATE product SET name='piyo'",
			},
			want: []product{
				{
					ID:   "id_0000",
					Name: "fuga",
				},
				{
					ID:   "id_0001",
					Name: "hoge",
				},
			},
			err: errors.New("failed to execcute transaction: Error 1146: Table 'app.product' doesn't exist"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fn := func(ctx context.Context, tx *sqlw.Tx) error {
				for _, sql := range tt.in {
					_, err := tx.Exec(ctx, sql)
					if err != nil {
						return err
					}
				}
				return nil
			}
			ctx := context.Background()
			if err := db.Transaction(ctx, fn); err != nil {
				log.Print(err)
			}

			rows, err := db.QueryForMaster(ctx, "SELECT * FROM products")
			if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			defer rows.Close()
			got := []product{}
			for rows.Next() {
				product := product{}
				if err := rows.Scan(&product.ID, &product.Name); err != nil {
					t.Error(err)
				}
				got = append(got, product)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("failed test %s: %v", tt.name, diff)
			}
		})
	}
}
