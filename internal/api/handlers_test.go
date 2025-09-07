package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/ankurkuriyal159/banking-ledger-service/internal/api"
	"github.com/ankurkuriyal159/banking-ledger-service/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTest creates a fresh in-memory DB and handler for each test
func setupTest(t *testing.T) (*api.Handler, *mux.Router) {
	// Open in-memory SQLite without cgo
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&models.Account{}, &models.Transaction{}); err != nil {
		t.Fatal(err)
	}

	handler := &api.Handler{DB: db} // Mongo/Kafka can be nil in tests
	r := mux.NewRouter()
	api.RegisterRoutes(r, db, nil)

	return handler, r
}

func TestCreateAccount(t *testing.T) {
	_, r := setupTest(t)

	reqBody := []byte(`{"name":"Alice","initial_balance":500}`)
	req := httptest.NewRequest("POST", "/accounts", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var acc models.Account
	err := json.Unmarshal(w.Body.Bytes(), &acc)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", acc.Name)
	assert.Equal(t, 500.0, acc.Balance)
}

func TestDepositFunds(t *testing.T) {
	_, r := setupTest(t)

	// create account first
	accReq := []byte(`{"name":"Bob","initial_balance":200}`)
	accReqObj := httptest.NewRequest("POST", "/accounts", bytes.NewReader(accReq))
	accReqObj.Header.Set("Content-Type", "application/json")
	accW := httptest.NewRecorder()
	r.ServeHTTP(accW, accReqObj)

	var acc models.Account
	_ = json.Unmarshal(accW.Body.Bytes(), &acc)

	// deposit
	depositReq := []byte(`{"account_id":1,"amount":300}`)
	req := httptest.NewRequest("POST", "/transactions/deposit", bytes.NewReader(depositReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// verify balance updated
	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(500), resp["balance"])
}

func TestWithdrawFunds(t *testing.T) {
	_, r := setupTest(t)

	// create account first
	accReq := []byte(`{"name":"Charlie","initial_balance":400}`)
	accReqObj := httptest.NewRequest("POST", "/accounts", bytes.NewReader(accReq))
	accReqObj.Header.Set("Content-Type", "application/json")
	accW := httptest.NewRecorder()
	r.ServeHTTP(accW, accReqObj)

	// withdraw
	withdrawReq := []byte(`{"account_id":1,"amount":150}`)
	req := httptest.NewRequest("POST", "/transactions/withdraw", bytes.NewReader(withdrawReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(250), resp["balance"])
}

func TestGetAccountLedger(t *testing.T) {
	_, r := setupTest(t)

	// create account
	accReq := []byte(`{"name":"Dave","initial_balance":100}`)
	accReqObj := httptest.NewRequest("POST", "/accounts", bytes.NewReader(accReq))
	accReqObj.Header.Set("Content-Type", "application/json")
	accW := httptest.NewRecorder()
	r.ServeHTTP(accW, accReqObj)

	// deposit
	depositReq := []byte(`{"account_id":1,"amount":50}`)
	depReqObj := httptest.NewRequest("POST", "/transactions/deposit", bytes.NewReader(depositReq))
	depReqObj.Header.Set("Content-Type", "application/json")
	depW := httptest.NewRecorder()
	r.ServeHTTP(depW, depReqObj)

	// GET ledger
	req := httptest.NewRequest("GET", "/accounts/1/ledger", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var transactions []models.Transaction
	err := json.Unmarshal(w.Body.Bytes(), &transactions)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(transactions), 1)
}
