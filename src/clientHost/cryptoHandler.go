package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"

	_ "github.com/mattn/go-sqlite3"
)

var clientPublicKey *rsa.PublicKey

func RSA_OAEP_Encrypt(secretMessage string, key *rsa.PublicKey) (string, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, key, []byte(secretMessage), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RSA_OAEP_Decrypt(cipherText string, privKey *rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, privKey, ct, label)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func encryptForClient(msg string) (string, error) {
	cypherText, err := RSA_OAEP_Encrypt(msg, clientPublicKey)
	if err != nil {
		return "", err
	}

	return cypherText, nil
}

func bytesToKey(bytes []byte) (*rsa.PublicKey, error) {
	key, e := x509.ParsePKIXPublicKey(bytes)
	return key.(*rsa.PublicKey), e
}

func keyToBytes(key *rsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(key)
}

func encryptForSending(msg string, userInfo *UserInfo) (string, error) {
	cypherText, err := RSA_OAEP_Encrypt(msg, userInfo.key)
	if err != nil {
		return "", err
	}

	return cypherText, nil
}

func generateSecret(userInfo *UserInfo) string {
	return "asdsad"
}

func verifySecret(userInfo *UserInfo, secret string) bool {
	return true
}

func authenticateClientToken(token []byte) bool {
	return true
}

func decryptToken(token []byte) ([]byte, error) {
	// todo
	return token, nil
}
