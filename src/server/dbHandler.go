package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

const (
	timeOut = 900000000000 // 15 minutes in nanoseconds
	dbName  = "smvs:password@tcp(localhost:3306)/smvsserver"
)

type User struct {
	UserID   int    `json:"userID"`
	Username string `json:"username"`
	PubKey   []byte `json:"pubKey"`
	IP       string `json:"ip"`
}
type Token struct {
	TokenID  int       `json:"tokenID"`
	UserID   int       `json:"userID"`
	Content  []byte    `json:"content"`
	Assigned time.Time `json:"assigned"`
}

/*
	Connects to the DB

	Return:
		error - If there's an error
*/
func connect() error {
	db, _ = sql.Open("mysql", dbName+"?parseTime=true")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	c := `
		CREATE TABLE IF NOT EXISTS users (
			userID INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
			username TEXT,
			pubKey BLOB,
			ip TEXT
		);
	`
	_, err = db.Exec(c)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c = `
		CREATE TABLE IF NOT EXISTS tokens (
			tokenID INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
			userID INTEGER,
			content BLOB,
			assigned TIMESTAMP,
			FOREIGN KEY (userID) REFERENCES users(userID)
		);
	`
	_, err = db.Exec(c)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c = `
		CREATE TRIGGER IF NOT EXISTS userDel
			AFTER DELETE ON users FOR EACH ROW
		BEGIN
			DELETE FROM tokens
			WHERE userID=OLD.userID;
		END;
	`
	_, err = db.Exec(string(c))
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Connected to db")
	return nil
}

/*
	Adds a user to the system

	Input:
		username - A unique username
		publickey - The user's public key
		ip - The user's host's IP address
	Return:
		error - If there's an error adding the user to the system
*/
func addUser(username string, publickey *rsa.PublicKey, ip string) error {
	keyBytes, _ := x509.MarshalPKIXPublicKey(publickey)
	_, err := db.Exec("INSERT INTO users (username, pubKey, ip) VALUES(?, ?, ?)",
		username, keyBytes, ip)
	if err != nil {
		return err
	}
	return nil
}

/*
	Adds a token to a user's profile
	The token will be sent encrypted by the user's public key, and must be decrypted to use
	The token will expire after 30 minutes of unuse

	Input:
		username - The profile's username
	Output:
		string - The user's encrypted token
		error  - nil unless error occurs
*/
func addToken(username string) ([]byte, error) {
	var userID int
	var key []byte
	// Gets the user's id and public key
	err := db.QueryRow("SELECT userID, pubKey FROM users WHERE username=?", username).Scan(&userID, &key)
	if err != nil {
		return nil, err
	}

	PKIXpubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	pubKey := PKIXpubKey.(*rsa.PublicKey)

	// Generates a 64 character long token
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 64)
	for i := 0; i < len(b); i++ {
		b[i] = letters[rand.Intn(len(letters))]
	}
	str := string(b)

	_, err = db.Exec("INSERT INTO tokens (userid, content) VALUES(?, ?)", userID, str)
	if err != nil {
		return nil, err
	}
	// Encrypts the token
	encryptedToken, er := rsa.EncryptOAEP(
		sha512.New(),
		crand.Reader,
		pubKey,
		[]byte(str),
		nil)
	if er != nil {
		return nil, er
	}

	go pruneTokens()
	return encryptedToken, nil
}

/*
	Checks if a user's sent token is valid and unexpired
	Updates the timer on the token

	Input:
		username - the profile's username
		token - the decrypted token
	Output:
		bool - Whether the token is accepted or not
*/
func checkToken(username string, token []byte) (bool, error) {
	pruneTokens()
	var userID int
	err := db.QueryRow("Select userID FROM users WHERE username=?", username).Scan(&userID)
	if err != nil {
		return false, err
	}
	err = db.QueryRow("Select tokenID FROM tokens WHERE userID=? AND content=?",
		userID, token).Scan(&token)
	if err != nil {
		if err != sql.ErrNoRows {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

/*
A function that continues pruning after parent function finishes
*/
func pruneTokens() {
	_, err := db.Exec("DELETE FROM tokens WHERE assigned<?", time.Now().Add(-time.Minute*15).Format(time.RFC3339))
	if err != nil {
		log.Fatal(err)
	}
}

/*
	Updates a user's IP address

	Input:
		username - the user's unique username
		ip - the user's new IP
	Output:
		error - If there's an error updating the IP for whatever reason
*/
func updateIP(username string, ip string) error {
	_, err := db.Exec("UPDATE users SET ip=? WHERE username=?", ip, username)
	return err
}

/*
	Updates a user's public key

	Input:
		username - The user's unique username
		publickey - The user's new public key
	Output:
		error - Reports if an error occured
*/
func updateKey(username string, publicKey *rsa.PublicKey) error {
	keyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	_, err := db.Exec("UPDATE users SET pubKey=? WHERE username=?", keyBytes, username)
	return err
}

// TODO: Write this
func searchUser(partialname string) ([]string, error) {
	return nil, nil
}

/*
	Gets a user's IP and public key

	Input:
		username - The target user's unique username
	Output:
		ip - The target user's Host IP address and port
		publickey - The target user's public key
		err - If there's an error finding the user
*/
func getUser(username string) (ip string, key []byte, err error) {
	err = db.QueryRow("SELECT pubKey, ip FROM users WHERE username=?", username).Scan(&key, &ip)
	return ip, key, err
}
