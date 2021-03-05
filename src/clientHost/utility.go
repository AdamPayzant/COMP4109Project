package main

import (
	"crypto/sha256"
	"encoding/base64"
	"regexp"
	"strings"
)

//Utility functions

//Used to prevent Cross Site Scripting (Make other functions)
//Messages should pass through this.
//So far tries to break tags
func cleanText(text string) string {

	strings.ReplaceAll(text, "<", "&lt;")
	strings.ReplaceAll(text, ">", "&gt;")

	return text

}

// Function to Convert a string into a sha256 hash
// Currently set for sha256 in base64 encoding
// Used for password function
//
func standardHashFunction(text string) string {

	return string(base64.StdEncoding.EncodeToString(sha256.New().Sum([]byte(text))))

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

func passwordChecker(text string) bool {

	if standardHashFunction(text) == standardHashFunction(debugPassword) {
		return true
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

	return "error", "null", "null"

}

func clientJSONCreator(functionType string, payload string) string {

	return "{\"type\":\"" + functionType + "\",\"payload\":\"" + payload + "\"}"

}
