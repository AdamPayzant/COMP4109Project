package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"log"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"

	_ "github.com/mattn/go-sqlite3"
)

var clientPublicKey *rsa.PublicKey

func RSA_OAEP_Encrypt(plaintext string, key *rsa.PublicKey) (string, error) {
	rng := rand.Reader
	encryptedMSG, er := rsa.EncryptOAEP(sha512.New(), rng, key, []byte(plaintext), nil)
	if er != nil {
		return "", er
	}
	return base64.StdEncoding.EncodeToString(encryptedMSG), nil
}

func RSA_OAEP_Decrypt(cipherText string, privKey *rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha512.New(), rng, privKey, ct, nil)
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

func verifySecret(userInfo *UserInfo, secret []byte) bool {
	hash := sha512.New()
	hash.Write([]byte(userInfo.username))
	err := rsa.VerifyPKCS1v15(userInfo.key, crypto.SHA512, hash.Sum(nil), secret)
	if err != nil {
		log.Printf("Error verifying Secret: %s\n", err)
		return false
	}

	return true
}

func authenticateClientToken(token []byte) bool {
	hash := sha512.New()
	hash.Write([]byte(settings.Token))
	err := rsa.VerifyPKCS1v15(clientPublicKey, crypto.SHA512, hash.Sum(nil), token)
	if err != nil {
		log.Printf("Error verifying Secret Token: %s\n", err)
		return false
	}

	return true
}

func decryptToken(token []byte) ([]byte, error) {
	// todo
	return token, nil
}
