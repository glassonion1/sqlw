package sqlw_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//nolint
func TestMain(m *testing.M) {
	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		log.Fatal(err)
	}

	if err := setup(master); err != nil {
		master.Close()
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		cleanup(db)
		db.Close()
	}(master)

	m.Run()
}

func setup(db *sql.DB) error {

	if err := cleanup(db); err != nil {
		return err
	}

	// creates tables users, companies, products
	// inserts 2 test data

	ddl1 := `CREATE TABLE IF NOT EXISTS users(
                id varchar(255),    
                name varchar(255))`
	if _, err := db.Exec(ddl1); err != nil {
		return err
	}

	dml1 := "INSERT INTO users(id, name) VALUES('id_0000', 'hoge')"
	if _, err := db.Exec(dml1); err != nil {
		return err
	}

	ddl2 := `CREATE TABLE IF NOT EXISTS companies(
                id varchar(255),    
                name varchar(255))`
	if _, err := db.Exec(ddl2); err != nil {
		return err
	}
	dml2 := "INSERT INTO companies(id, name) VALUES('id_0000', 'hoge')"
	if _, err := db.Exec(dml2); err != nil {
		return err
	}

	ddl3 := `CREATE TABLE IF NOT EXISTS products(
                id varchar(255),    
                name varchar(255))`
	if _, err := db.Exec(ddl3); err != nil {
		return err
	}

	log.Println("ready for test")

	// Replica lag measures
	time.Sleep(time.Second * 1)

	return nil
}

func cleanup(db *sql.DB) error {
	ddl1 := "DROP TABLE IF EXISTS users, companies, products"
	if _, err := db.Exec(ddl1); err != nil {
		return err
	}

	return nil
}
