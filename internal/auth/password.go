package auth

import "golang.org/x/crypto/bcrypt"

// GenerateHash создает bcrypt хэш
func GenerateHash(password string) (string, error) {
	saltedBytes := []byte(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

// Compare сравнивает хэш с паролем
func Compare(hash, password string) bool {
	incoming := []byte(password)
	existing := []byte(hash)

	err := bcrypt.CompareHashAndPassword(existing, incoming)
	return err == nil
}
