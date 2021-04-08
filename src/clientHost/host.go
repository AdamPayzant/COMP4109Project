package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"time"

	pb_host "pb_host"
	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"
	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type ClientHostSettings struct {
	PublicKeyPath    string `json:"publicKeyPath"`
	PrivateKeyPath   string `json:"privateKeyPath"`
	CertDir          string `json:"certDir"`
	DataBasePath     string `json:"dataBasePath"`
	ServerIP         string `json:"serverIP"`
	Username         string `json:"username"`
	CentrialServerIP string `json:"centrialServerIP"`
	TokenPath        string `json:"tokenPath"`
}

var settings ClientHostSettings

var db *sql.DB
var clientPublicKey *rsa.PublicKey
var clientPrivateKey *rsa.PrivateKey

var server pb_server.ServerClient = nil
var port = ":9090"

const (
	serverAddress = "localhost:8080"
	username      = "Tester"
)

type host struct {
	pb_host.UnimplementedClientHostServer
}

type UserInfo struct {
	name     string
	msgCount int
	ip       string
	key      *rsa.PublicKey
}

var userInfoCache map[string]*UserInfo

func updateUserInfoDatabase(user string, ip string, publicKey string) error {
	msgCount := 0
	rows, _ := db.Query("SELECT msgCount FROM userInfo WHERE user='" + user + "'")
	if rows.Next() {
		rows.Scan(&msgCount)
		db.Exec("DELETE FROM userInfo WHERE user='" + user + "'")
	}
	rows.Close()

	statement, e := db.Prepare("INSERT INTO userInfo (user, msgCount, ip, key) VALUES (?, ?, ?, ?)")
	if e != nil {
		return e
	}
	_, er := statement.Exec(user, msgCount, ip, publicKey)
	if er != nil {
		return er
	}

	return nil
}

func connectToCentralSever() (pb_server.ServerClient, *grpc.ClientConn, error) {
	// Connects to the central server
	// Current uses self-signed TLS for this, I'd rather not go through a CA unless this is actually deployed
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(settings.CentrialServerIP, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return nil, nil, err
	}

	return pb_server.NewServerClient(conn), conn, nil
}

func getUserInfoFromSever(user string) (bool, error) {
	server, conn, e := connectToCentralSever()
	ui, e := server.GetUser(context.Background(), &pb_server.Username{Username: user})
	if e != nil {
		return false, e
	}
	defer conn.Close()

	PKIXkey, err := x509.ParsePKIXPublicKey([]byte(ui.PublicKey))
	if err != nil {
		return false, err
	}
	rsakey := PKIXkey.(*rsa.PublicKey)

	e = updateUserInfoDatabase(user, ui.IP, string(ui.PublicKey))
	if e != nil {
		return false, e
	}

	msgCount := 0
	rows, _ := db.Query("SELECT msgCount FROM userInfo WHERE user='" + user + "'")
	if rows.Next() {
		rows.Scan(&msgCount)
	}
	rows.Close()

	userInfo := &UserInfo{name: user, msgCount: msgCount, ip: ui.IP, key: rsakey}
	userInfoCache[user] = userInfo

	return true, nil
}

func loadUserInfo(user string) (bool, error) {
	var name string
	var msgCount int
	var ip string
	var key string
	rows, _ := db.Query("SELECT user, msgCount, ip, key FROM userInfo WHERE user='" + user + "'")
	if rows.Next() {
		rows.Scan(&name, &msgCount, &ip, &key)
		rows.Close()

		PKIXkey, err := x509.ParsePKIXPublicKey([]byte(key))
		if err != nil {
			return false, err
		}
		rsakey := PKIXkey.(*rsa.PublicKey)

		userInfo := &UserInfo{name: name, msgCount: msgCount, ip: ip, key: rsakey}
		userInfoCache[user] = userInfo
	} else {
		rows.Close()
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
		if !hasLoad {
			return "", e
		}
	}
	ip = userInfoCache[user].ip
	return ip, nil
}

func getKey(user string) (*rsa.PublicKey, error) {
	var key *rsa.PublicKey
	if userInfoCache[user] == nil {
		hasLoad, e := loadUserInfo(user)
		if !hasLoad {
			return nil, e
		}
	}
	key = userInfoCache[user].key
	return key, nil
}

