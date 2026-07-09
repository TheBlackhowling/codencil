package db

import (
	"fmt"

	"github.com/TheBlackHowling/typrow"
	_ "github.com/lib/pq"
)

// Open connects to Postgres using database/sql.
func Open(databaseURL string) (*typrow.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database URL is required")
	}
	return typrow.Open("postgres", databaseURL)
}
