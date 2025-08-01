package types

import (
	"sync/atomic"

	"github.com/wfcornelissen/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	dbquery        *database.Queries
}

type Chirp struct {
	Body         string `json:"body" default:""`
	Cleaned_body string `json:"cleaned_body" default:""`
}
