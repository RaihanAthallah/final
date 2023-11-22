package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
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
	padSize := des.BlockSize - (len(input) % des.BlockSize)
	if padSize > 0 {
		padding := make([]byte, padSize)
		input += string(padding)
	}

	mode := cipher.NewCBCEncrypter(block, make([]byte, des.BlockSize))

	ciphertext := make([]byte, len(input))
	mode.CryptBlocks(ciphertext, []byte(input))

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

	mode := cipher.NewCBCEncrypter(block, make([]byte, des.BlockSize))

	plaintext := make([]byte, len(encryptedData))
	mode.CryptBlocks(plaintext, []byte(encryptedData))

	// Trim any trailing null bytes (padding)
	for i := len(plaintext) - 1; i >= 0; i-- {
		if plaintext[i] != 0 {
			plaintext = plaintext[:i+1]
			break
		}
	}

	return string(plaintext), nil
}

func GenerateKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	publicKey := privateKey.PublicKey

	// Convert to easily stored format
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyDER}
	publicKeyPEM := &pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyDER}

	publicKeyString := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(publicKeyPEM))
	privateKeyString := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(privateKeyPEM))

	return privateKeyString, publicKeyString, nil
}

func EncryptRSA(data string, publicKeyString string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyString))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing public key")
	}

	// Parse the public key.
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse public key:", err)
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data))

	return string(ciphertext), err
}

func DecryptRSA(ciphertext string, privateKeyString string) (string, error) {
	block, _ := pem.Decode([]byte(privateKeyString))
	if block == nil {
		return "", errors.New("failed to parse PEM block containing private key")
	}

	// Parse the private key.
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse private key:", err)
		return "", err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, []byte(ciphertext))

	return string(plaintext), err
}
