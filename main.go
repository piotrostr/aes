package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"io/ioutil"
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

func GetGCM() *cipher.AEAD {
	key := ReadKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Errorw(
			"error when creating a block",
			"error", err.Error(),
		)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Errorw(
			"error when creating a GCM-wrapped gcm",
			"error", err.Error(),
		)
	}

	return &gcm
}

func GenerateNonce(size int) *[]byte {
	nonce := make([]byte, size)
	n, err := rand.Read(nonce)
	if err != nil {
		logger.Errorw(err.Error())
	}
	logger.Infof("generated nonce: (%d bytes)", n)
	return &nonce
}

func main() {
	payload := *flag.String("payload", "", "string to encrypt")
	file := *flag.String("file", "", "file to encrypt")
	encrypt := *flag.Bool("encrypt", false, "encrypt")
	decrypt := *flag.Bool("decrypt", false, "decrypt")
	flag.Parse()

	if encrypt && decrypt {
		logger.Fatalw("cannot encrypt and decrypt at the same time")
	} else if !encrypt && !decrypt {
		logger.Fatalw("specify either encrypt or decrypt")
	}

	src := []byte(payload)
	if len(src) == 0 {
		logger.Fatalw("no payload provided")
	} else if len(src) < aes.BlockSize {
		logger.Fatalw(
			"payload is too short",
			"payload length:", len(src),
		)
	}

	gcm := *GetGCM()
	if encrypt {
		nonce := *GenerateNonce(gcm.NonceSize())
		encrypted := gcm.Seal(
			nil, // dst
			nonce,
			src,
			nil, // additional data
		)
		logger.Infof("encrypted: %s", encrypted)
	} else if decrypt && file != "" {
		if _, err := os.Stat(file); err != nil {
			logger.Fatalw(
				"file does not exist",
				"file:", file,
			)
		}

		// read in file
		encrypted, err := ioutil.ReadFile(file)
		if err != nil {
			logger.Errorw(err.Error())
		}
		decrypted, err := gcm.Open(
			nil, // dst
			nil, // nonce,
			encrypted,
			nil, // additional data
		)
		if err != nil {
			logger.Errorw(err.Error())
		}
		logger.Infof("decrypted: %s", decrypted)
	}
	/*
		cipher.Decrypt(dst, dst)
		logger.Infof("decrypted: %s", string(dst))
	*/
}
