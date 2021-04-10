package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"crypto/rsa"
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"
)

const (
	port = ":9090"
)

type server struct {
	pb_server.UnimplementedServerServer
}

func (s *server) Register(ctx context.Context, reg *pb_server.UserReg) (*pb_server.Status, error) {
	PKIXkey, err := x509.ParsePKIXPublicKey(reg.GetKey())
	if err != nil {
		return &pb_server.Status{Status: 2}, err
	}
	key := PKIXkey.(*rsa.PublicKey)
	err = addUser(reg.GetUsername(), key, reg.GetIp())
	if err != nil {
		return &pb_server.Status{Status: 1}, err
	}
	return &pb_server.Status{Status: 0}, nil
}

func (s *server) GetToken(ctx context.Context, req *pb_server.Username) (*pb_server.AuthKey, error) {

	token, err := addToken(req.Username)
	if err != nil {
		return nil, err
	}

	return &pb_server.AuthKey{AuthKey: token}, nil
}

func (s *server) UpdateIP(ctx context.Context, req *pb_server.IPupdate) (*pb_server.Status, error) {
	validated, err := checkToken(req.Username, req.AuthKey)
	if !validated || err != nil {
		return &pb_server.Status{Status: 1}, err
	}

	err = updateIP(req.Username, req.NewIP)
	if err != nil {
		return &pb_server.Status{Status: 2}, err
	}

	return &pb_server.Status{Status: 0}, nil
}

func (s *server) UpdateKey(ctx context.Context, req *pb_server.KeyUpdate) (*pb_server.Status, error) {
	validated, err := checkToken(req.Username, req.AuthKey)
	if !validated || err != nil {
		return &pb_server.Status{Status: 1}, err
	}

	PKIXkey, err := x509.ParsePKIXPublicKey(req.GetNewKey())
	if err != nil {
		return &pb_server.Status{Status: 3}, errors.New("non key passed as key")
	}
	key := PKIXkey.(*rsa.PublicKey)
	err = updateKey(req.Username, key)
	if err != nil {
		return &pb_server.Status{Status: 2}, err
	}

	return &pb_server.Status{Status: 0}, nil
}

// TODO: Implement the DB side function then this
func (s *server) SearchUser(ctx context.Context, req *pb_server.UserQuery) (*pb_server.UserList, error) {
	return nil, nil
}

func (s *server) GetUser(ctx context.Context, req *pb_server.Username) (*pb_server.UserInfo, error) {
	ip, key, err := getUser(req.Username)
	if err != nil {
		return nil, err
	}
	return &pb_server.UserInfo{
		PublicKey: key,
		IP:        ip,
	}, nil
}

func main() {
	fmt.Println("Starting server")

	creds, err := credentials.NewServerTLSFromFile("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	//s := grpc.NewServer()
	pb_server.RegisterServerServer(s, &server{})

	err = connect()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
