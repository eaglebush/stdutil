package stdutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"errors"
	"io"
)

// commandSeedKey - the universal seed key for command
var commandSeedKey = []byte{
	0x9f, 0x9b, 0xf0, 0x4b, 0xd7, 0xa4, 0xae, 0x27, 0xe6, 0x34, 0xf6, 0x7d, 0x07, 0x9b, 0xcc, 0x92,
	0xdd, 0x9b, 0x78, 0xcb, 0x9a, 0x83, 0x47, 0xfd, 0x74, 0xaa, 0xa8, 0x1d, 0xf4, 0x90, 0xb9, 0xae,
}

// Encrypt - encrypt a string using AES
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt - decrypt a string using AES
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// EncodeText encodes plain text with a key and returns an encrypted and base64-encoded string
func EncodeText(plainText string, key []byte) string {
	if plainText == "" {
		return plainText
	}
	enc, _ := Encrypt([]byte(plainText), key)
	return b64.RawURLEncoding.WithPadding(b64.NoPadding).EncodeToString(enc)
}

// DecodeText decodes an encypted base64-encoded text with a key and returns a decrypted string
func DecodeText(encoded string, key []byte) string {
	if encoded == "" {
		return encoded
	}
	benc, _ := b64.RawURLEncoding.WithPadding(b64.NoPadding).DecodeString(encoded)
	dec, _ := Decrypt(benc, key)
	return string(dec)
}

// EncodeCommand encodes a command to be decoded by the receiving page
//
// # It uses the library key to encrypt the string and later encoded with base64
//
// Deprecated: Use the EncodeText and DecodeText functions with a 32-bit application-defined key
func EncodeCommand(command string) string {
	if command == "" {
		return command
	}
	enc, _ := Encrypt([]byte(command), commandSeedKey)
	return b64.RawURLEncoding.WithPadding(b64.NoPadding).EncodeToString(enc)
}

// DecodeCommand decodes an encrypted and base64-encoded text for the receiving page
//
// # It uses the library key to decrypt the string after the input has been decoded from base64
//
// Deprecated: Use the EncodeText and DecodeText functions with a 32-bit application-defined key
func DecodeCommand(encoded string) string {
	if encoded == "" {
		return encoded
	}
	benc, _ := b64.RawURLEncoding.WithPadding(b64.NoPadding).DecodeString(encoded)
	dec, _ := Decrypt(benc, commandSeedKey)
	return string(dec)
}
