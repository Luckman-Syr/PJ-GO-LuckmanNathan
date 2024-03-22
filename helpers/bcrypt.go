package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) string {
	salt := 10
	password := []byte(pass)
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, salt)
	return string(hashedPassword)
}

func CheckPassword(hashedPassword, password string) bool {
	byteHash := []byte(hashedPassword)
	bytePass := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePass)

	return err == nil
}
