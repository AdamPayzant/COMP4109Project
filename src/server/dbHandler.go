package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	crand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection = nil

type User struct {
	Username string
	Key      rsa.PublicKey
	IP       string
	Tokens   []string
}

/*
	Connects to the DB

	Return:
		error - If there's an error
*/
func connect() error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connects to the DB server
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	// Checks the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	collection = client.Database("SMVS").Collection("Users")

	fmt.Println("Connected to DB")
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
func Register(username string, publickey *rsa.PublicKey, ip string) error {
	// Does some error checking
	if collection == nil {
		return errors.New("Collection Not Defined")
	}
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.TODO(), filter)
	if err == nil {
		return errors.New("User Already Exists")
	}

	user := User{username, *publickey, ip, nil}
	insertResults, er := collection.InsertOne(context.TODO(), user)
	if er != nil {
		return er
	}

	fmt.Printf("Successfully added user %q", insertResults.InsertedID)
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
	if collection == nil {
		return nil, errors.New("Collection Not Defined")
	}

	// Gets the user
	var user User
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, errors.New("User Does Not Exist")
	}

	// Generates a 64 character long token
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 64)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	str := string(b)

	var keys = make([]string, len(user.Tokens))
	for i := range b {
		keys[i] = user.Tokens[i]
	}
	keys[len(keys)-1] = str
	update := bson.M{"$set": bson.M{"Tokens": keys}}

	// Encrypts the token
	encryptedToken, er := rsa.EncryptOAEP(
		sha256.New(),
		crand.Reader,
		&user.Key,
		[]byte(str),
		nil)
	if er != nil {
		return nil, er
	}

	// If there's been no errors up until this point, adds token to user
	collection.UpdateOne(context.TODO(), filter, update)

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
func checkToken(username string, token string) (bool, error) {
	if collection == nil {
		return false, errors.New("Collection Not Defined")
	}

	// Gets the user
	var user User
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return false, errors.New("User Does Not Exist")
	}

	for _, b := range user.Tokens {
		if token == b {
			return true, nil
		}
	}

	return false, nil
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
	if collection == nil {
		return errors.New("Collection Not Defined")
	}

	// Gets the user
	var user User
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return errors.New("User Does Not Exist")
	}

	update := bson.M{"$set": bson.M{"IP": ip}}
	collection.UpdateOne(context.TODO(), filter, update)

	return nil
}

/*
	Updates a user's public key

	Input:
		username - The user's unique username
		publickey - The user's new public key
	Output:
		error - Reports if an error occured
*/
func updateKey(username string, publicKey rsa.PublicKey) error {
	if collection == nil {
		return errors.New("Collection Not Defined")
	}

	// Gets the user
	var user User
	filter := bson.D{{Key: "username", Value: username}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return errors.New("User Does Not Exist")
	}

	update := bson.M{"$set": bson.M{"Key": publicKey}}
	collection.UpdateOne(context.TODO(), filter, update)

	return nil
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
func getUser(username string) (ip string, key rsa.PublicKey, err error) {
	if collection == nil {
		return "", key, errors.New("Collection Not Defined")
	}

	// Gets the user
	var user User
	filter := bson.D{{Key: "username", Value: username}}
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return ip, key, errors.New("User Does Not Exist")
	}

	ip = user.IP
	key = user.Key

	return ip, key, nil
}
