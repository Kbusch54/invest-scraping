package assets

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const EventTopic = "scraping.price"
const EventName = "assets.updated.price"

type Event struct {
	ID          string    `json:"id"`
	AssetName   string    `json:"asset_name"`
	AssetSymbol string    `json:"asset_symbol"`
	Price       float64   `json:"price"`
	Event       string    `json:"event"`
	Time        time.Time `json:"time"`
}

func NewEvent(assetName string, assetSymbol string, price float64, event string, createdAt time.Time) *Event {
	id := generateID(assetName, event, createdAt.String())
	return &Event{
		ID:          id,
		AssetName:   assetName,
		AssetSymbol: assetSymbol,
		Price:       price,
		Event:       event,
		Time:        createdAt,
	}
}

func generateID(assetName, event string, time string) string {
	// Concatenate the input strings
	combined := assetName + event + time

	// Compute the SHA-256 hash of the combined string
	hash := sha256.Sum256([]byte(combined))

	// Convert the hash to a hexadecimal string
	return hex.EncodeToString(hash[:])
}
