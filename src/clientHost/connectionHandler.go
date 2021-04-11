package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	pb_host "pb_host"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"
	pb_client "github.com/AdamPayzant/COMP4109Project/src/protos/smvsclient"
	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	_ "github.com/mattn/go-sqlite3"
)

type UserConnection struct {
	userInfo  *UserInfo
	user      pb_host.ClientHostClient
	conn      *grpc.ClientConn
	idleCount int
}

type ServerConnection struct {
	server pb_server.ServerClient
	conn   *grpc.ClientConn
}

type ClientConnection struct {
	client pb_client.ClientClient
	conn   *grpc.ClientConn
}

var centralServer *ServerConnection
var clientConnection *ClientConnection

const maxConnections int = 10

var connectionCount int = 0
var userConnection map[string]*UserConnection

func initConnectionPool() {
	userConnection = make(map[string]*UserConnection)
}

func closeConnectionPool() {
	for _, userConnection := range userConnection {
		userConnection.conn.Close()
	}
}

func getConnectionToClient() *ClientConnection {
	return clientConnection
}

func getConnectionToUser(username string) (*UserConnection, error) {
	var err error
	connection := userConnection[username]
	if connection == nil {
		if connectionCount > maxConnections {
			var greatestIdleCount = 0
			var mostIdled *UserConnection
			for _, userConnection := range userConnection {
				if userConnection.idleCount > greatestIdleCount {
					greatestIdleCount = userConnection.idleCount
					mostIdled = userConnection
				}
			}

			mostIdled.conn.Close()
			delete(userConnection, mostIdled.userInfo.username)
			connectionCount = connectionCount - 1
		}

		connection, err = connectToUser(username)
		if err == nil {
			userConnection[username] = connection
			connectionCount = connectionCount + 1
		}
	}

	connection.idleCount = 0
	return connection, err
}

func getConnectionToCentralServer() (*ServerConnection, error) {
	var err error
	if server == nil {
		centralServer, err = connectToCentralServer()
	}
	return centralServer, err
}

func connectToCentralServer() (*ServerConnection, error) {
	// Connects to the central server
	// Current uses self-signed TLS for this, I'd rather not go through a CA unless this is actually deployed

	pemServerCA, err := ioutil.ReadFile(settings.CentrialServerCACert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		InsecureSkipVerify: true, //Set to false when real CAs are being used
		RootCAs:            certPool,
	}

	conn, err := grpc.Dial(settings.CentrialServerIP, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return nil, err
	}

	return &ServerConnection{server: pb_server.NewServerClient(conn), conn: conn}, nil
}

func closeConnectionToCentralServer() {
	if clientConnection != nil {
		clientConnection.conn.Close()
		clientConnection = nil
	}
}

func connectToUser(user string) (*UserConnection, error) {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	connectedWithLocalInform := false
	var connection pb_host.ClientHostClient
	var conn *grpc.ClientConn
	var userInfo *UserInfo
	var err error

	userInfo, err = getUserInfo(user)
	if err == nil {
		conn, err = grpc.Dial(userInfo.ip, grpc.WithTransportCredentials(credentials.NewTLS(config)))
		if err == nil {
			connection = pb_host.NewClientHostClient(conn)
			connectedWithLocalInform = true
		}
	}

	if !connectedWithLocalInform {
		userInfo, err = getUserInfoFromSever(user)
		if err != nil {
			return nil, err
		}

		conn, err = grpc.Dial(userInfo.ip, grpc.WithTransportCredentials(credentials.NewTLS(config)))
		if err != nil {
			return nil, err
		}
		connection = pb_host.NewClientHostClient(conn)
	}

	_ = updateOrAddUserInfo(userInfo)

	return &UserConnection{userInfo: userInfo, user: connection, conn: conn, idleCount: 0}, nil
}

func connectToClient(ip string) error {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		return err
	}

	clientConnection = &ClientConnection{client: pb_client.NewClientClient(conn), conn: conn}

	return nil
}

func closeClientConnection() {
	if centralServer != nil {
		centralServer.conn.Close()
		centralServer = nil
	}
}

func getUserInfoFromSever(user string) (*UserInfo, error) {
	centralServer, e := getConnectionToCentralServer()
	ui, e := centralServer.server.GetUser(context.Background(), &pb_server.Username{Username: user})
	if e != nil {
		log.Printf("Could not connect to central server: %v", e)
		return nil, e
	}

	PKIXkey, err := x509.ParsePKIXPublicKey([]byte(ui.PublicKey))
	if err != nil {
		return nil, err
	}
	rsakey := PKIXkey.(*rsa.PublicKey)

	userInfo := &UserInfo{username: user, msgCount: 0, ip: ui.IP, key: rsakey}
	return userInfo, nil
}
