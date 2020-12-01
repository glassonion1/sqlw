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
	teardown := setup()
	defer teardown()

	m.Run()

	// パニックが発生するとteardownが実行されないのでひろう
	if err := recover(); err != nil {
		log.Println(err)
	}
}

func setup() func() {
	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		panic(err)
	}
	defer master.Close()

	// テストテーブル3つ users, companies, products
	// テストデータ2件

	ddl1 := `CREATE TABLE IF NOT EXISTS users(
                id varchar(255),    
                name varchar(255))`
	if _, err := master.Exec(ddl1); err != nil {
		panic(err)
	}

	dml1 := "INSERT INTO users(id, name) VALUES('id_0000', 'hoge')"
	if _, err := master.Exec(dml1); err != nil {
		panic(err)
	}

	ddl2 := `CREATE TABLE IF NOT EXISTS companies(
                id varchar(255),    
                name varchar(255))`
	if _, err := master.Exec(ddl2); err != nil {
		panic(err)
	}
	dml2 := "INSERT INTO companies(id, name) VALUES('id_0000', 'hoge')"
	if _, err := master.Exec(dml2); err != nil {
		panic(err)
	}

	ddl3 := `CREATE TABLE IF NOT EXISTS products(
                id varchar(255),    
                name varchar(255))`
	if _, err := master.Exec(ddl3); err != nil {
		panic(err)
	}

	log.Println("ready for test")

	// スレーブ遅延対策
	time.Sleep(time.Second * 1)

	return func() {
		master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
		if err != nil {
			log.Fatal(err)
		}
		defer master.Close()
		if err != nil {
			log.Fatal(err)
		}
		ddl1 := "DROP TABLE users"
		if _, err := master.Exec(ddl1); err != nil {
			log.Fatal(err)
		}
		ddl2 := "DROP TABLE companies"
		if _, err := master.Exec(ddl2); err != nil {
			log.Fatal(err)
		}
		ddl3 := "DROP TABLE products"
		if _, err := master.Exec(ddl3); err != nil {
			log.Fatal(err)
		}
	}
}
