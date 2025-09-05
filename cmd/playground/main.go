package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

func connectToDatabase() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("MAIN_DATABASE_DSN"))

	if err != nil {
		slog.Error("Error connecting to database: " + err.Error())
	}

	return db
}

func bootstrap() (*http.ServeMux, *ServiceMap, *sql.DB) {
	// load environmental variables
	slog.Info("Loading environment")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	slog.Info("Connecting to database")
	// get database
	db := connectToDatabase()

	slog.Info("Initializing services")
	// initialize services
	serviceMap := BuildServiceMap(db)

	// initialize controllers
	slog.Info("Initializing controllers")
	rootMux := RegisterHandlers(*serviceMap)

	return rootMux, serviceMap, db
}

func main() {
	// setup logger
	w := os.Stderr
	logger := slog.New(tint.NewHandler(w, &tint.Options{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	_, serviceMap, db := bootstrap()

	slog.Info("Listening on port 3000")
	// http.ListenAndServe(":3000", rootMux)

	playground(serviceMap, db)
}

func playground(serviceMap *ServiceMap, db *sql.DB) {

}
