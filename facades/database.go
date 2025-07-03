package facades

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	SqlDB *sql.DB
	once  sync.Once
)

func ConnectDB(envFiles ...string) *gorm.DB {
	once.Do(func() {
		if len(envFiles) > 0 {
			err := godotenv.Load(envFiles...)
			if err != nil {
				log.Printf("Warning: No .env file found. Using environment variables instead. Error: %v", err)
			} else {
				log.Println(".env file loaded successfully")
			}
		}

		requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_DB", "DB_USER"}
		for _, v := range requiredEnvVars {
			if os.Getenv(v) == "" {
				log.Fatalf("Error: Required environment variable %s is not set", v)
			}
		}

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DB"),
		)

		var err error
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					SlowThreshold: 500 * time.Millisecond,
					LogLevel:      logger.Warn,
					Colorful:      true,
				},
			),
			PrepareStmt: true,
		})
		if err != nil {
			log.Fatalf("Error: failed to connect to the database: %v", err)
		}

		SqlDB, err = DB.DB()
		if err != nil {
			log.Fatalf("Error: failed to get SQL DB object from GORM: %v", err)
		}

		SqlDB.SetMaxIdleConns(10)
		SqlDB.SetMaxOpenConns(200)
		SqlDB.SetConnMaxLifetime(15 * time.Minute)
		SqlDB.SetConnMaxIdleTime(5 * time.Minute)

		log.Println("Database connection successfully established")
	})

	return DB
}

func CloseDB() {
	if SqlDB != nil {
		err := SqlDB.Close()
		if err != nil {
			log.Printf("Error closing the database connection: %v", err)
		} else {
			log.Println("Database connection successfully closed")
		}
	}
}
