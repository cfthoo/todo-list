package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func ConnectDatabase() (*sql.DB, error) {

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	dbport, _ := strconv.Atoi(port)

	connInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, dbname)

	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Established a successful DB connection!")
	return db, nil
}

func RunDBMigration(migrationURL string, dbString string) {
	migration, err := migrate.New(migrationURL, dbString)
	if err != nil {
		log.Fatal("Failed to create new migrate instance:", err)
	}

	err = migration.Up()
	if err != nil {
		if err.Error() == "no change" {
			log.Println("No changes on migration.")
		} else {
			log.Fatal("Failed to run migrate up:", err)
		}

	}
	log.Println("Successfully migrated DB")
}
