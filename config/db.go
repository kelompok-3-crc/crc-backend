package config

import (
	"fmt"
	"log"
	"ml-prediction/internal/app/model"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
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
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get *sql.DB from GORM: %v", err)
	}

	driver, err := migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	if err != nil {
		log.Fatalf("could not create postgres driver: %v", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	fmt.Println("Current working directory:", wd)
	migrationsPath := "file://" + filepath.Join(wd, "migrations")
	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		c.Postgres.PostgresqlDbname,
		driver,
	)
	if err != nil {
		log.Fatalf("migration init failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	} else {
		log.Println("migrations applied successfully or no changes.")
	}
	seedAdminUser(db)

	return db
}

// seedAdminUser creates an admin user if one doesn't exist yet
func seedAdminUser(db *gorm.DB) {
	const (
		adminName     = "Admin1"
		adminNIP      = "ADM001"
		adminRole     = "admin"
		adminPassword = "admin123" // Default password - can be changed later
	)

	// Check if admin user already exists
	var count int64
	db.Model(&model.User{}).Where("nip = ? AND role = ?", adminNIP, adminRole).Count(&count)
	if count > 0 {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return
	}

	// Create the admin user
	adminUser := model.User{
		Nama:     adminName,
		NIP:      adminNIP,
		Role:     adminRole,
		Password: string(hashedPassword),
	}

	result := db.Create(&adminUser)
	if result.Error != nil {
		log.Printf("Error creating admin user: %v", result.Error)
		return
	}

	log.Printf("Admin user created successfully with NIP: %s and password: %s", adminNIP, adminPassword)
}
