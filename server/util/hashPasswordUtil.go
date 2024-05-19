package util

import (
	"social/project/initializer"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	logger := initializer.InitializeLogger()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		logger.Fatal(err)
	}
	return string(hashedPassword), nil
}
