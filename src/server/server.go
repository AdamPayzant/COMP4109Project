package main

import (
	"context"
	"errors"
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
		return &pb_server.Status{Status: 2}, err
	}
	err = addUser(reg.GetUsername(), key, reg.GetIp())
	if err != nil {
		return &pb_server.Status{Status: 1}, err
	}
	return &pb_server.Status{Status: 0}, nil
}

func (s *server) getToken(ctx context.Context, req *pb_server.Username) (*pb_server.AuthKey, error) {
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

func (s *server) updateKey(ctx context.Context, req *pb_server.KeyUpdate) (*pb_server.Status, error) {
	validated, err := checkToken(req.Username, req.AuthKey)
	if !validated || err != nil {
		return &pb_server.Status{Status: 1}, err
	}

	key, err := x509.ParsePKCS1PublicKey(req.GetNewKey())
	if err != nil {
		return &pb_server.Status{Status: 3}, errors.New("Non key passed as key")
	}
	err = updateKey(req.Username, key)
	if err != nil {
		return &pb_server.Status{Status: 2}, err
	}

	return &pb_server.Status{Status: 0}, nil
}

// TODO: Implement the DB side function then this
func (s *server) searchUser(ctx context.Context, req *pb_server.UserQuery) (*pb_server.Status, error) {
	return nil, nil
}

func (s *server) getUser(ctx context.Context, req *pb_server.Username) (*pb_server.UserInfo, error) {
	ip, key, err := getUser(req.Username)
	if err != nil {
		return nil, err
	}
	return &pb_server.UserInfo{
		PublicKey: x509.MarshalPKCS1PublicKey(&key),
		IP:        ip,
	}, nil
}

func main() {
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
