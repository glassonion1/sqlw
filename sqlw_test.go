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
	setup()

	m.Run()
}

func setup() error {
	if err := cleanup(); err != nil {
		return err
	}

	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		return err
	}
	defer master.Close()

	// Creates three tables users, companies, products
	// Creates two test data

	ddl1 := `CREATE TABLE IF NOT EXISTS users(
                id varchar(255),    
                name varchar(255))`
	if _, err := master.Exec(ddl1); err != nil {
		return err
	}

	dml1 := "INSERT INTO users(id, name) VALUES('id_0000', 'hoge')"
	if _, err := master.Exec(dml1); err != nil {
		return err
	}

	ddl2 := `CREATE TABLE IF NOT EXISTS companies(
                id varchar(255),    
                name varchar(255))`
	if _, err := master.Exec(ddl2); err != nil {
		return err
	}
	dml2 := "INSERT INTO companies(id, name) VALUES('id_0000', 'hoge')"
	if _, err := master.Exec(dml2); err != nil {
		return err
	}

	ddl3 := `CREATE TABLE IF NOT EXISTS products(
                id varchar(255),    
                name varchar(255))`
	if _, err := master.Exec(ddl3); err != nil {
		return err
	}

	log.Println("ready for test")

	// For replica rag
	time.Sleep(time.Second * 1)

	return nil
}

func cleanup() error {
	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		return err
	}
	defer master.Close()

	ddl := "DROP TABLE IF EXISTS users, companies, products"
	if _, err := master.Exec(ddl); err != nil {
		return err
	}
	return nil
}
