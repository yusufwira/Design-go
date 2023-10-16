package connection

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dB struct {
	Db            *gorm.DB
	StorageClient *storage.Client
}

type Config struct {
	PGHost     string
	PGPort     string
	PGUser     string
	PGPassword string
	PGDB       string
}

func loadConfig() Config {
	return Config{
		PGHost:     os.Getenv("DB_HOST"),
		PGPort:     os.Getenv("DB_PORT"),
		PGUser:     os.Getenv("DB_USERNAME"),
		PGPassword: os.Getenv("DB_PASSWORD"),
		PGDB:       os.Getenv("DB_DATABASE"),
	}
}

func Database() *dB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	config := loadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", config.PGHost, config.PGPort, config.PGUser, config.PGPassword, config.PGDB)
	db, errs := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// dsn := "postgres://postgres:V3ry5tr0n94dm1nP@$$w0rd@192.168.188.232:5432/pi-smart"
	//dsn := "postgres://postgres:postgres@localhost:5432/postgres"
	// db, errs := gorm.Open(postgres.Open(dsn),
	// 	&gorm.Config{NowFunc: func() time.Time {
	// 		return time.Now().UTC()
	// 	},
	// 	})
	// db, errs := gorm.Open(postgres.Open("postgres://postgres:postgres@localhost:5432/postgres"), &gorm.Config{})
	if errs != nil {
		panic("failed to connect database")
	}

	// Initialize google storage client
	log.Printf("Connecting to Cloud Storage\n")
	ctx := context.Background()
	client, errss := storage.NewClient(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel() // release resources if slowOperation complete before timeout elapses
	// storage, errss := storage.NewClient(ctx)

	if errss != nil {
		panic("Error Creating cloud storage client")
	}

	return &dB{
		Db:            db,
		StorageClient: client,
	}
}
