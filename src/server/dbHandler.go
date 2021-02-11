package main

import (
	"errors"
	"fmt"
)

/*
Adds a user to the system

Input:
	username - A unique username
	publickey - The user's public key
	ip - The user's host's IP address
Return:
	error - If there's an error adding the user to the system
*/
func addUser(username string, publickey string, ip string) error {
	if username == "" {
		return errors.New("Empty Username")
	}
	fmt.Printf("Successfully added user %q", username)
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
*/
func addToken(username string) string {
	return ""
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
func checkToken(username string, token string) bool {
	return true
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
func updateKey(username string, publicKey string) error {
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
func getUser(username string) (ip string, publickey string, err error) {
	ip = "000.000.000.000:111"
	publickey = ""
	err = nil

	return ip, publickey, err
}
