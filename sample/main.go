package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/glassonion1/sqlw"
	"github.com/glassonion1/sqlw/sample/port/repository"
	"github.com/glassonion1/sqlw/sample/port/rest"
	"github.com/glassonion1/sqlw/sample/usecase/interactor"
)

func main() {

	time.Local = nil

	setup()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Master
	master := sqlw.Config{
		User: "root", Password: "password",
		Host: "127.0.0.1", Port: "3306", DBName: "app",
	}
	// Replica1
	rep1 := sqlw.Config{
		User: "root", Password: "password",
		Host: "127.0.0.1", Port: "3307", DBName: "app",
	}
	// Replica2
	rep2 := sqlw.Config{
		User: "root", Password: "password",
		Host: "127.0.0.1", Port: "3308", DBName: "app",
	}
	db, err := sqlw.NewMySQLDB(master, rep1, rep2)
	if err != nil {
		panic(err)
	}

	repo := repository.NewItem(db)
	ii := interactor.NewItem(repo)
	h := rest.NewItemHandler(ii)

	e.GET("/items", h.List())
	e.GET("/items/:item_id", h.Get())
	e.POST("/items", h.Create())

	e.Start(":8080")
}

// Creates sample data
func setup() {
	master, err := sql.Open("mysql", "root:password@tcp(:3306)/app")
	if err != nil {
		panic(err)
	}
	defer master.Close()
	ddl1 := `CREATE TABLE IF NOT EXISTS items(
                  id varchar(255),    
                  name varchar(255))`
	if _, err := master.Exec(ddl1); err != nil {
		panic(err)
	}
}
