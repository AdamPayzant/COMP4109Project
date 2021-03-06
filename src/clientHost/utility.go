package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	//"code.google.com/p/go.crypto/scrypt"
)

//Utility functions

//Used to prevent Cross Site Scripting (Make other functions)
//Messages should pass through this.
//So far tries to break tags
func cleanText(text string) string {

	var tempText string
	tempText = strings.ReplaceAll(text, "<", "&lt;")
	tempText = strings.ReplaceAll(tempText, ">", "&gt;")

	//return tempText
	return text
}

// Function to Convert a string into a sha256 hash
// Currently set for sha256 in base64 encoding
// Used for password function
//
func standardHashFunction(text string) string {

	return base64.StdEncoding.EncodeToString(sha256.New().Sum([]byte(text)))

}

// Function to check if a string is a valid Hash (If change if the hash Function changes)
// Currently set for sha256 in base64 encoding
//
func validHashCheck(text string) bool {

	result, err := regexp.MatchString("[a-fA-F0-9]{64}", text)

	if err != nil {
		return false
	}

	return result

}

var debugPassword string = "12345678"

func createPassword(password string) string {
	return standardHashFunction(password)
}

func createPasswordSalt(password string) (string, string) {

	noise := make([]byte, 2)
	_, err := io.ReadFull(rand.Reader, noise)
	if err != nil {
		log.Fatal(err)
	}

	pepper := string(noise)
	//pepper := string(hex.EncodeToString(noise))
	//pepper := base64.StdEncoding.EncodeToString(noise)

	var hash = standardHashFunction(pepper + password)
	return hash, pepper

}

func passwordChecker(text string) bool {

	//Old Function should not be needed
	if standardHashFunction(text) == standardHashFunction(debugPassword) {
		return true
	}

	return false

}

//Keep Here
var testPassword, passwordSalt = createPasswordSalt(debugPassword)

func passwordCheckerWithSalt(password string, salt string) bool {

	if standardHashFunction(salt+password) == testPassword {
		fmt.Print("It works")
		return true
	}
	return false
}

func passwordCheckerWithPepper(text string) bool {

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {

			//To generate All Salts
			b := make([]byte, 2)
			b[0] = byte(i)
			b[1] = byte(j)
			tempPepper := string(b)

			//Check if Valid
			if passwordCheckerWithSalt(text, tempPepper) {
				fmt.Print("It works")
				return true
			}
		}
	}

	return false

}

//Function to turn a json object sent from the browser into 3 strings of text
//
//Output format:
//    functionClass: Top level classification of the request made. Used to select router function
//    functionName: Used to specify an operation preformed
//    payload: The arguments for the function (further parsing maybe required)
//
func clienetJSONParser(inputText string) (string, string, string) {

	type asdklsadllsdakl struct {
		A string `json:"functionClass"`
		B string `json:"functionName"`
		C string `json:"payload"`
	}

	var sdassd asdklsadllsdakl

	if !json.Valid([]byte(inputText)) {
		return "Error", "Error", "Warning"
	}

	err := json.Unmarshal([]byte(inputText), &sdassd)

	if err != nil {
		print(err)
	}

	if &sdassd == nil {

		return "A", "A", "s"
	}

	fmt.Println(sdassd.A)

	return sdassd.A, sdassd.B, sdassd.C

}

func clientJSONCreator(functionType string, payload string) string {

	return "{\"type\":\"" + functionType + "\",\"payload\":\"" + strings.ReplaceAll(payload, "\"", "\\\"") + "\"}"

}
