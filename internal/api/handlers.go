package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ankurkuriyal159/banking-ledger-service/internal/models"
	"github.com/ankurkuriyal159/banking-ledger-service/internal/queue"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Handler struct {
	DB       *gorm.DB
	MongoDB  *mongo.Database
	Producer *queue.KafkaProducer
}

func NewHandlers(db *gorm.DB, mongoDB *mongo.Database, producer *queue.KafkaProducer) *Handler {
	return &Handler{
		DB:       db,
		MongoDB:  mongoDB,
		Producer: producer,
	}
}

type CreateAccountRequest struct {
	Name           string  `json:"name"`
	InitialBalance float64 `json:"initial_balance"`
}

type TransactionRequest struct {
	AccountID uint    `json:"account_id"`
	Amount    float64 `json:"amount"`
}

// ---- Handlers ----

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	account := models.Account{
		Name:    req.Name,
		Balance: req.InitialBalance,
	}

	if err := h.DB.Create(&account).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(account)
}

func (h *Handler) DepositFunds(w http.ResponseWriter, r *http.Request) {
	h.processTransaction(w, r, "deposit")
}

func (h *Handler) WithdrawFunds(w http.ResponseWriter, r *http.Request) {
	h.processTransaction(w, r, "withdraw")
}

func (h *Handler) GetAccountLedger(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	cursor, err := h.MongoDB.Collection("transactions").
		Find(context.Background(), bson.M{"account_id": accountID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var txs []models.Transaction
	if err := cursor.All(context.Background(), &txs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(txs)
}

// ---- helper ----

func (h *Handler) processTransaction(w http.ResponseWriter, r *http.Request, txType string) {
	var req TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var account models.Account
	if err := h.DB.First(&account, req.AccountID).Error; err != nil {
		http.Error(w, "account not found", http.StatusNotFound)
		return
	}

	if txType == "withdraw" && account.Balance < req.Amount {
		http.Error(w, "insufficient funds", http.StatusBadRequest)
		return
	}

	// update MySQL balance
	if txType == "deposit" {
		account.Balance += req.Amount
	} else {
		account.Balance -= req.Amount
	}

	if err := h.DB.Save(&account).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// insert transaction log in MongoDB
	tx := models.Transaction{
		AccountID: req.AccountID,
		Type:      txType,
		Amount:    req.Amount,
		Timestamp: time.Now(),
	}
	_, err := h.MongoDB.Collection("transactions").
		InsertOne(context.Background(), tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "transaction successful",
		"accountId": account.ID,
		"balance":   account.Balance,
	})
}

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
