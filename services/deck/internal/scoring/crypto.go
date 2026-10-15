package scoring

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"unicode/utf8"

	"github.com/azihsoyn/rijndael256"
)

// TODO: Document when the keys were changed
const (
	encryptionKeyV1 = "h89f2-890h2h89b34g-h80g134n90133"
	encryptionKeyV2 = "osu!-scoreburgr---------"
)

func DecryptString(encoded string, iv []byte, key string) (string, error) {
	if len(key) != rijndael256.BlockSize {
		return "", fmt.Errorf("invalid encryption key length %d", len(key))
	}
	if len(iv) != rijndael256.BlockSize {
		return "", fmt.Errorf("invalid IV length %d", len(iv))
	}

	// Score submission receives data as b64 encoded ciphertext
	// encrypted with rijndael256 in CBC mode with PKCS7 padding

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}
	if len(ciphertext) == 0 || len(ciphertext)%rijndael256.BlockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext length %d", len(ciphertext))
	}

	block, err := rijndael256.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("create rijndael cipher: %w", err)
	}
	plaintext := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, ciphertext)

	plaintext, err = removePKCS7Padding(plaintext, rijndael256.BlockSize)
	if err != nil {
		return "", err
	}
	if !utf8.Valid(plaintext) {
		return "", fmt.Errorf("decrypted value is not valid utf8")
	}
	return string(plaintext), nil
}

func removePKCS7Padding(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded data length %d", len(data))
	}
	// Source: https://wgallagher86.medium.com/pkcs-7-padding-in-go-6da5d1d14590

	// Ensure that we got a valid padding length
	paddingLength := int(data[len(data)-1])
	if paddingLength <= 0 || paddingLength > len(data) || paddingLength > blockSize {
		return nil, fmt.Errorf("invalid PKCS7 padding length %d", paddingLength)
	}

	// Validate that the padding bytes are all correct
	for _, value := range data[len(data)-paddingLength:] {
		if int(value) != paddingLength {
			return nil, fmt.Errorf("invalid PKCS7 padding")
		}
	}

	// Return the unpadded data
	return data[:len(data)-paddingLength], nil
}
