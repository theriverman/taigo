package main

import (
	"log"

	"github.com/denisbrodbeck/machineid"
	"github.com/theriverman/taigo/cli/passwordbasedencryption"
)

func PasswordEncrypt(plainPassword string) string {
	machineID, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}
	s, err := passwordbasedencryption.Encrypt(machineID, 5, plainPassword)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func PasswordDecrypt(encryptedPassword string) string {
	machineID, err := machineid.ID()
	if err != nil {
		log.Fatal(err)
	}
	s, err := passwordbasedencryption.Decrypt(machineID, 5, encryptedPassword)
	if err != nil {
		log.Fatal(err)
	}
	return s
}
