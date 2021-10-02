package storage

import (
	"fmt"
	"time"
)

var ErrNoLink = fmt.Errorf("no link for the given ID found")

type Service interface {
	Save(string, time.Time) (string, error)
	Load(string) (string, error)
	LoadInfo(string) (*Item, error)
	Close() error
}

type Item struct {
	Id      uint64 `json:"id" redis:"id"`
	URL     string `json:"url" redis:"url"`
	Expires string `json:"expires" redis:"expires"`
	Visits  uint64 `json:"visits" redis:"visits"`
}