func connectToUser(user string) (pb_host.ClientHostClient, *grpc.ClientConn, error) {
	ip, e := getIp(user)
	if e != nil {
		return nil, nil, e
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return nil, nil, err
	}
	connection := pb_host.NewClientHostClient(conn)

	return connection, conn, nil
}

func RSA_OAEP_Encrypt(secretMessage string, key *rsa.PublicKey) (string, error) {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, key, []byte(secretMessage), label)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func RSA_OAEP_Decrypt(cipherText string, privKey *rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, privKey, ct, label)
	if err != nil {
		return "", err
	}
	// fmt.Println("Plaintext:", string(plaintext))
	return string(plaintext), nil
}

func encryptForClient(msg string) (string, error) {
	cypherText, err := RSA_OAEP_Encrypt(msg, clientPublicKey)
	if err != nil {
		return "", err
	}

	return cypherText, nil
}

func encryptForSending(msg string, user string) (string, error) {
	key, e := getKey(user)
	if e != nil {
		return "", e
	}

	cypherText, err := RSA_OAEP_Encrypt(msg, key)
	if err != nil {
		return "", err
	}

	return cypherText, nil
}

func decryptForClient(msg string) (string, error) {
	text, err := RSA_OAEP_Decrypt(msg, clientPrivateKey)
	if err != nil {
		return "", err
	}

	return text, nil
}

func generateSecret(user string) string {

	return "asdsad"
}

func verifySecret(user string, secret string) bool {
	return true
}

func authenticateClientToken(token string) bool {
	return true
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

// ####################################################################################################################

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
	keyBytes, _ := x509.MarshalPKIXPublicKey(&newkey.PublicKey)
	r, err := server.UpdateKey(ctx, &pb_server.KeyUpdate{Username: username,
		AuthKey: authToken,
		NewKey:  keyBytes})
	if err != nil {
		return &pb_host.Status{Status: 2}, err
	}
	return &pb_host.Status{Status: r.Status}, nil
}

func (h *host) DeleteMessage(ctx context.Context, req *pb_host.DeleteReq) (*pb_host.Status, error) {
	if authenticateClientToken(req.Token) {
		_, e := db.Exec("DELETE FROM conversations WHERE user='" + req.User + "' AND id='" + string(req.MessageID) + "'")
		if e != nil {
			return &pb_host.Status{Status: 2}, e
		}
	}

	return &pb_host.Status{Status: 0}, nil
}

