package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/wfcornelissen/chirpy/internal/database"
	"github.com/wfcornelissen/chirpy/types"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bitPass, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return "", err
	}
	return string(bitPass), nil
}

func CompareHashAndPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Hour)

	newToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   userID.String(),
		})

	return newToken.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(userID)
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return "", nil
	}

	return token, nil
}

func MakeRefreshToken(cfg *types.ApiConfig, userID uuid.UUID, req *http.Request) (types.RefreshToken, error) {

	random := make([]byte, 32)
	_, err := rand.Read(random)
	if err != nil {
		return types.RefreshToken{}, err
	}
	result := hex.EncodeToString(random)

	timeCreated := time.Now()
	token := types.RefreshToken{
		Token:     result,
		CreatedAt: timeCreated,
		UpdatedAt: timeCreated,
		UserID:    userID,
		ExpiresAt: timeCreated.Add(time.Hour * 24 * 60),
		RevokedAt: sql.NullTime{Valid: false},
	}

	dbTokenParams := database.CreateTokenParams{
		Token:     token.Token,
		CreatedAt: sql.NullTime{Time: token.CreatedAt, Valid: true},
		UpdatedAt: sql.NullTime{Time: token.UpdatedAt, Valid: true},
		UserID:    userID,
		ExpiresAt: sql.NullTime{Time: token.ExpiresAt, Valid: true},
		RevokedAt: sql.NullTime{Valid: false},
	}

	_, err = cfg.Dbquery.CreateToken(req.Context(), dbTokenParams)
	if err != nil {
		return types.RefreshToken{}, err
	}

	return token, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	token := strings.TrimPrefix(authHeader, "ApiKey ")
	if token == authHeader {
		return "", nil
	}

	return token, nil
}
