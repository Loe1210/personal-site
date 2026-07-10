package configs

import "os"

func LoadMySQLConfig() MySQLConfig {
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = "root"
	}

	charset := os.Getenv("MYSQL_CHARSET")
	if charset == "" {
		charset = "utf8mb4"
	}

	return MySQLConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName:   os.Getenv("MYSQL_DBNAME"),
		Charset:  charset,
	}
}