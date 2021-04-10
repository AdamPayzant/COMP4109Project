package main

import (
	"errors"
	"log"

	"crypto/rsa"
	"crypto/x509"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type UserInfo struct {
	username string
	msgCount int
	ip       string
	key      *rsa.PublicKey
}

type Message struct {
	Order       int    `json:"order"`
	Speaker     bool   `json:"speaker"`
	MessageText string `json:"messageText"`
}

const (
	timeOut = 900000000000 // 15 minutes in nanoseconds
)

var db *sql.DB = nil

func updateOrAddUserInfo(userInfo *UserInfo) error {
	var e error
	e = addUserInfo(userInfo)
	if e != nil {
		r, er := updateUserInfo(userInfo)
		if er != nil {
			log.Printf("Could not save user info: %s : %v", userInfo.username, er)
			e = er
		}

		count, err := r.RowsAffected()

		if err != nil {
			log.Fatalf("Connot get local user info db update results!: %v", err)
		}

		if count > 0 {
			log.Printf("Could not save user info!")
		}
	}

	return e
}

/*
	Connects to the DB

	Return:
		error - If there's an error
*/
func connect(dbName string) error {
	db, _ = sql.Open("mysql", dbName+"?parseTime=true")
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	c := `
		CREATE TABLE IF NOT EXISTS userInfo (
			userID INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
			username TINYTEXT NOT NULL UNIQUE,
			msgCount INTEGER NOT NULL,
			pubKey BLOB NOT NULL,
			ip TINYTEXT NOT NULL
		);
	`
	_, err = db.Exec(c)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c = `
		CREATE TABLE IF NOT EXISTS conversation (
			userID INTEGER NOT NULL,
			username TINYTEXT NOT NULL,
			id INTEGER NOT NULL,
			speaker BOOL NOT NULL,
			assigned TIMESTAMP,
			msg TEXT NOT NULL,
			FOREIGN KEY (userID) REFERENCES userInfo(userID),
			PRIMARY KEY (userID, id)
		);
	`
	_, err = db.Exec(c)
	if err != nil {
		log.Fatal(err)
		return err
	}

	c = `
		CREATE TRIGGER IF NOT EXISTS countMSG
			BEFORE INSERT ON conversation FOR EACH ROW
		BEGIN
			UPDATE userInfo SET msgCount = msgCount + 1
			WHERE userInfo.userID = NEW.userID;
			UPDATE userInfo
			SET NEW.id=userInfo.msgCount, NEW.username = userInfo.username
			WHERE userInfo.userID = NEW.userID;
		END;
	`
	_, err = db.Exec(string(c))
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Connected to db: " + dbName)
	return nil
}

/*
	Adds a user to the system

	Input:
		UserInfo - The stuct that contains the information about the user
	Return:
		error - If there's an error adding the user to the system
*/
func addUserInfo(userInfo *UserInfo) error {
	keyBytes, e := x509.MarshalPKIXPublicKey(userInfo.key)
	if e != nil {
		return e
	}
	_, err := db.Exec("INSERT INTO userInfo (username, msgCount, pubKey, ip) VALUES(?, ?, ?, ?)", userInfo.username, 0, keyBytes, userInfo.ip)
	if err != nil {
		return err
	}
	return nil
}

/*
	Updates a user's public key and ip

	Input:
		UserInfo - The stuct that contains the information about the user
	Output:
		error - Reports if an error occured
*/
func updateUserInfo(userInfo *UserInfo) (sql.Result, error) {
	var err error
	var res sql.Result = nil
	if userInfo.key != nil {
		var keyBytes []byte
		keyBytes, err = x509.MarshalPKIXPublicKey(userInfo.key)
		if err != nil {
			return res, err
		}

		if userInfo.ip != "" {
			res, err = db.Exec("UPDATE userInfo SET pubKey=?, ip=? WHERE username=?", keyBytes, userInfo.ip, userInfo.username)
		} else {
			res, err = db.Exec("UPDATE userInfo SET pubKey=?, WHERE username=?", keyBytes, userInfo.username)
		}
	} else {
		if userInfo.ip != "" {
			res, err = db.Exec("UPDATE userInfo SET ip=? WHERE username=?", userInfo.ip, userInfo.username)
		}
	}

	return res, err
}

/*
	Gets a user's IP and public key

	Input:
		username - The target user's unique username
	Output:
		UserInfo - The stuct that contains the information about the target user
		err - If there's an error finding the user
*/
func getUserInfo(username string) (*UserInfo, error) {
	var err error = errors.New("No such user in DB")
	var msgCount int
	var keyBytes []byte
	var ip string
	err = db.QueryRow("SELECT msgCount, pubKey, ip FROM userInfo WHERE username=?", username).Scan(&msgCount, &keyBytes, &ip)
	if err != nil {
		return nil, err
	}
	key, e := x509.ParsePKIXPublicKey(keyBytes)
	if e != nil {
		log.Printf("could not get userInfo: %v", e)
		return nil, e
	}
	return &UserInfo{username: username, msgCount: msgCount, ip: ip, key: key.(*rsa.PublicKey)}, err
}

/*
	Adds a messages to the system

	Input:
		username - A unique username
		msg - The Message struct that contains all the imformation about a message
	Return:
		error - If there's an error adding the user to the system
*/
func addMessage(username string, msg *Message) (sql.Result, error) {
	var err error
	var res sql.Result = nil
	var userID int
	err = db.QueryRow("SELECT userID FROM userInfo WHERE username=?", username).Scan(&userID)
	if err != nil {
		return res, err
	}
	res, err = db.Exec("INSERT INTO conversation (userID, username, speaker, msg) VALUES(?, ?, ?, ?)", userID, username, msg.Speaker, msg.MessageText)
	if err != nil {
		return res, err
	}
	return res, nil
}

/*
	Gets all messages between a user

	Input:
		username - The target user's unique username
	Output:
		conversation - The list of messages from and to the target user
		err - If there's an error finding the user
*/
func getConversationFromDB(username string) ([]Message, error) {
	var conversation []Message = nil
	rowCount, err := db.Query("SELECT COUNT(*) AS count FROM conversation WHERE username=?", username)
	if err != nil {
		return nil, err
	}

	if rowCount.Next() {
		var count int
		rowCount.Scan(&count)
		rowCount.Close()

		conversation = make([]Message, count)

		msgIndex := 0
		rows, err := db.Query("SELECT id, speaker, msg FROM conversation WHERE username=?", username)
		for rows.Next() {
			var id int
			var speaker bool
			var msg string
			rows.Scan(&id, &speaker, &msg)
			if err != nil {
				log.Println(err)
			}
			conversation[msgIndex].Order = id
			conversation[msgIndex].Speaker = speaker
			conversation[msgIndex].MessageText = msg
			msgIndex = msgIndex + 1
		}
		rows.Close()
	}

	return conversation, err
}

/*
	Delete a messages from a user

	Input:
		username - The user's unique username
		id - The message id
	Output:
		error - Reports if an error occured
*/
func deleteMessageFromDB(username string, id int64) error {
	_, e := db.Exec("DELETE FROM conversation WHERE username=? AND id=?", username, id)
	if e != nil {
		return e
	}
	return nil
}
