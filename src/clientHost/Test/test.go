package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
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

func main() {

	type Message struct {
		Order       int    `json:"order"`
		Speaker     bool   `json:"speaker"`
		MessageText string `json:"messageText"`
	}

	type Convo struct {
		Messages []Message `json:"messages"`
	}

	config := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(credentials.NewTLS(config)))
	// conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	server = pb_host.NewClientHostClient(conn)

	//SEND TEST
	response, err := server.SendText(context.Background(), &pb_host.ClientText{TargetUser: "Default1", Message: &pb_host.ListofMessages{Messages: []string{"test", "test", "test"}}, Token: "TEst"})
	if err != nil {
		log.Fatalf("Error when calling SendText: %s", err)
	}
	log.Printf("Response from server: %s", response.Status)

	//GET TEST
	res, er := server.GetConversation(context.Background(), &pb_host.Username{Username: "Default1", Token: "TEst"})
	if er != nil {
		log.Fatalf("Error when calling GetConversation: %s", er)
	}

	var convo Convo
	var wholeMSG string
	for _, msg := range res.Convo.Messages {
		wholeMSG = wholeMSG + msg
	}

	json.Unmarshal([]byte(wholeMSG), &convo)

	for _, m := range convo.Messages {
		fmt.Println(m.Order)
	}

	messageID := int64(convo.Messages[0].Order)

	//DELETE TEST
	response, err = server.DeleteMessage(context.Background(), &pb_host.DeleteReq{User: "Default1", MessageID: messageID, Token: "TEst"})
	if err != nil {
		log.Fatalf("Error when calling DeleteMessage: %s", err)
	}
	log.Printf("Response from server: %s", response.Status)

	//GET TEST
	res, er = server.GetConversation(context.Background(), &pb_host.Username{Username: "Default1", Token: "TEst"})
	if er != nil {
		log.Fatalf("Error when calling GetConversation: %s", er)
	}

	for _, msg := range res.Convo.Messages {
		wholeMSG = wholeMSG + msg
	}

	json.Unmarshal([]byte(wholeMSG), &convo)

	for _, m := range convo.Messages {
		fmt.Println(m.Order)
	}

	//ca, err := server.GetCA(context.Background(), &pb_host.Empty{})
	//fmt.Println(ca)
	// err1 := ioutil.WriteFile("./test/ca-cert.pem", ca.Ca, 0644)
	// if err1 != nil {
	// 	panic(err1)
	// }
}
