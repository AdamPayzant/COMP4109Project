//https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go
package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	//sqlite3 "github.com/mattn/go-sqlite3"
)

//Tried to impliment the example program from the repository used for socket.io go library used
//Geting an error that indicates that the server (this) is getting packets from the client
//The connection never starts

//Will be looking into an alternative library
//Example from: https://github.com/googollee/go-socket.io/tree/master/_example

func main() {

	server := socketio.NewServer(nil)

	//Socket IO work
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	//defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./webView/")))
	log.Println("Serving at http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}
