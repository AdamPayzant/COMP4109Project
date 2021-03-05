package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"crypto/x509"

	"google.golang.org/grpc"

	pb_server "smvsserver"
)

const (
	port = ":8080"
)

type server struct {
	pb_server.UnimplementedServerServer
}

func (s *server) Register(ctx context.Context, reg *pb_server.UserReg) (*pb_server.Status, error) {
	key, err := x509.ParsePKCS1PublicKey(reg.GetKey())
	if err != nil {
		return 2, err
	}
	err = addUser(reg.GetUsername(), key, reg.GetIp())
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func (s *server) getToken(ctx context.Context, req *pb_server.Username) (*pb_server.AuthKey, error) {
	return nil, nil
}

func (s *server) UpdateIP(ctx context.Context, req *pb_server.IPupdate) (*pb_server.Status, error) {
	return nil, nil
}

func (s *server) updateKey(ctx context.Context, req *pb_server.KeyUpdate) (*pb_server.Status, error) {
	return nil, nil
}

func (s *server) searchUser(ctx context.Context, req *pb_server.UserQuery) (*pb_server.Status, error) {
	return nil, nil
}

func (s *server) getUser(ctx context.Context, req *pb_server.Username) (*pb_server.UserInfo, error) {
	return nil, nil
}

func main() {
	fmt.Println("It's working")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb_server.RegisterServerServer(s, &server{})

	err = connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
