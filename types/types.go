package types

import (
	"database/sql"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Dbquery        *database.Queries
	Platform       string
	Secret         string
}

type Chirp struct {
	Id   uuid.UUID `json:"id"`
	Body string    `json:"body" default:""`
	// CleanedBody string     `json:"cleaned_body" default:""` --DEPRECATED
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
}

type User struct {
	UserID       uuid.UUID    `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	Password     string       `json:"password"`
	EAddress     string       `json:"email"`
	AccessToken  string       `json:"token"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

type CreatedUserResponse struct {
	UserID    uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	IsRedUser bool      `json:"is_chirpy_red"`
}

type LoginUserResponse struct {
	UserID       uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	AccessToken  string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsRedUser    bool      `json:"is_chirpy_red"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Bearer   string `json:"Authorisation"`
}

type RefreshToken struct {
	Token     string       `json:"token"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserID    uuid.UUID    `json:"user_id"`
	ExpiresAt time.Time    `json:"expires_at"`
	RevokedAt sql.NullTime `json:"revoked_at"`
}

type NewAccessTokenResponse struct {
	Token string `json:"token"`
}

type Webhook struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}
