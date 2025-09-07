// internal/api/routes.go
package api

import (
	"github.com/ankurkuriyal159/banking-ledger-service/internal/queue"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func RegisterRoutes(r *mux.Router, db *gorm.DB, mongoDB *mongo.Database, producer *queue.KafkaProducer) {
	h := NewHandlers(db, mongoDB, producer)

	r.HandleFunc("/health", h.HealthHandler).Methods("GET")
	r.HandleFunc("/accounts", h.CreateAccount).Methods("POST")
	r.HandleFunc("/transactions/deposit", h.DepositFunds).Methods("POST")
	r.HandleFunc("/transactions/withdraw", h.WithdrawFunds).Methods("POST")
	r.HandleFunc("/accounts/{id}/ledger", h.GetAccountLedger).Methods("GET")
}
