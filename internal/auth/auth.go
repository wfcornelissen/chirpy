package auth

import (
	"github.com/pingcap/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bitPass, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		log.Errorf("Error hashing password.")
	}
	return string(bitPass), nil
}

func CompareHashAndPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	if err != nil {
		return err
	}
	return nil
}
