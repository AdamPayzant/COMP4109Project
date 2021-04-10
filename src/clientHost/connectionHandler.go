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
	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	_ "github.com/mattn/go-sqlite3"
)

type ClientConnection struct {
	userInfo  *UserInfo
	client    pb_host.ClientHostClient
	conn      *grpc.ClientConn
	idleCount int
}

type ServerConnection struct {
	server pb_server.ServerClient
	conn   *grpc.ClientConn
}

var centralServer *ServerConnection

const maxConnections int = 10

var connectionCount int = 0
var clientConnections map[string]*ClientConnection

func initConnectionPool() {
	clientConnections = make(map[string]*ClientConnection)
}

func closeConnectionPool() {
	for _, userConnection := range clientConnections {
		userConnection.conn.Close()
	}
}

func getConnectionToUser(username string) (*ClientConnection, error) {
	var err error
	connection := clientConnections[username]
	if connection == nil {
		if connectionCount > maxConnections {
			var greatestIdleCount = 0
			var mostIdled *ClientConnection
			for _, userConnection := range clientConnections {
				if userConnection.idleCount > greatestIdleCount {
					greatestIdleCount = userConnection.idleCount
					mostIdled = userConnection
				}
			}

			mostIdled.conn.Close()
			delete(clientConnections, mostIdled.userInfo.username)
			connectionCount = connectionCount - 1
		}

		connection, err = connectToUser(username)
		if err == nil {
			clientConnections[username] = connection
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

func connectToUser(user string) (*ClientConnection, error) {
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

	return &ClientConnection{userInfo: userInfo, client: connection, conn: conn, idleCount: 0}, nil
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
