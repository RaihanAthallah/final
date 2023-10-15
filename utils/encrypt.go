package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
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
	if err != nil {
		log.Printf("Error creating GCM: %v\n", err)
	}

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
		return nil, fmt.Errorf("error creating cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error creating nonce: %v", err)
	}

	// Read the contents of the input file
	fileStat, err := inputFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file information: %v", err)
	}

	fileSize := fileStat.Size()
	fileData := make([]byte, fileSize)
	_, err = inputFile.Read(fileData)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Encrypt the file contents
	encryptedData := gcm.Seal(nonce, nonce, fileData, nil)

	return encryptedData, nil
}

func DecryptAES(encryptedFilePath string) (string, error) {

	key := os.Getenv("AES_KEY")

	ciphertext, err := hex.DecodeString(encryptedFilePath)
	if err != nil {
		return "", fmt.Errorf("error decoding hex string: %v", err)
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error creating cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("error creating GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext is too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting: %v", err)
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

func EncryptDES(input string) (string, error) {
	key := os.Getenv("DES_KEY")
	if len(key) != 8 {
		return "", fmt.Errorf("dES key must be 8 bytes long")
	}

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error creating DES cipher: %v", err)
	}

	// Ensure the input is a multiple of 8 bytes (the DES block size)
	padSize := 8 - (len(input) % 8)
	if padSize > 0 {
		padding := make([]byte, padSize)
		input += string(padding)
	}

	ciphertext := make([]byte, len(input))
	block.Encrypt(ciphertext, []byte(input))

	return string(ciphertext), nil
}

func DecryptDES(encryptedData string) (string, error) {
	key := os.Getenv("DES_KEY")

	if len(key) != 8 {
		return "", fmt.Errorf("dES key must be 8 bytes long")
	}

	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("error creating DES cipher: %v", err)
	}

	// Ensure the input is a multiple of 8 bytes (the DES block size)
	if len(encryptedData)%8 != 0 {
		return "", fmt.Errorf("invalid encrypted data length")
	}

	plaintext := make([]byte, len(encryptedData))
	block.Decrypt(plaintext, []byte(encryptedData))

	// Trim any trailing null bytes (padding)
	for i := len(plaintext) - 1; i >= 0; i-- {
		if plaintext[i] != 0 {
			plaintext = plaintext[:i+1]
			break
		}
	}

	return string(plaintext), nil
}
