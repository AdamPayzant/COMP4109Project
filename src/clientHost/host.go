package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	// pb_host "pb_host"
	pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"
	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var encryptMSG bool
var clientPublicKey *rsa.PublicKey
var hostPrivateKey *rsa.PrivateKey
var clientPrivateKey *rsa.PrivateKey

type UserInfo struct {
	name       string
	ip         string
	key        *rsa.PublicKey
	connection *pb_host.ClientHostClient
}

var userInfoCache map[string]*UserInfo

var server pb_server.ServerClient = nil
var port = ":9090"

const (
	serverAddress = "localhost:8080"
	username      = "Tester"
)

type host struct {
	pb_host.UnimplementedClientHostServer
}

func RSA_OAEP_Encrypt(secretMessage string, key rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	if err != nil {
		log.Fatalf("Failed to EncryptOAEP: %v", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func RSA_OAEP_Decrypt(cipherText string, privKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		log.Fatalf("Failed to DecryptOAEP: %v", err)
	}
	// fmt.Println("Plaintext:", string(plaintext))
	return string(plaintext)
}

func genSecret(user string) string {
	return "asdasdsad"
}

func verifySecret(user string, secret string) bool {
	return true
}

func getUserInfoFromSever(user string) (bool, error) {
	ui, e := server.GetUser(context.Background(), &pb_server.Username{Username: user})
	if e != nil {
		return false, e
	}

	rsakey, err := x509.ParsePKCS1PublicKey([]byte(ui.PublicKey))
	if err != nil {
		panic(err)
	}

	userInfo := &UserInfo{name: user, ip: ui.IP, key: rsakey}
	userInfoCache[user] = userInfo

	rows, _ := db.Query("SELECT user, ip, key FROM conversations WHERE user='" + user + "'")
	if rows.Next() {
		db.Exec("DELETE FROM userInfo WHERE user='" + user + "'")
	}

	statement, _ := db.Prepare("INSERT INTO userInfo (user, ip, key) VALUES (?, ?, ?)")
	_, er := statement.Exec(user, ui.IP, string(ui.PublicKey))
	if er != nil {
		return false, er
	}

	return true, nil
}

func loadUserInfo(user string) (bool, error) {
	var name string
	var ip string
	var key string
	rows, _ := db.Query("SELECT user, ip, key FROM conversations WHERE user='" + user + "'")
	if rows.Next() {
		rows.Scan(&name, &ip, &key)

		rsakey, err := x509.ParsePKCS1PublicKey([]byte(key))
		if err != nil {
			panic(err)
		}

		userInfo := &UserInfo{name: name, ip: ip, key: rsakey}
		userInfoCache[user] = userInfo
	} else {
		_, e := getUserInfoFromSever(name)
		if e != nil {
			return false, e
		}
	}
	return true, nil
}

func getIp(user string) (string, error) {
	var ip string
	if userInfoCache[user] == nil {
		hasLoad, e := loadUserInfo(user)
		if hasLoad {
			ip = userInfoCache[user].ip
		} else {
			return "", e
		}
	}
	ip = userInfoCache[user].ip
	return ip, nil
}

func connectToUser(user string) (*pb_host.ClientHostClient, error) {
	if userInfoCache[user] == nil {
		ip, e := getIp(user)
		if e != nil {
			return nil, e
		}

		config := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := grpc.Dial(ip+":9090", grpc.WithTransportCredentials(credentials.NewTLS(config)))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		}
		defer conn.Close()
		connection := pb_host.NewClientHostClient(conn)
		userInfoCache[user].connection = &connection
	}

	return userInfoCache[user].connection, nil
}

func getTimeStamp() (int, string, int, int, int, int) {
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	var hour int
	var minute int
	var second int
	hour, minute, second = now.Clock()

	return now.Year(), now.Month().String(), now.Day(), hour, minute, second
}

/*
	Message pipe line
		Text
			Sending: plain text -> RSA -> TLS/SSL (gRPC)
			Receving: TLS/SSL (gRPC) -> RSA -> plain text
		Video
			Sending: stream -> TLS/SSL (gRPC)
			Receving: TLS/SSL (gRPC) -> steam
*/

func (h *host) ReKey(ctx context.Context, req *pb_host.Token) (*pb_host.Status, error) {
	/*
		This function is used to update the public key for the RSA encryption.
		The public key can only be changed if the correct auth key is provided to the main server
		The private key should only exist on TRUSTED end user clients.
	*/

	/*
		This is currently just in a state to demo gRPC call
		Plenty of stuff still to do
		TODO:
			Implement key management system
			Implement token management system
	*/
	newkey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return &pb_host.Status{Status: 1}, errors.New("Key gen error")
	}
	// Get authtoken
	authToken := make([]byte, 64)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := server.UpdateKey(ctx, &pb_server.KeyUpdate{Username: username,
		AuthKey: authToken,
		NewKey:  x509.MarshalPKCS1PublicKey(&newkey.PublicKey)})
	if err != nil {
		return &pb_host.Status{Status: 2}, err
	}
	return &pb_host.Status{Status: r.Status}, nil
}

