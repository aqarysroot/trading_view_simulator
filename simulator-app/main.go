package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type WebhookPayload struct {
	StrategyID string  `json:"strategy_id"`
	SecretKey  string  `json:"secret_key"`
	Action     string  `json:"action"`
	Price      float64 `json:"price"`
}

func randomAction() string {
	actions := []string{"buy", "sell"}
	return actions[rand.Intn(len(actions))]
}

func randomPrice() float64 {
	return rand.Float64()*10 + 100
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		payload := WebhookPayload{
			StrategyID: "BIOL-LH-10m",
			SecretKey:  "dU5TyY6ZEgiihmT4wdHGN3j7G5kbwS",
			Action:     randomAction(),
			Price:      randomPrice(),
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}

		resp, err := http.Post("http://goapp:8080/webhook", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		println("Webhook sent, response status:", resp.Status)
	}
}
