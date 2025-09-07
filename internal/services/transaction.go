package services

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"

	"github.com/ankurkuriyal159/banking-ledger-service/internal/queue"
)

type TransactionProcessor struct {
	mysqlDB  *gorm.DB
	mongoDB  *mongo.Database
	consumer *queue.KafkaConsumer
}

func NewTransactionProcessor(mysqlDB *gorm.DB, mongoDB *mongo.Database, consumer *queue.KafkaConsumer) *TransactionProcessor {
	return &TransactionProcessor{mysqlDB: mysqlDB, mongoDB: mongoDB, consumer: consumer}
}

func (tp *TransactionProcessor) Start() {
	fmt.Println("Transaction processor started")
	// Here you start reading from Kafka and updating DBs
}
