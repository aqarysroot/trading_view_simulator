package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type WebhookPayload struct {
	StrategyID string  `json:"strategy_id"`
	SecretKey  string  `json:"secret_key"`
	Action     string  `json:"action"`
	Price      float64 `json:"price"`
}

type Trade struct {
	StrategyID string
	Action     string
	OpenPrice  float64
	ClosePrice float64
	EntryTime  time.Time
	ExitTime   time.Time
	ProfitLoss float64
}

type Strategy struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Symbol           string  `json:"symbol"`
	InitialBalance   float64 `json:"initial_balance"`
	Lots             int     `json:"lots"`
	EquityPercent    float64 `json:"equity_percent"`
	Inverse          bool    `json:"inverse"`
	Pyramid          bool    `json:"pyramid"`
	MarketDataSource string  `json:"market_data_source"`
}

var db *sql.DB

func main() {
	var err error
	connStr := "host=db user=simulator password=password dbname=trading_simulator sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Fail to connect: %v", err)
	}

	http.HandleFunc("/create_strategy", createStrategyHandler)
	http.HandleFunc("/list_strategies", listStrategiesHandler)
	http.HandleFunc("/webhook", webhookHandler)
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createStrategyHandler(w http.ResponseWriter, r *http.Request) {
	var s Strategy
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.ID = uuid.New().String()

	_, err := db.Exec("INSERT INTO trading_strategy (id, name, symbol, initial_balance, lots, equity_percent, inverse, pyramid, market_data_source) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		s.ID, s.Name, s.Symbol, s.InitialBalance, s.Lots, s.EquityPercent, s.Inverse, s.Pyramid, s.MarketDataSource)
	if err != nil {
		http.Error(w, "Failed to insert strategy", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(s)
}

func listStrategiesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, symbol, initial_balance, lots, equity_percent, inverse, pyramid, market_data_source FROM trading_strategy")
	if err != nil {
		http.Error(w, "Failed to query strategies", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var strategies []Strategy
	for rows.Next() {
		var s Strategy
		if err := rows.Scan(&s.ID, &s.Name, &s.Symbol, &s.InitialBalance, &s.Lots, &s.EquityPercent, &s.Inverse, &s.Pyramid, &s.MarketDataSource); err != nil {
			http.Error(w, "Failed to scan strategy", http.StatusInternalServerError)
			return
		}
		strategies = append(strategies, s)
	}

	json.NewEncoder(w).Encode(strategies)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if payload.SecretKey != "X$#Ksecret" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	trade := calculateTradeDetails(payload)

	if err := saveTrade(trade); err != nil {
		http.Error(w, "Error saving trade", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Trade processed"))
}

func calculateTradeDetails(payload WebhookPayload) Trade {
	var openPrice, closePrice float64
	entryTime := time.Now()
	exitTime := entryTime.Add(time.Hour * 24)

	if payload.Action == "buy" {
		openPrice = payload.Price
		closePrice = openPrice * 1.05
	} else {
		openPrice = payload.Price
		closePrice = openPrice * 0.95
	}

	pnl := (closePrice - openPrice) * 100

	return Trade{
		StrategyID: payload.StrategyID,
		Action:     payload.Action,
		OpenPrice:  openPrice,
		ClosePrice: closePrice,
		EntryTime:  entryTime,
		ExitTime:   exitTime,
		ProfitLoss: pnl,
	}
}

func saveTrade(trade Trade) error {

	stmt, err := db.Prepare("INSERT INTO trades (strategy_id, action, open_price, close_price, entry_time, exit_time, pnl) VALUES ($1, $2, $3, $4, $5, $6, $7)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(trade.StrategyID, trade.Action, trade.OpenPrice, trade.ClosePrice, trade.EntryTime, trade.ExitTime, trade.ProfitLoss)
	if err != nil {
		return err
	}

	return nil
}
