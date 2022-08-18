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

// GCM struct holds the AEAD instance and the nonce.
type GCM struct {
	aead  cipher.AEAD
	nonce []byte
}

// Initialize initializes the GCM with a nonce and AEAD instance.
func (gcm *GCM) Initialize() {
	gcm.aead = *GetAEAD()
	gcm.nonce = *ReadNonce(gcm.aead.NonceSize())
}

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
		logger.Fatalw("failed to write key", "error", err.Error())
	}
	logger.Infof("key saved: (%d bytes)", len(key))
}

// CreateAndSaveNonce creates a new nonce and saves it to a file.
func CreateAndSaveNonce(size int) {
	nonce := *GenerateNonce(size)
	err := os.WriteFile("nonce", nonce, 0o644)
	if err != nil {
		logger.Fatalw("failed to write nonce", "error", err.Error())
	}
	logger.Infof("nonce saved: (%d bytes)", len(nonce))
}

// ReadKey checks if the key file exists, otherwise creates it. Afterwards, it
// reads the key and nonce from the file.
func ReadKey() []byte {
	path := "key"
	if _, err := os.Stat(path); err != nil {
		CreateAndSaveKey()
	}
	key, err := os.ReadFile(path)
	if err != nil {
		logger.Fatalw(
			"failed to read key",
			"error", err.Error(),
		)
	}

	if len(key) != aes.BlockSize {
		logger.Fatalw(
			"key is not aes.BlockSize",
			"key length:", len(key),
		)
	}

	logger.Infof("read key: (%d bytes)", len(key))
	return key
}

func ReadNonce(size int) *[]byte {
	path := "nonce"
	if _, err := os.Stat(path); err != nil {
		logger.Infow("nonce does not exist, creating new one")
		CreateAndSaveNonce(size)
	}
	nonce, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Fatalw("failed to read nonce", "error", err.Error())
	}

	// validate the nonce size is right
	if len(nonce) != size {
		logger.Fatalw(
			"nonce is not the right length",
			"got", len(nonce),
			"want", size,
		)
	}

	logger.Infof("read nonce: (%d bytes)", len(nonce))
	return &nonce
}

// GetAEAD a new instance of GCM having read the key and nonce from the key file.
func GetAEAD() *cipher.AEAD {
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
func (gcm *GCM) Encrypt(plaintext []byte) []byte {
	ciphertext := gcm.aead.Seal(
		nil, // dst
		gcm.nonce,
		plaintext,
		nil, // additional data
	)
	logger.Infof("encrypted: %s", ciphertext)
	return ciphertext
}

// Decrypt decrypts ciphertext with the GCM and returns the plaintext.
func (gcm *GCM) Decrypt(ciphertext []byte) []byte {
	plaintext, err := gcm.aead.Open(
		nil,       // dst
		gcm.nonce, // nonce,
		ciphertext,
		nil, // additional data
	)
	if err != nil {
		logger.Errorw(err.Error())
	}
	logger.Infof("decrypted: %s", plaintext)
	return plaintext
}
