package types

import "sync/atomic"

type ApiConfig struct {
	FileserverHits atomic.Int32
}

type Chirp struct {
	Body string `json:"body" default:""`
	Valid bool `json:"valid" default:"false"`
}