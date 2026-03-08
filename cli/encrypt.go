package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/theriverman/taigo/cli/passwordbasedencryption"
)

const encryptionVersionPrefix = "v2:"

func encryptionKey(machineID string) []byte {
	sum := sha256.Sum256([]byte("taigo-cli:" + machineID))
	return sum[:]
}

func PasswordEncrypt(plainPassword string) string {
	machineID, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(encryptionKey(machineID))
	if err != nil {
		log.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plainPassword), nil)
	payload := append(nonce, ciphertext...)
	return encryptionVersionPrefix + base64.StdEncoding.EncodeToString(payload)
}

func PasswordDecrypt(encryptedPassword string) string {
	machineID, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(encryptedPassword, encryptionVersionPrefix) {
		rawPayload := strings.TrimPrefix(encryptedPassword, encryptionVersionPrefix)
		payload, err := base64.StdEncoding.DecodeString(rawPayload)
		if err != nil {
			log.Fatal(err)
		}

		block, err := aes.NewCipher(encryptionKey(machineID))
		if err != nil {
			log.Fatal(err)
		}
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			log.Fatal(err)
		}
		if len(payload) < gcm.NonceSize() {
			log.Fatal("encrypted payload is too short")
		}
		nonce := payload[:gcm.NonceSize()]
		ciphertext := payload[gcm.NonceSize():]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			log.Fatal(err)
		}
		return string(plaintext)
	}

	// Backwards compatibility for pre-v2 CLI config files.
	s, err := passwordbasedencryption.Decrypt(machineID, 5, encryptedPassword)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
