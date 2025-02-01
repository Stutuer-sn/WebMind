// Package utils 负责数据库的连接和操作Postgresql
package utils

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectPostgresql 连接Postgresql数据库
func ConnectPostgresql() (*gorm.DB, error) {
	dsn := "user=postgres password=123456 dbname=webmind host=localhost port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}
