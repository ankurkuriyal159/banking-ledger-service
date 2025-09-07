package models

import "time"

type Transaction struct {
	AccountID uint      `bson:"account_id"`
	Type      string    `bson:"type"` // deposit or withdraw
	Amount    float64   `bson:"amount"`
	Timestamp time.Time `bson:"timestamp"`
}
