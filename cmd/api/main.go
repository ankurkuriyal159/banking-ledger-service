// cmd/api/main.go
package main

import (
	"log"
	"net/http"

	"github.com/ankurkuriyal159/banking-ledger-service/internal/api"
	"github.com/ankurkuriyal159/banking-ledger-service/internal/db"
	"github.com/ankurkuriyal159/banking-ledger-service/internal/queue"
	"github.com/gorilla/mux"
)

func main() {
	mysqlDB, err := db.InitMySQL()
	if err != nil {
		log.Fatalf("failed to init mysql: %v", err)
	}

	mongoDB, err := db.InitMongo()
	if err != nil {
		log.Fatalf("failed to init mongo: %v", err)
	}

	producer, err := queue.InitKafkaProducer()
	if err != nil {
		log.Fatalf("failed to init kafka producer: %v", err)
	}

	r := mux.NewRouter()
	api.RegisterRoutes(r, mysqlDB, mongoDB, producer)

	log.Println("API server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
