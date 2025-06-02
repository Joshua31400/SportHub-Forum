package authentification

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", fmt.Errorf("le mot de passe doit contenir au moins 6 caractères")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsPasswordStrong(password string) bool {
	return len(password) >= 8 &&
		containsUpper(password) &&
		containsDigit(password)
}

func containsUpper(password string) bool {
	return strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func containsDigit(password string) bool {
	return strings.ContainsAny(password, "0123456789")
}
