/*
	Original author:	LucasSloan/passwordbasedencryption
	Original GitGub:	https://github.com/LucasSloan/passwordbasedencryption
*/

package passwordbasedencryption

import (
	"strings"

	"encoding/base64"

	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
)

func getDerivedKey(password string, salt string, count int) ([]byte, []byte) {
	key := md5.Sum([]byte(password + salt))
	for i := 0; i < count-1; i++ {
		key = md5.Sum(key[:])
	}
	return key[:8], key[8:]
}

func Encrypt(password string, obtenationIterations int, plainText string) (string, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	encText, err := doEncrypt(password, plainText, salt, obtenationIterations)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(append(salt, encText...)), nil
}

func Decrypt(password string, obtenationIterations int, cipherText string) (string, error) {
	msgBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	salt := msgBytes[:8]
	encText := msgBytes[8:]
	return doDecrypt(password, encText, salt, obtenationIterations)
}

func EncryptWithFixedSalt(password string, obtenationIterations int, plainText string, fixedSalt string) (string, error) {
	salt := make([]byte, 8)
	copy(salt[:], fixedSalt)

	encText, err := doEncrypt(password, plainText, salt, obtenationIterations)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encText), nil
}

func DecryptWithFixedSalt(password string, obtenationIterations int, cipherText string, fixedSalt string) (string, error) {
	msgBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	salt := make([]byte, 8)
	copy(salt[:], fixedSalt)
	encText := msgBytes[:]
	return doDecrypt(password, encText, salt, obtenationIterations)
}

func doEncrypt(password, plainText string, salt []byte, obtenationIterations int) ([]byte, error) {
	padNum := byte(8 - len(plainText)%8)
	for i := byte(0); i < padNum; i++ {
		plainText += string(padNum)
	}

	dk, iv := getDerivedKey(password, string(salt), obtenationIterations)
	block, err := des.NewCipher(dk)
	if err != nil {
		return nil, err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plainText))
	encrypter.CryptBlocks(encrypted, []byte(plainText))

	return encrypted, nil
}

func doDecrypt(password string, encText, salt []byte, obtenationIterations int) (string, error) {
	dk, iv := getDerivedKey(password, string(salt), obtenationIterations)
	block, err := des.NewCipher(dk)

	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encText))
	decrypter.CryptBlocks(decrypted, encText)

	decryptedString := strings.TrimRight(string(decrypted), "\x01\x02\x03\x04\x05\x06\x07\x08")

	return decryptedString, nil
}
