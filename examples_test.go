package sqlw_test

import (
	"context"
	"log"

	"github.com/glassonion1/sqlw"
)

var m, r1, r2 sqlw.Config

func ExampleNewMySQLDB() {
	// Master
	m := sqlw.Config{
		User: "root", Password: "password",
		Host: "127.0.0.1", Port: "3306", DBName: "app",
	}
	// Replica1
	r1 := sqlw.Config{
		User: "root", Password: "password",
		Host: "127.0.0.1", Port: "3307", DBName: "app",
	}
	// Replica2
	r2 := sqlw.Config{
		User: "root", Password: "password",
		Host: "127.0.0.1", Port: "3308", DBName: "app",
	}
	// Connects to MySQL
	db, err := sqlw.NewMySQLDB(m, r1, r2)
	if err != nil {
		// TODO: Handle error.
	}
	_ = db // TODO: Use db.
}

func ExampleDB_Exec() {
	db, err := sqlw.NewMySQLDB(m, r1, r2)
	if err != nil {
		// TODO: Handle error.
	}

	ctx := context.Background()

	// Executes mutation query on the master database
	res, err := db.Exec(ctx, "INSERT INTO users(id, name) VALUES(?, ?)", "id:001", "hoge")
	if err != nil {
		// TODO: Handle error.
	}
	_ = res // TODO: Use res
}

func ExampleDB_Query() {
	db, err := sqlw.NewMySQLDB(m, r1, r2)
	if err != nil {
		// TODO: Handle error.
	}

	// Model
	type User struct {
		ID   string
		Name string
	}

	ctx := context.Background()

	// Executes query on the replica database
	rows, err := db.Query(ctx, "SELECT * FROM users WHERE name = ?", "hoge")
	if err != nil {
		// TODO: Handle error.
	}
	defer rows.Close()

	// Mapping data to model
	users := []User{}
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			// TODO: Handle error.
		}
		users = append(users, user)
	}
	log.Printf("users: %v", users)
}

func ExampleDB_Transaction() {
	db, err := sqlw.NewMySQLDB(m, r1, r2)
	if err != nil {
		// TODO: Handle error.
	}

	// Executes multiple queries in database transaction
	fn := func(ctx context.Context, tx *sqlw.Tx) error {
		_, err := tx.Exec(ctx, "INSER INTO users(id, name) VALUES(?, ?)", "id:001", "hoge")
		if err != nil {
			// Rollbacks automatically
			return err
		}
		_, err = tx.Exec(ctx, "UPDATE users SET name=? WHERE id=?", "piyo", "id:001")
		if err != nil {
			// Rollbacks automatically
			return err
		}

		// Warn: this query is executed outside of transaction
		_, _ = db.Exec(ctx, "UPDATE hoge SET name='foo'")

		return nil
	}

	ctx := context.Background()

	// Executes transaction and commits automatically if no errors
	if err := db.Transaction(ctx, fn); err != nil {
		// TODO: Handle error.
	}

	// Executes query for master database
	rows, err := db.QueryForMaster(ctx, "SELECT * FROM user")
	if err != nil {
		// TODO: Handle error.
	}
	_ = rows // TODO: Use rows
}
