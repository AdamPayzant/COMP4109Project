package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	pb_host "pb_host"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"
	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	_ "github.com/mattn/go-sqlite3"
)

type ClientHostSettings struct {
	ClientPublicKeyPath  string `json:"clientPublicKeyPath"`
	ServerCert           string `json:"serverCert"`
	ServerKey            string `json:"serverKey"`
	DB                   string `json:"DB"`
	ServerIP             string `json:"serverIP"`
	Username             string `json:"username"`
	CentrialServerIP     string `json:"centrialServerIP"`
	CentrialServerCACert string `json:"centrialServerCACert"`
}

var settings ClientHostSettings

func startClientHost(ip string) {
	serverCert, err := tls.LoadX509KeyPair(settings.ServerCert, settings.ServerKey)
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	lis, err := net.Listen("tcp", ip)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))
	pb_host.RegisterClientHostServer(s, &host{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}

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
	key, e := x509.ParsePKIXPublicKey(block.Bytes)
	if e != nil {
		log.Fatalln(e)
	}

	return key.(*rsa.PublicKey)
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

func decryptToken(token []byte) ([]byte, error) {
	// todo
	return token, nil
}

func getTokenFromServer() ([]byte, error) {
	centralServer, err := getConnectionToCentralServer()
	if err != nil {
		log.Printf("Failed to get Token from Server: %v", err)
	}
	defer centralServer.conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pd_server_token, e := centralServer.server.GetToken(ctx, &pb_server.Username{Username: settings.Username})
	if e != nil {
		log.Printf("Failed to retrieve token: %v", e)
		return nil, e
	}

	token, _ := decryptToken(pd_server_token.AuthKey)
	return token, nil
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

// else {
// 	token, _ := getTokenFromServer()
// 	PKIXkey, err := x509.ParsePKIXPublicKey(pb_server_userInfo.PublicKey)
// 	if err != nil {
// 		log.Println("Failed to parse public key from  centrial Server: %v", err)
// 	}
// 	key := PKIXkey.(*rsa.PublicKey)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	clientPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
// 	if key == nil || !bytes.Equal(pb_server_userInfo.PublicKey, clientPublicKeyBytes) {
// 		newKey, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
// 		key_status, key_e := server.UpdateKey(ctx, &pb_server.KeyUpdate{Username: username,
// 			AuthKey: token,
// 			NewKey:  newKey})
// 		if key_e != nil {
// 			log.Printf("Failed to update Key: return state: %s", key_status)
// 		} else {
// 			if key_status.Status == 0 {
// 				log.Printf("Updated Key!")
// 			}
// 		}
// 	}

// 	if pb_server_userInfo.IP != settings.ServerIP {
// 		ip_status, ip_e := server.UpdateIP(ctx, &pb_server.IPupdate{Username: username,
// 			AuthKey: token,
// 			NewIP:   settings.ServerIP})
// 		if ip_e != nil {
// 			log.Printf("Failed to update IP: return state: %s", ip_status)
// 		} else {
// 			if ip_status.Status == 0 {
// 				log.Printf("Updated IP!")
// 			}
// 		}
// 	}
// }

func main() {
	settingsPath := os.Args[1]

	initConnectionPool()

	loadSettings(settingsPath)
	clientPublicKey = tryLoadClientPublicKey(settings.ClientPublicKeyPath)
	connect(settings.DB)
	registerIfNeeded()
	startClientHost(settings.ServerIP)

	closeConnectionPool()
	db.Close()
}
