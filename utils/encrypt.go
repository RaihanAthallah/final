package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rc4"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func EncryptAES(filePath string) string {
	key := os.Getenv("AES_KEY")

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Printf("Error creating cipher: %v\n", err)
	}

	gcm, err := cipher.NewGCM(c)

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Printf("Error creating nonce: %v\n", err)
	}

	fmt.Println("Encrypted file path: ", hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(filePath), nil)))
	// return hex string
	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(filePath), nil))
}

func EncryptAES2(inputFile *os.File) ([]byte, error) {
	key := os.Getenv("AES_KEY")

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("Error creating cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("Error creating GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("Error creating nonce: %v", err)
	}

	// Read the contents of the input file
	fileStat, err := inputFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("Error getting file information: %v", err)
	}

	fileSize := fileStat.Size()
	fileData := make([]byte, fileSize)
	_, err = inputFile.Read(fileData)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %v", err)
	}

	// Encrypt the file contents
	encryptedData := gcm.Seal(nonce, nonce, fileData, nil)

	return encryptedData, nil
}

func DecryptAES(encryptedFilePath string) (string, error) {

	key := os.Getenv("AES_KEY")

	ciphertext, err := hex.DecodeString(encryptedFilePath)
	if err != nil {
		return "", fmt.Errorf("Error decoding hex string: %v", err)
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("Error creating cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("Error creating GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("Ciphertext is too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("Error decrypting: %v", err)
	}

	return string(plaintext), nil

}

func EncryptRC4(password string) (string, error) {
	key := os.Getenv("RC4_KEY")
	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	cipherText := make([]byte, len(password))
	cipher.XORKeyStream(cipherText, []byte(password))
	return hex.EncodeToString(cipherText), nil
}

func DecryptRC4(encryptedPassword string) (string, error) {
	key := os.Getenv("RC4_KEY")

	// Decode the hex-encoded ciphertext back to bytes
	ciphertext, err := hex.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}

	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext by XORing it with the RC4 cipher
	plaintext := make([]byte, len(ciphertext))
	cipher.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}
