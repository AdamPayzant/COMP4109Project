package main

import (
	"log"
	"net/http"
)

func main() {

	//Servers Files from ./webView Directory
	http.Handle("/", http.FileServer(http.Dir("./webView")))

	//To connect to server use url https://127.0.0.1:3030
	//Will reject http requests
	log.Fatal(http.ListenAndServeTLS(":3030", "server.crt", "server.key", nil))

}
