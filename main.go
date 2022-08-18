package main

import (
	"crypto/aes"
	"crypto/rand"
	"flag"
	"os"
)

var logger = SetupLogger()

// CreateAndSaveKey creates a new key and saves it to a file.
func CreateAndSaveKey() {
	key := make([]byte, aes.BlockSize)
	_, err := rand.Read(key)
	if err != nil {
		logger.Errorw(err.Error())
	}
	err = os.WriteFile("key", key, 0o644)
	if err != nil {
		logger.Fatalw("failed to write key", "error:", err.Error())
	}
	logger.Infof("key saved: (%d bytes)", len(key))
}

// ReadKey checks if the key exists, otherwise creates it. Afterwards, it reads
// the key from the file.
func ReadKey() []byte {
	if _, err := os.Stat("key"); err != nil {
		CreateAndSaveKey()
	}
	key, err := os.ReadFile("key")
	if err != nil {
		logger.Fatalw(
			"failed to read key",
			"error:", err.Error(),
		)
	}

	if len(key) != aes.BlockSize {
		logger.Fatalw(
			"key is not aes.BlockSize",
			"key length:", len(key),
		)
	}
	return key
}

// ParsePayload parses the payload from the command line arguments.
func ParsePayload() []byte {
	payload := flag.String("payload", "", "string to encrypt")
	flag.Parse()

	src := []byte(*payload)
	if len(src) == 0 {
		logger.Fatalw("no payload provided")
	} else if len(src)%aes.BlockSize != 0 {
		logger.Fatalw(
			"payload is not a multiple of aes.BlockSize",
			"payload length:", len(src),
		)
	}
	return src
}

func main() {
	key := ReadKey()
	cipher, err := aes.NewCipher(key)
	if err != nil {
		logger.Errorw(
			"error when creating a cipher",
			"error message", err.Error(),
		)
	}

	src := ParsePayload()
	dst := make([]byte, len(src))

	{
		cipher.Encrypt(dst, src)
		logger.Infof("encrypted: %s", string(dst))
		cipher.Decrypt(dst, dst)
		logger.Infof("decrypted: %s", string(dst))
	}
}
