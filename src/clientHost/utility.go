package main

import "strings"

//Utility functions

//Used to prevent Cross Site Scripting (Make other functions)
//Messages should pass through this.
//So far tries to break tags
func cleanText(text string) string {

	strings.ReplaceAll(text, "<", "&lt;")
	strings.ReplaceAll(text, ">", "&gt;")

	return text

}