func (h *host) SendText(ctx context.Context, req *pb_host.ClientText) (*pb_host.Status, error) {
	if req.TargetUser != username && authenticateClientToken(req.Token) {
		secret := generateSecret(req.TargetUser)
		connection, conn, e := connectToUser(req.TargetUser)
		if e != nil {
			log.Println(e)
			return &pb_host.Status{Status: 1}, e
		}

		var wholeMsg string
		for _, msg := range req.Message.Messages {
			wholeMsg = wholeMsg + msg
		}

		wholeMsgEncrypted, err := encryptForSending(wholeMsg, req.TargetUser)
		if err != nil {
			log.Println(err)
			return &pb_host.Status{Status: 1}, nil
		}

		chunckSize := 255
		lastIndex := 0
		msgLeng := len(wholeMsgEncrypted)
		segmentCount := int(math.Ceil(float64(msgLeng) / float64(chunckSize)))
		msgSegments := make([]string, segmentCount)
		i := 0
		for lastIndex < msgLeng {
			nextIndex := int(math.Min(float64(msgLeng), float64(lastIndex+chunckSize)))
			msgSegments[i] = wholeMsgEncrypted[lastIndex:nextIndex]
			lastIndex = nextIndex
			i = i + 1
		}

		startus, err := connection.RecieveText(context.Background(), &pb_host.H2HText{Message: &pb_host.ListofMessages{Messages: msgSegments}, User: req.TargetUser, Secret: secret})
		if err != nil {
			return startus, err
		}
		defer conn.Close()

		year, month, day, hour, minute, second := getTimeStamp()
		id := userInfoCache[req.TargetUser].msgCount

		statement, e := db.Prepare("INSERT INTO conversations (user, id, sender, year, month, day, hour, minute, second, msg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		wholeMsgClientEncrypted, errr := encryptForClient(wholeMsg)
		if errr != nil {
			log.Println(err)
			return &pb_host.Status{Status: 1}, nil
		}

		id = id + 1
		statement.Exec(req.TargetUser, id, true, year, month, day, hour, minute, second, wholeMsgClientEncrypted)
		db.Exec("UPDATE userInfo SET msgCount=" + strconv.Itoa(id) + " WHERE user='" + req.TargetUser + "'")
		return &pb_host.Status{Status: 0}, nil
	} else {
		return &pb_host.Status{Status: 1}, nil
	}
}

func (h *host) RecieveText(ctx context.Context, req *pb_host.H2HText) (*pb_host.Status, error) {
	if verifySecret(req.User, req.Secret) {
		year, month, day, hour, minute, second := getTimeStamp()
		id := -1

		rows, _ := db.Query("SELECT msgCount FROM userInfo WHERE user='" + req.User + "'")
		if rows.Next() {
			rows.Scan(&id)
		}
		rows.Close()

		statement, _ := db.Prepare("INSERT INTO conversations (user, id, sender, year, month, day, hour, minute, second, msg) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		var wholeMsg string
		for _, msg := range req.Message.Messages {
			wholeMsg = wholeMsg + msg
		}

		id = id + 1
		_, e := statement.Exec(req.User, id, false, year, month, day, hour, minute, second, wholeMsg)
		if e != nil {
			log.Fatalf("Error when adding to conversations table: %s", e)
		}

		db.Exec("UPDATE userInfo SET msgCount=" + strconv.Itoa(id) + " WHERE user='" + req.User + "'")

		return &pb_host.Status{Status: 0}, nil
	} else {
		return &pb_host.Status{Status: 1}, nil
	}
}

func (h *host) GetConversation(ctx context.Context, req *pb_host.Username) (*pb_host.Conversation, error) {
	if authenticateClientToken(req.Token) {

		type Message struct {
			Order       int    `json:"order"`
			Speaker     bool   `json:"speeker"`
			MessageText string `json:"messageText"`
		}

		type Convo struct {
			Messages []Message `json:"messages"`
		}

		var response *pb_host.Conversation

		rowCount, _ := db.Query("SELECT COUNT(*) AS msgCount FROM conversations WHERE user='" + req.Username + "'")
		if rowCount.Next() {
			var convo []Message
			var msgCount int

			rowCount.Scan(&msgCount)
			rowCount.Close()

			convo = make([]Message, msgCount)

			msgIndex := 0
			rows, _ := db.Query("SELECT id, sender, msg  FROM conversations WHERE user='" + req.Username + "'")
			for rows.Next() {
				var id int
				var sender bool
				var msg string
				rows.Scan(&id, &sender, &msg)
				decryptedMSG, err := decryptForClient(msg)
				if err != nil {
					log.Println(err)
				}
				convo[msgIndex].Order = id
				convo[msgIndex].Speaker = sender
				convo[msgIndex].MessageText = decryptedMSG
				msgIndex = msgIndex + 1
			}
			rows.Close()

			convoJsonBytes, _ := json.Marshal(&Convo{Messages: convo})
			convoJson := string(convoJsonBytes)

			chunckSize := 255
			lastIndex := 0
			convoJsonLeng := len(convoJson)
			segmentCount := int(math.Ceil(float64(convoJsonLeng) / float64(chunckSize)))
			convoJsonSegments := make([]string, segmentCount)
			i := 0
			for lastIndex < convoJsonLeng {
				nextIndex := int(math.Min(float64(convoJsonLeng), float64(lastIndex+chunckSize)))
				convoJsonSegments[i] = convoJson[lastIndex:nextIndex]
				lastIndex = nextIndex
				i = i + 1
			}

			response = &pb_host.Conversation{Convo: &pb_host.ListofMessages{Messages: convoJsonSegments}}
		} else {
			response = nil
			rowCount.Close()
		}

		return response, nil
	} else {
		return nil, nil
	}
}

// ####################################################################################################################

func initDB(file string) {
	var err error
	db, err = sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalln(err)
	}

	db.Exec("create table if not exists conversations (user text not null, id integer not null, sender boolean not null, year integer, month text, day integer, hour integer, minute integer, second integer, msg text not null, PRIMARY key(user, id))")
	db.Exec("create table if not exists userInfo (user text not null primary key, msgCount integer, ip integer not null, key text not null)")
}

func startClientHost(ip string) {
	serverCert, err := tls.LoadX509KeyPair("./certs/server-cert.pem", "./certs/server-key.pem")
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
		log.Fatalln("unable to open publicKey from path")
	}

	raw, _ := ioutil.ReadFile(file)
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		log.Fatalln("unable to decode publicKey")
	}
	key, e := x509.ParsePKIXPublicKey(block.Bytes)
	if e != nil {
		log.Fatalln(e)
	}

	return key.(*rsa.PublicKey)
}

