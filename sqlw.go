package sqlw

import "fmt"

// Config holds the database configuration information.
type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func (c Config) mysqlStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (c Config) postgresStr() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}
