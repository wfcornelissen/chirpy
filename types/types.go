package types

import "sync/atomic"

type ApiConfig struct {
	FileserverHits atomic.Int32
}

type Chirp struct {
	Body         string `json:"body" default:""`
	Cleaned_body string `json:"cleaned_body" default:""`
}