func (h *host) DeleteMessage(ctx context.Context, req *pb_host.DeleteReq) (*pb_host.Status, error) {
	if auth(req.Token) {
		_, e := db.Exec("DELETE FROM conversations WHERE user='" + req.User + "' AND id='" + string(req.MessageID) + "'")
		if e != nil {
			return &pb_host.Status{Status: 2}, e
		}
	}

	return &pb_host.Status{Status: 0}, nil
}

func (h *host) InitializeConvo(ctx context.Context, req *pb_host.InitMessage) (*pb_host.Status, error) {
	fmt.Println("Test")
	return &pb_host.Status{Status: 0}, nil
}

func (h *host) ConfirmConvo(ctx context.Context, req *pb_host.InitMessage) (*pb_host.Status, error) {
	return nil, nil
}

func (h *host) SendText(ctx context.Context, req *pb_host.ClientText) (*pb_host.Status, error) {
	if auth(req.Token) {
		secret := genSecret(req.TargetUser)
		connection, e := connectToUser(req.TargetUser)
		if e != nil {
			return &pb_host.Status{Status: 1}, e
		}

		sendMSGs := make([]string, len(req.Message.Messages))
		for i, msg := range req.Message.Messages {
			sendMSGs[i] = RSA_OAEP_Encrypt(msg, *userInfoCache[req.TargetUser].key)
		}

		startus, err := (*connection).RecieveText(context.Background(), &pb_host.H2HText{Message: &pb_host.ListofMessages{Messages: sendMSGs}, User: req.TargetUser, Secret: secret})
		if err != nil {
			return startus, e
		}

		year, month, day, hour, minute, second := getTimeStamp()
		var id int

		rows, _ := db.Query("SELECT MAX(id) AS maxID FROM conversations WHERE user='" + req.TargetUser + "'")
		if rows != nil && rows.Next() {
			rows.Scan(&id)
		}
		rows.Close()

		statement, _ := db.Prepare("INSERT INTO conversations (user, id, to, year, month, day, hour, minute, second, msg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		for _, msg := range req.Message.Messages {
			if encryptMSG {
				msg = RSA_OAEP_Encrypt(msg, *clientPublicKey)
			}
			id = id + 1
			statement.Exec(req.TargetUser, id, true, year, month, day, hour, minute, second, msg)
		}
		return &pb_host.Status{Status: 0}, nil
	} else {
		return &pb_host.Status{Status: 1}, nil
	}
}

func (h *host) RecieveText(ctx context.Context, req *pb_host.H2HText) (*pb_host.Status, error) {
	if verifySecret(req.User, req.Secret) {
		year, month, day, hour, minute, second := getTimeStamp()
		id := -1

		rows, _ := db.Query("SELECT MAX(id) AS maxID FROM conversations WHERE user='" + req.User + "'")
		if rows != nil && rows.Next() {
			rows.Scan(&id)
		}
		rows.Close()

		statement, _ := db.Prepare("INSERT INTO conversations (user, id, sender, year, month, day, hour, minute, second, msg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		for _, msg := range req.Message.Messages {
			if encryptMSG {
				msg = RSA_OAEP_Encrypt(msg, *clientPublicKey)
			}

			id = id + 1
			fmt.Println(msg)
			_, e := statement.Exec(req.User, id, false, year, month, day, hour, minute, second, msg)
			if e != nil {
				log.Fatalf("Error when adding to conversations table: %s", e)
			}

		}

		return &pb_host.Status{Status: 0}, nil
	} else {
		return &pb_host.Status{Status: 1}, nil
	}
}

func (h *host) GetConversation(ctx context.Context, req *pb_host.Username) (*pb_host.Conversation, error) {
	return nil, nil
}

func auth(token string) bool {
	return true
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./data.db")
	if err != nil {
		fmt.Println(err)
	}

	db.Exec("create table if not exists conversations (user text not null, id integer not null, sender boolean not null, year integer, month text, day integer, hour integer, minute integer, second integer, msg text not null, PRIMARY key(user, id))")
	db.Exec("create table if not exists userInfo (user text not null primary key, ip integer not null, key text not null)")
}

func startClientHost() {
	serverCert, err := tls.LoadX509KeyPair("./certs/server-cert.pem", "./certs/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(credentials.NewTLS(config)))
	pb_host.RegisterClientHostServer(s, &host{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}

func connectToCentralSever() {
	// Connects to the central server
	// Current uses self-signed TLS for this, I'd rather not go through a CA unless this is actually deployed
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	server = pb_server.NewServerClient(conn)
}

func tryLoadPublicKey() {
	_, err := os.Stat("./client_public.pem")

	if err == nil {
		raw, _ := ioutil.ReadFile("./client_public.pem")
		block, _ := pem.Decode([]byte(raw))
		if block == nil {
			fmt.Println("unable to decode publicKey to request")
		}
		key, e := x509.ParsePKIXPublicKey(block.Bytes)
		if e != nil {
			log.Fatalf("%v", e)
		}

		clientPublicKey = key.(*rsa.PublicKey)
		encryptMSG = true
	}
}

func main() {
	encryptMSG = false
	userInfoCache = make(map[string]*UserInfo)
	// generateClientKeys()
	tryLoadPublicKey()
	initDB()

	// connectToCentralSever()

	startClientHost()
	db.Close()
}
