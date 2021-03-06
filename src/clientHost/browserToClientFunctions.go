package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//var connection *websocket.Conn = nil

func websocketRequestHander(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	//Reject if another client is connected
	//if connection != nil {
	//	conn.Close()
	//	return
	//}

	if authenticationFunction(conn) == -1 {
		return
	}

	for {

		if -1 == clientMessageRouterFunction(conn) {
			break
		}

	}

}

func authenticationFunction(conn *websocket.Conn) int {

	// Print the message to the console
	//fmt.Printf("Password sent: %s\n", string(password))

	var MAX_ATTEMPTS = 5

	//Loop for login
	var attemps = 0
	for {

		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return -1
		}

		fmt.Printf("%s\n", string(msg))

		functionClass, functionName, payload := clienetJSONParser(string(msg))

		fmt.Printf("%s :: %s :: %s\n", functionClass, functionName, payload)

		if functionClass == "login" {

			//var status = authenticationFunction(payload)

			if passwordCheckerWithPepper(string(payload)) {
				setConnectionHandlers(conn)
				conn.WriteMessage(1, []byte(clientJSONCreator("login", "1")))
				return 1
			}

			fmt.Printf("Attempt %d\n", attemps)

			attemps++
			conn.WriteMessage(1, []byte(clientJSONCreator("echo", fmt.Sprintf("Failed attempt: %d", attemps))))
			fmt.Printf("Failed Attempt %d\n", attemps)

			//Five failed attempts Special status
			if attemps >= MAX_ATTEMPTS { //|| (connection != nil) {
				conn.WriteMessage(1, []byte(clientJSONCreator("login", "2")))
				conn.Close()
				return -1
			}

		}

	}

}

func clientMessageRouterFunction(conn *websocket.Conn) int {

	// Read message from browser
	msgType, msg, err := conn.ReadMessage()
	if err != nil {
		return -1
	}

	_, _, payload := clienetJSONParser(string(msg))

	//fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

	if payload == "end" {
		conn.Close()
		return -1
	}

	fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), cleanText(payload))

	// Write message back to browser
	if err = conn.WriteMessage(msgType, []byte(clientJSONCreator("echo", cleanText(payload)))); err != nil {
		return 1
	}

	return 0

}

func setConnectionHandlers(conn *websocket.Conn) {

	conn.SetCloseHandler(custonSocketCloseHandler)

}

func custonSocketCloseHandler(code int, text string) error {

	var error2 error = nil

	//conn.Close()
	//fmt.Printf("User %d has disconnected\n", localvalue)

	return error2
}
