package main

import (
	"log"

	"github.com/ankurkuriyal159/banking-ledger-service/internal/db"
	"github.com/ankurkuriyal159/banking-ledger-service/internal/queue"
	"github.com/ankurkuriyal159/banking-ledger-service/internal/services"
)

func main() {
	// Init MySQL
	mysqlDB, err := db.InitMySQL()
	if err != nil {
		log.Fatalf("failed to init mysql: %v", err)
	}

	// Init Mongo
	mongoDB, err := db.InitMongo()
	if err != nil {
		log.Fatalf("failed to init mongo: %v", err)
	}

	// Init Kafka Consumer
	consumer, err := queue.InitKafkaConsumer()
	if err != nil {
		log.Fatalf("failed to init kafka consumer: %v", err)
	}

	// Create Transaction Processor
	processor := services.NewTransactionProcessor(mysqlDB, mongoDB, consumer)

	// Start Processor Loop
	processor.Start()
}
