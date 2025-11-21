package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/setiadijoe/tibd-rate-limiter/config"
	"github.com/setiadijoe/tibd-rate-limiter/handlers"
	"github.com/setiadijoe/tibd-rate-limiter/ratelimit"
)

func openDB() (*sql.DB, error) {
	cfg, err := config.LoadConfig("config/files/default.yaml")
	if err != nil {
		return nil, err
	}

	dbCfg := cfg.Database

	// TiDB Cloud Serverless biasanya butuh TLS
	// pastikan DSN-nya pakai tls=true dan parseTime=true
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?tls=true&parseTime=true&charset=utf8mb4",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db err: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db err: %v", err)
	}

	return db, nil
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("error connection db: %+v", err)
	}
	rlSvc := ratelimit.New(db)
	handler := handlers.New(rlSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	log.Println("listening :8080")
	if err := http.ListenAndServe(":8080", handler.RateLimitMiddleware(mux)); err != nil {
		log.Fatalf("error starting http server: %+v", err)
	}
}