func tryLoadClientPrivateKey(file string) *rsa.PrivateKey {
	_, err := os.Stat(file)
	if err != nil {
		log.Fatalln("unable to open privateKey from path")
	}

	raw, _ := ioutil.ReadFile(file)
	block, _ := pem.Decode([]byte(raw))
	if block == nil {
		log.Fatalln("unable to decode publicKey to request")
	}
	key, e := x509.ParsePKCS1PrivateKey(block.Bytes)
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

func decryptToken(token []byte) ([]byte, error) {
	// todo
	return token, nil
}

func getTokenFromServer() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	pd_server_token, e := server.GetToken(ctx, &pb_server.Username{Username: settings.Username})
	if e != nil {
		log.Printf("Failed to retrieve token: %v", e)
		return nil, e
	}

	token, _ := decryptToken(pd_server_token.AuthKey)
	return token, nil
}

func registerOrUpdateUserIfNeeded() {
	server, conn, e := connectToCentralSever()
	if e != nil {
		log.Printf("Could verify if user needs to be registered: %v", e)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		pb_server_userInfo, err := server.GetUser(ctx, &pb_server.Username{Username: settings.Username})
		if err != nil {
			log.Printf("Failed to retrieve UserInfo: %v", err)
		}
		defer conn.Close()

		if pb_server_userInfo == nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			keyBytes, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
			status, er := server.Register(ctx, &pb_server.UserReg{Username: settings.Username, Key: keyBytes, Ip: settings.ServerIP})
			if er != nil {
				log.Printf("Failed to register User: %v", er)
			}

			if status.Status != 0 {
				log.Printf("Failed to register User: return state: %s", status)
			}
		} else {
			token, _ := getTokenFromServer()
			PKIXkey, _ := x509.ParsePKIXPublicKey(pb_server_userInfo.PublicKey)
			key := PKIXkey.(*rsa.PublicKey)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			clientPublicKeyBytes, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
			if key == nil || !bytes.Equal(pb_server_userInfo.PublicKey, clientPublicKeyBytes) {
				newKey, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
				key_status, key_e := server.UpdateKey(ctx, &pb_server.KeyUpdate{Username: username,
					AuthKey: token,
					NewKey:  newKey})
				if key_e != nil {
					log.Printf("Failed to update Key: return state: %s", key_status)
				}
			}

			if pb_server_userInfo.IP != settings.ServerIP {
				ip_status, ip_e := server.UpdateIP(ctx, &pb_server.IPupdate{Username: username,
					AuthKey: token,
					NewIP:   settings.ServerIP})
				if ip_e != nil {
					log.Printf("Failed to update IP: return state: %s", ip_status)
				}
			}
		}
	}
}

func main() {
	settingsPath := os.Args[1]
	userInfoCache = make(map[string]*UserInfo)

	loadSettings(settingsPath)
	clientPublicKey = tryLoadClientPublicKey(settings.PublicKeyPath)
	clientPrivateKey = tryLoadClientPrivateKey(settings.PrivateKeyPath)
	// registerOrUpdateUserIfNeeded()
	initDB(settings.DataBasePath)

	keyBytes, _ := x509.MarshalPKIXPublicKey(clientPublicKey)
	updateUserInfoDatabase("Tester", ":8080", string(keyBytes))
	updateUserInfoDatabase("Tester1", ":7070", string(keyBytes))

	startClientHost(settings.ServerIP)

	db.Close()
}
