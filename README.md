# Banking Ledger Service (Case Study)

A **sample backend service** in Go for banking ledger.  
This is a **case study project** focused on showing architecture and best practices in a simple way.

## Features
- Create accounts with initial balance (MySQL).
- Deposit & Withdraw funds (stored in DB)).
- Transaction log stored in MongoDB.
- Transaction log queued to MongoDB(Pending).
- All services run with `docker-compose`.

## Tech Stack
- **Go** – API + Processor(For kafka message processing)
- **MySQL** – Account balances
- **MongoDB** – Ledger (transaction history)
- **Kafka** – Queue for transactions
- **Docker Compose** – Service orchestration

## Run Locally
```bash
docker-compose up --build

## Database Schema

The service uses a simple schema for accounts:

```sql
CREATE TABLE accounts (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

**## API Requests Endpoints:**

**Header** Content-Type: application/json
**Create Account:**
POST http://localhost:8080/accounts
{
  "name": "Alice",
  "initial_balance": 1000
}


**Deposit Funds:**
POST http://localhost:8080/transactions/deposit
{
  "account_id": 1,
  "amount": 500
}

**Withdraw Funds:**
POST http://localhost:8080/transactions/withdraw
{
  "account_id": 1,
  "amount": 200
}

**Ledger(Transaction History):**
GET http://localhost:8080/accounts/1/ledger

