package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//var remoteConnection = nil

//Using example for Websocket as start for client application
//WebSocket used to communicate between the CleintDevice and hostDevice
//At the moment it will be assumed that One client will run off of one host at a time.
//
// Src: https://gowebexamples.com/websockets/

func main() {

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.Handle("/", http.FileServer(http.Dir("./webView/")))
	//http.ListenAndServeTLS(":3030", "server.crt", "server.key", nil)
	log.Println("Serving at http://localhost:3030")
	http.ListenAndServe(":3030", nil)
}
