package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var key = []byte("erIje49Kl;weoKoerkKEmwmdKEJr849P")

// HashPassword will encrypt the password and secure it
func HashPassword(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		panic(1)
	}
	return hash
}

// CheckPassword will simply check the password to ensure it's correct
func CheckPassword(password string, hashedPassword []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return err
	}
	return nil
}

// Encrypt will simply do as you'd think
func Encrypt(messages string) ([]byte, error) {
	cphr, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return []byte(messages), err
	}
	gcm, err := cipher.NewGCM(cphr)
	if err != nil {
		log.Fatal(err)
		return []byte(messages), err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
		return []byte(messages), err
	}
	return gcm.Seal(nonce, nonce, []byte(messages), nil), nil
}

// Decrypt will simply do as you'd think
func Decrypt(hash []byte) (string, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return string(hash), err
	}
	gcmDecrypt, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
		return string(hash), err
	}
	nonceSize := gcmDecrypt.NonceSize()
	if len(hash) < nonceSize {
		log.Fatal(err)
		return string(hash), err
	}
	nonce, encryptedMessage := hash[:nonceSize], hash[nonceSize:]
	plaintext, err := gcmDecrypt.Open(nil, nonce, encryptedMessage, nil)
	if err != nil {
		log.Fatal(err)
		return string(hash), err
	}
	return string(plaintext), nil
}
