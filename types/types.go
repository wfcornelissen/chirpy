package types

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Dbquery        *database.Queries
	Platform       string
}

type Chirp struct {
	Body         string `json:"body" default:""`
	Cleaned_body string `json:"cleaned_body" default:""`
}

type User struct {
	UserID    uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	EAddress  string    `json:"email"`
}
