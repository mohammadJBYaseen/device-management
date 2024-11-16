package config

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func initDBConnection() *gorm.DB {
	var sqlDB *sql.DB
	var err error
	var db *gorm.DB
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		ApplicationProperties.Database.Host,
		ApplicationProperties.Database.Username,
		ApplicationProperties.Database.Password,
		ApplicationProperties.Database.Database,
		ApplicationProperties.Database.Port,
		ApplicationProperties.Database.SSLMode)
	if db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dns,
		PreferSimpleProtocol: true,
	}), &gorm.Config{}); err != nil {
		log.Fatalf("postgresql database connection exception:%v", err)
	} else {
		if sqlDB, err = db.DB(); err != nil {
			log.Fatalf("postgresql database connection exception:%v", err)
		}
		if err = sqlDB.Ping(); err != nil {
			log.Fatalf("postgresql database connection exception:%v", err)
		}
		sqlDB.SetMaxIdleConns(ApplicationProperties.Database.MaxIdleConn)
		sqlDB.SetMaxOpenConns(ApplicationProperties.Database.MaxOpenConn)
		sqlDB.SetConnMaxIdleTime(ApplicationProperties.Database.GetConnMaxIdleTime())
		sqlDB.SetMaxOpenConns(ApplicationProperties.Database.MaxOpenConn)
		log.Printf("postgresql database connection has been intialized...")
	}
	return db
}

func closeDBConnection() {
	if DB == nil {
		return
	}
	if sqlDB, err := DB.DB(); err != nil {
		log.Fatalf("postgresql database connection exception:%v", err)
	} else if err = sqlDB.Close(); err != nil {
		log.Fatalf("postgresql database connection exception:%v", err)
	}
}
