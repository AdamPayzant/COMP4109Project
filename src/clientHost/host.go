package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net"
	"time"

	pb_host "smvshost"
	pb_server "smvsserver"

	"google.golang.org/grpc"
)

var server pb_server.ServerClient = nil
var port = ":9090"

const (
	serverAddress = "localhost:8080"
	username      = "Tester"
)

type host struct {
	pb_host.UnimplementedClientHostServer
}

func (h *host) ReKey(ctx context.Context, req *pb_host.Token) (*pb_host.Status, error) {
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
	return nil, nil
}

func (h *host) InitializeConvo(ctx context.Context, req *pb_host.InitMessage) (*pb_host.Status, error) {
	return nil, nil
}

func (h *host) ConfirmConvo(ctx context.Context, req *pb_host.InitMessage) (*pb_host.Status, error) {
	return nil, nil
}

func (h *host) SendText(ctx context.Context, req *pb_host.ClientText) (*pb_host.Status, error) {
	return nil, nil
}

func (h *host) RecieveText(ctx context.Context, req *pb_host.H2HText) (*pb_host.Status, error) {
	return nil, nil
}

func (h *host) GetConversation(ctx context.Context, req *pb_host.Username) (*pb_host.Conversation, error) {
	return nil, nil
}

func main() {
	// Connects to the central server
	// Current uses self-signed TLS for this, I'd rather not go through a CA unless this is actually deployed
	config := &tls.Config{
		InsecureSkipVerify: false,
	}
	conn, err := grpc.Dial(serverAddress, grpc.WithTransportCredentials(crendentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	server = pb_server.NewServerClient(conn)

	// Starts up host server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb_host.RegisterClientHostServer(s, &host{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
