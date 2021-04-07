package main

import (
	"context"
	"crypto/tls"
	"log"
	pb_host "pb_host"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type UserInfo struct {
	name string
	ip   string
	key  string
}

var userInfoCache map[string]UserInfo

var server pb_host.ClientHostClient = nil
var port = ":9090"

const (
	serverAddress = "localhost:9090"
	username      = "Tester"
)

type host struct {
	pb_host.UnimplementedClientHostServer
}

type cahost struct {
	pb_host.UnimplementedClientCAHostServer
}

func main() {
	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(":9090", grpc.WithTransportCredentials(credentials.NewTLS(config)))
	// conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	server = pb_host.NewClientHostClient(conn)

	response, err := server.RecieveText(context.Background(), &pb_host.H2HText{Message: &pb_host.ListofMessages{Messages: []string{"test", "test", "test"}}, User: "Me", Secret: "TEst"})
	if err != nil {
		log.Fatalf("Error when calling InitializeConvo: %s", err)
	}
	log.Printf("Response from server: %s", response.Status)

	// ca, err := server.GetCA(context.Background(), &pb_host.Empty{})
	// fmt.Println(ca)
	// err1 := ioutil.WriteFile("./test/ca-cert.pem", ca.Ca, 0644)
	// if err1 != nil {
	// 	panic(err1)
	// }
}
