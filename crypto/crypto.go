package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io/ioutil"
	"os"

	logging "github.com/piotrostr/aes/logging"
)

var logger = logging.SetupLogger()

// GetFileSize checks if file exists and reads it, otherwise exits.
func GetFileContents(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Fatalw("file does not exist", "path", path)
	}
	if err != nil {
		logger.Errorw(err.Error())
	}
	return file
}

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

// GetGCM a new instance of GCM having read the key from the key file.
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

// GenerateNonce generates random bytes of size n.
func GenerateNonce(size int) *[]byte {
	nonce := make([]byte, size)
	n, err := rand.Read(nonce)
	if err != nil {
		logger.Errorw(err.Error())
	}
	logger.Infof("generated nonce: (%d bytes)", n)
	return &nonce
}

// Encrypt encrypts plaintext with the GCM and returns the ciphertext.
func Encrypt(plaintext []byte) []byte {
	gcm := *GetGCM()
	nonce := *GenerateNonce(gcm.NonceSize())
	ciphertext := gcm.Seal(
		nil, // dst
		nonce,
		plaintext,
		nil, // additional data
	)
	logger.Infof("encrypted: %s", ciphertext)
	return ciphertext
}

// Decrypt decrypts ciphertext with the GCM and returns the plaintext.
func Decrypt(ciphertext []byte) []byte {
	gcm := *GetGCM()
	nonce := *GenerateNonce(gcm.NonceSize())
	plaintext, err := gcm.Open(
		nil,   // dst
		nonce, // nonce,
		ciphertext,
		nil, // additional data
	)
	if err != nil {
		logger.Errorw(err.Error())
	}
	logger.Infof("decrypted: %s", plaintext)
	return plaintext
}
