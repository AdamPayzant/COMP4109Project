package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"

	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"
	"google.golang.org/grpc"
)

func testServer() bool {
	var serverAddress = "localhost:50051"

	/*
		config := &tls.Config{
			InsecureSkipVerify: false,
		}
		conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	*/
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	fmt.Println("Connected")
	defer conn.Close()
	server := pb_server.NewServerClient(conn)

	fmt.Println("Testing register")
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	var res = false
	res = testRegister(server, key)
	if !res {
		fmt.Println("Failed to register user")
		return false
	}
	fmt.Println("Passed")
	fmt.Println("Testing get Token")
	res, token := testGetToken(server)
	if !res {
		fmt.Printf("Failed to generate token")
		return false
	}
	fmt.Println("Passed")
	fmt.Println("Testing update IP")
	tok, _ := rsa.DecryptOAEP(sha512.New(), rand.Reader, key, token, nil)
	res = testUpdateIP(server, tok)
	if !res {
		fmt.Println("Failed to update IP")
		return false
	}
	fmt.Println("Passed")
	fmt.Println("Testing update key")
	newKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	res = testUpdateKey(server, tok, newKey)
	if !res {
		fmt.Println("Failed to update key")
		return false
	}
	fmt.Println("Passed")
	fmt.Println("Testing get user")
	res = testGetUser(server, newKey)
	if !res {
		fmt.Println("Failed to get user")
		return false
	}
	fmt.Println("Passed")
	return true
}

func testRegister(server pb_server.ServerClient, key *rsa.PrivateKey) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb_server.UserReg{
		Username: "Tester",
		Key:      x509.MarshalPKCS1PublicKey(&key.PublicKey),
		Ip:       "1.1.1.1:1111",
	}
	r, err := server.Register(ctx, &req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if r.Status != 0 {
		return false
	}
	return true
}
func testGetToken(server pb_server.ServerClient) (bool, []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb_server.Username{
		Username: "Tester",
	}
	r, err := server.GetToken(ctx, &req)
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, r.AuthKey
}
func testUpdateIP(server pb_server.ServerClient, token []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb_server.IPupdate{
		Username: "Tester",
		AuthKey:  token,
		NewIP:    "1.1.1.1:1112",
	}
	r, err := server.UpdateIP(ctx, &req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if r.Status != 0 {
		return false
	}
	return true
}
func testUpdateKey(server pb_server.ServerClient, token []byte, newKey *rsa.PrivateKey) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb_server.KeyUpdate{
		Username: "Tester",
		AuthKey:  token,
		NewKey:   x509.MarshalPKCS1PublicKey(&newKey.PublicKey),
	}
	r, err := server.UpdateKey(ctx, &req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if r.Status != 0 {
		return false
	}
	return true
}
func testGetUser(server pb_server.ServerClient, key *rsa.PrivateKey) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := pb_server.Username{
		Username: "Tester",
	}
	r, err := server.GetUser(ctx, &req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if r.IP != "1.1.1.1:1112" {
		fmt.Println("Incorrect IP")
		return false
	}
	if !bytes.Equal(r.PublicKey, x509.MarshalPKCS1PublicKey(&key.PublicKey)) {
		fmt.Println("Incorrect key")
		return false
	}
	return true
}

func main() {
	fmt.Println("Testing server:")
	res := testServer()
	if res {
		fmt.Println("Server has passed test")
	}
}
