package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

func generateHashPassword(password string) (string, error) {

	//::::CREATING salt
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	//:::GENERATING HASH
	hash, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err

	}
	//:::RETURN PASSWORD
	saltedHash := fmt.Sprintf("%s.%s", hex.EncodeToString(hash), hex.EncodeToString(salt))
	return saltedHash, nil
}

func Compare(storedPassword string, suppliedPassword string) (bool, error) {
	passwordSaltedBytes := strings.Split(storedPassword, ".")
	saltBytes, err := hex.DecodeString(passwordSaltedBytes[1])
	if err != nil {
		return false, err
	}

	hash, err := scrypt.Key([]byte(suppliedPassword), saltBytes, 32768, 8, 1, 32)
	if err != nil {
		return false, err

	}
	return hex.EncodeToString(hash) == passwordSaltedBytes[0], nil

}
