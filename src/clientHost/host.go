package main

//"encoding/json"
import (
	"log"
	"net/http"
)

//Using example for Websocket as start for client application
//WebSocket used to communicate between the CleintDevice and hostDevice
//At the moment it will be assumed that One client will run off of one host at a time.
//
// Src: https://gowebexamples.com/websockets/

//var index int = 0

//var connection *websocket.Conn = nil

func main() {

	http.HandleFunc("/echo", websocketRequestHander)
	/*
		go func() {

			time.Sleep(9 * time.Second)

			if connection != nil {

				fmt.Printf("skyrim %t %t %t %t \n", passwordChecker("12345678"), passwordChecker("badpassword"), validHashCheck("1"), validHashCheck("ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"))

				connection.WriteMessage(1, []byte("Finaly awake..."))

			}
		}()
	*/
	http.Handle("/", http.FileServer(http.Dir("./webView/")))
	//http.ListenAndServeTLS(":3030", "server.crt", "server.key", nil)
	log.Println("Serving at http://localhost:3030")
	http.ListenAndServe(":3030", nil)
}
