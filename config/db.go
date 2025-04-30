package config

import (
	"fmt"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres" // or mysql, sqlite, etc.
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(c *Configuration) *gorm.DB {
	dataSourceName := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.Postgres.PostgresqlHost,
		c.Postgres.PostgresqlUser,
		c.Postgres.PostgresqlDbname,
		c.Postgres.PostgresqlPassword,
		c.Postgres.PostgresqlPort,
		c.Postgres.PostgresParams,
	)

	db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	panic("failed to get db instance: " + err.Error())
	// }

	// sqlDB.SetMaxOpenConns(6)
	// sqlDB.SetConnMaxLifetime(time.Hour)
	// sqlDB.SetConnMaxIdleTime(time.Minute * 30)
	// sqlDB.SetMaxIdleConns(6)

	return db
}
