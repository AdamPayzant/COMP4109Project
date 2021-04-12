package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
	"time"

	// pb_host "pb_host"

	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	_ "github.com/mattn/go-sqlite3"
)

type ClientHostSettings struct {
	ClientPublicKeyPath  string `json:"clientPublicKeyPath"`
	ServerCert           string `json:"serverCert"`
	ServerKey            string `json:"serverKey"`
	DB                   string `json:"DB"`
	ServerIP             string `json:"serverIP"`
	Username             string `json:"username"`
	CentrialServerIP     string `json:"centralServerIP"`
	CentrialServerCACert string `json:"centralServerCACert"`
	Token                string `json:"token"`
}

var settings ClientHostSettings

func tryLoadClientPublicKey(file string) *rsa.PublicKey {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatalf("Unable to open publicKey from path: %v", err)
	}

	raw, _ := ioutil.ReadFile(file)
	block, er := pem.Decode([]byte(raw))
	if err != nil {
		log.Fatalf("Unable to decode publicKey: %v", er)
	}
	key, e := bytesToKey(block.Bytes)
	if e != nil {
		log.Fatalln(e)
	}

	return key
}

func loadSettings(file string) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	e := json.Unmarshal(raw, &settings)
	if e != nil {
		log.Fatalf("Could not load settings: %v", e)
	}
}

func registerIfNeeded() {
	centralServer, e := getConnectionToCentralServer()
	if e != nil {
		log.Printf("Could not verify if user needs to be registered: %v", e)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		pb_server_userInfo, err := centralServer.server.GetUser(ctx, &pb_server.Username{Username: settings.Username})
		if err != nil {
			log.Printf("Failed to retrieve UserInfo: %v", err)
		}
		defer centralServer.conn.Close()

		if pb_server_userInfo == nil {
			keyBytes, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
			status, er := centralServer.server.Register(ctx, &pb_server.UserReg{Username: settings.Username, Key: keyBytes, Ip: settings.ServerIP})
			if er != nil {
				log.Printf("Failed to register User: %v", er)
			} else {
				if status.Status != 0 {
					log.Printf("Failed to register User: return state: %s", status)
				} else {
					log.Printf("Registed User!")
				}
			}
		}
	}
}

func main() {
	settingsPath := os.Args[1]

	initConnectionPool()

	loadSettings(settingsPath)
	clientPublicKey = tryLoadClientPublicKey(settings.ClientPublicKeyPath)
	connect(settings.DB)
	registerIfNeeded()
	startClientHost(settings.ServerIP)

	//closeConnectionPool()
	//closeConnectionToCentralServer()
	//db.Close()
}
