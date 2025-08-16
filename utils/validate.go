package utils

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func ValidateReq(typeErr, field string) (err error) {
	if typeErr == "" {
		return
	}
	if typeErr != "required" {
		err = MakeError("Invalid Format", field)
		return
	}
	return MakeError("Invalid Mandatory", field)

}

func IsNumeric(s string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	return re.MatchString(s)
}

// HashPassword hashes the plain password using bcrypt.
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("Failed to hash password: " + err.Error())
	}
	return string(hash)
}

func CheckPasswordHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
