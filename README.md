# sqlw

[![Godoc](https://img.shields.io/badge/godoc-reference-blue)](https://godoc.org/github.com/glassonion1/sqlw)
[![GitHub license](https://img.shields.io/github/license/glassonion1/sqlw)](https://github.com/glassonion1/sqlw/blob/main/LICENSE)

database/sql wrapper library for Go  
Manages automatically on the master/replica databases

## Install
```
$ go get github.com/glassonion1/sqlw
```

## Usage
### Database connection

Connects to database.
```go
// Settings for master
master := sqlw.Config{
  User: "root", Password: "password",
  Host: "127.0.0.1", Port: "3306", DBName: "app",
}
// Settings for replica1
rep1 := sqlw.Config{
  User: "root", Password: "password",
  Host: "127.0.0.1", Port: "3307", DBName: "app",
}
// Settings for replica2
rep2 := sqlw.Config{
  User: "root", Password: "password",
  Host: "127.0.0.1", Port: "3308", DBName: "app",
}
// Connects to mysql
db, err := sqlw.NewMySQLDB(master, rep1, rep2)
if err != nil {
  // TODO: Handle error.
}
```

To confirm the database connection.
```go
db, err := sqlw.NewMySQLDB(master, rep1, rep2)
if err != nil {
  // TODO: Handle error.
}
// Is it readable?
if err := db.Readable(); err != nil {
  // not readable
}
// Is it writable?
if err := db.Writable(); err != nil {
  // not writable
}
```

### Executes query

Query the database
```go
db, err := sqlw.NewMySQLDB(master, rep1, rep2)
if err != nil {
  // TODO: Handle error.
}

// table definition
type User struct {
  ID   string
  Name string
}

// Query the database(exec on replica)
rows, err := db.Query("SELECT * FROM users WHERE name = ?", "hoge")
if err != nil {
  // TODO: Handle error.
}
defer rows.Close()

// Scan for the selected data
users := []User{}
for rows.Next() {
  user := User{}
  if err := rows.Scan(&user.ID, &user.Name); err != nil {
    // TODO: Handle error.
  }
  users = append(users, user)
}
```

Query the database uses prepare method(exec on replica)
```go
// Instanciates statement object
stmt, err := db.PrepareQuery("SELECT * FROM users WHERE name = ?")
if err != nil {
  // TODO: Handle error.
}
defer stmt.Close()
// Executes query
rows, err := stmt.Query("hoge")
if err != nil {
  // TODO: Handle error.
}
defer rows.Close()

users := []User{}
for rows.Next() {
  user := User{}
  if err := rows.Scan(&user.ID, &user.Name); err != nil {
    // TODO: Handle error.
  }
  users = append(users, user)
}
```

Executes the mutation query(exec on master)
```go
db, err := sqlw.NewMySQLDB(master, rep1, rep2)
if err != nil {
  // TODO: Handle error.
}

res, err := db.Exec("INSERT INTO users(id, name) VALUES(?, ?)", "id:001", "hoge")
if err != nil {
  // TODO: Handle error.
}
```

Executes the mutation query uses prepare method(exec on master)
```go
// Instanciates statement object
stmt, err := db.PrepareMutation("INSERT INTO users(id, name) VALUES(?, ?)")
if err != nil {
  // TODO: Handle error.
}
defer stmt.Close()

res, err := stmt.Exec("id:001", "hoge")
if err != nil {
  // TODO: Handle error.
}
```

### Transaction

Automatically commit or rollback on transaction
```go
db, err := sqlw.NewMySQLDB(master, rep1, rep2)
if err != nil {
  // TODO: Handle error.
}
// Processes the transaction on the function
fn := func(tx *sqlw.Tx) error {
  _, err := tx.Exec("INSER INTO users(id, name) VALUES(?, ?)", "id:001", "hoge")
  if err != nil {
    // rollback on automatically
    return err
  }
  _, err := tx.Exec("UPDATE users SET name=? WHERE id=?", "piyo", "id:001")
  if err != nil {
    // rollback on automatically
    return err
  }
  return nil
}
// Executes transaction function
if err := db.Transaction(fn); err != nil {
  // TODO: Handle error.
}

// Query the master database
rows, err := db.QueryForMaster("SELECT * FROM user")
```

## Unit tests

Executes unit tests
```
$ cd mysql
$ docker-compose up -d
$ cd ../
$ go test -v ./...
```