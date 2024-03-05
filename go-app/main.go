package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type WebhookPayload struct {
	StrategyID string  `json:"strategy_id"`
	SecretKey  string  `json:"secret_key"`
	Action     string  `json:"action"`
	Price      float64 `json:"price"`
}

type Trade struct {
	StrategyID     string    `json:"strategy_id"`
	Action         string    `json:"action"`
	Quantity       int       `json:"quantity"` // Assume this is determined by your logic
	Direction      string    `json:"direction"`
	EntryTime      time.Time `json:"entry_time"`
	ExitTime       time.Time `json:"exit_time"`
	OpenPrice      float64   `json:"open_price"`
	ClosePrice     float64   `json:"close_price"`
	ProfitLoss     float64   `json:"profit_loss"`
	PnLPercent     float64   `json:"pnl_percent"`
	Status         string    `json:"status"`
	InitialBalance float64   `json:"initial_balance"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "host=db user=simulator password=password dbname=trading_simulator sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Fail to connect: %v", err)
	}

	http.HandleFunc("/webhook", webhookHandler)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if payload.SecretKey != "dU5TyY6ZEgiihmT4wdHGN3j7G5kbwS" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	trade := Trade{
		StrategyID: payload.StrategyID,
		Action:     payload.Action,
		OpenPrice:  payload.Price,
		EntryTime:  time.Now(),
		Status:     "open",
	}

	calculateTradeDetails(&trade)

	if err := saveTrade(trade); err != nil {
		http.Error(w, "Error saving trade to database", http.StatusInternalServerError)
		log.Printf("Error saving trade: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Trade processed"))
}

func calculateTradeDetails(trade *Trade) {
	var lastTrade Trade
	err := db.QueryRow("SELECT strategy_id, action, quantity, direction, entry_time, open_price FROM trades WHERE strategy_id = $1 AND status = 'open' ORDER BY entry_time DESC LIMIT 1", trade.StrategyID).Scan(&lastTrade.StrategyID, &lastTrade.Action, &lastTrade.Quantity, &lastTrade.Direction, &lastTrade.EntryTime, &lastTrade.OpenPrice)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error fetching last trade: %v\n", err)
		return
	}

	if trade.Action == "sell" && err != sql.ErrNoRows {
		trade.ClosePrice = trade.OpenPrice
		trade.ExitTime = time.Now()
		trade.Status = "closed"

		trade.ProfitLoss = float64(trade.Quantity) * (trade.ClosePrice - lastTrade.OpenPrice)
		trade.PnLPercent = (trade.ProfitLoss / (lastTrade.OpenPrice * float64(trade.Quantity))) * 100
	} else if trade.Action == "buy" {
		trade.EntryTime = time.Now()
		trade.Status = "open"
	}

}

func saveTrade(trade Trade) error {
	stmt, err := db.Prepare("INSERT INTO trades (strategy_id, action, quantity, direction, entry_time, exit_time, open_price, close_price, profit_loss, pnl_percent, status, initial_balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(trade.StrategyID, trade.Action, trade.Quantity, trade.Direction, trade.EntryTime, trade.ExitTime, trade.OpenPrice, trade.ClosePrice, trade.ProfitLoss, trade.PnLPercent, trade.Status, trade.InitialBalance)
	if err != nil {
		return err
	}

	return nil
}
