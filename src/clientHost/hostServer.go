package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math"
	"time"

	pb_host "pb_host"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"
	pb_client "github.com/AdamPayzant/COMP4109Project/src/protos/smvsclient"
	pb_server "github.com/AdamPayzant/COMP4109Project/src/protos/smvsserver"

	_ "github.com/mattn/go-sqlite3"
)

type host struct {
	pb_host.UnimplementedClientHostServer
}

var server pb_server.ServerClient = nil

func stringToListofMessages(msg string) *pb_host.ListofMessages {
	msgLen := len(msg)
	if msgLen > 0 {
		chunckSize := 255
		lastIndex := 0
		segmentCount := int(math.Ceil(float64(msgLen) / float64(chunckSize)))
		msgSegments := make([]string, segmentCount)
		i := 0
		for i < segmentCount {
			nextIndex := int(math.Min(float64(msgLen), float64(lastIndex+chunckSize)))
			msgSegments[i] = msg[lastIndex:nextIndex]
			lastIndex = nextIndex
			i = i + 1
		}

		return &pb_host.ListofMessages{Messages: msgSegments}
	} else {
		return nil
	}
}

func listofMessagesToString(list *pb_host.ListofMessages) string {
	message := ""
	for _, part := range list.Messages {
		message = message + part
	}
	return message
}

func storeMesssage(user string, speaker bool, msg string) {
	res, err := addMessage(user, &Message{Speaker: speaker, MessageText: msg})

	if err != nil {
		log.Printf("A message was not saved!: %v", err)
	} else {
		count, errr := res.RowsAffected()
		if errr != nil {
			log.Fatalf("Connot get results to adding message to local db!: %v", errr)
		}

		if count <= 0 {
			log.Printf("A message was not saved!")
		}
	}
}

func forawrdMessageToClient(msg *pb_host.ListofMessages) {
	client := getConnectionToClient()
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		status, err := client.client.RecieveMessage(ctx, &pb_client.ListofMessages{Messages: msg.Messages})
		if err != nil {
			log.Printf("Could not forward message: %v \n", err)
		}

		if status.Status != 0 {
			log.Println("Could not forward message")
		}
	}
}

func LogIn(ctx context.Context, req *pb_host.ClientInfo) (*pb_host.Status, error) {
	if req.Username != settings.Username {
		return &pb_host.Status{Status: 1}, errors.New("User does not match")
	}

	if !authenticateClientToken(req.Token) {
		return &pb_host.Status{Status: 1}, errors.New("Failed to authenticate Token")
	}

	err := connectToClient(req.Ip)
	if err != nil {
		return &pb_host.Status{Status: 1}, err
	}

	return &pb_host.Status{Status: 0}, nil
}

func LogOut(ctx context.Context, req *pb_host.ClientInfo) (*pb_host.Status, error) {
	if req.Username != settings.Username {
		return &pb_host.Status{Status: 1}, errors.New("User does not match")
	}

	if !authenticateClientToken(req.Token) {
		return &pb_host.Status{Status: 1}, errors.New("Failed to authenticate Token")
	}

	closeClientConnection()

	return &pb_host.Status{Status: 0}, nil
}

func UpdateKey(ctx context.Context, req *pb_host.PublicKeyInfo) (*pb_host.Status, error) {
	if !authenticateClientToken(req.Token) {
		return &pb_host.Status{Status: 1}, errors.New("Failed to authenticate Token")
	}

	centralServer, err := getConnectionToCentralServer()
	if err != nil {
		return &pb_host.Status{Status: 1}, err
	}

	_, err = bytesToKey(req.Key)
	if err != nil {
		return &pb_host.Status{Status: 1}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	status, err := centralServer.server.UpdateKey(ctx, &pb_server.KeyUpdate{Username: settings.Username, AuthKey: req.Token, NewKey: req.Key})
	return &pb_host.Status{Status: status.Status}, err
}

func UserPing(ctx context.Context, req *pb_host.Username) (*pb_host.Status, error) {
	log.Println("Pinging...")
	if !authenticateClientToken(req.Token) {
		return &pb_host.Status{Status: 1}, errors.New("Failed to authenticate Token")
	}

	userConnection, err := getConnectionToUser(req.Username)
	if err != nil {
		return &pb_host.Status{Status: 1}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	status, _ := userConnection.user.Ping(ctx, &pb_host.Empty{})
	if status.Status == 0 {
		log.Println("Pong!")
	}
	return status, nil
}

func Ping(ctx context.Context, req *pb_host.Empty) *pb_host.Status {
	log.Println("Ping!")
	return &pb_host.Status{Status: 0}
}

func (h *host) DeleteMessage(ctx context.Context, req *pb_host.DeleteReq) (*pb_host.Status, error) {
	if authenticateClientToken(req.Token) {
		e := deleteMessageFromDB(req.User, req.MessageID)
		if e != nil {
			return &pb_host.Status{Status: 2}, e
		}
	}

	return &pb_host.Status{Status: 0}, nil
}

func (h *host) SendText(ctx context.Context, req *pb_host.ClientText) (*pb_host.Status, error) {
	if req.TargetUser != settings.Username && authenticateClientToken(req.Token) {
		userConnection, e := getConnectionToUser(req.TargetUser)
		if e != nil {
			log.Printf("Failed to connect to user: %v", e)
			return &pb_host.Status{Status: 1}, e
		}

		secret := generateSecret(userConnection.userInfo)
		msg := listofMessagesToString(req.Message)
		msgEncryptedForSending, er := encryptForSending(msg, userConnection.userInfo)
		if er != nil {
			log.Println(er)
			return &pb_host.Status{Status: 1}, nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		status, err := userConnection.user.RecieveText(ctx, &pb_host.H2HText{Message: stringToListofMessages(msgEncryptedForSending), User: settings.Username, Secret: secret})
		if err != nil {
			return status, err
		}

		msgEncryptedForClient, errr := encryptForClient(msg)
		if errr != nil {
			log.Printf("Message was not able to be Encrypted for storage!: %v", errr)
		} else {
			storeMesssage(req.TargetUser, true, msgEncryptedForClient)
		}

		return &pb_host.Status{Status: 0}, nil
	} else {
		return &pb_host.Status{Status: 1}, nil
	}
}

func (h *host) RecieveText(ctx context.Context, req *pb_host.H2HText) (*pb_host.Status, error) {
	recivedWithLocalInfo := false
	var status *pb_host.Status
	var userInfo *UserInfo
	var err error
	userInfo, err = getUserInfo(req.User)
	if err == nil {
		if verifySecret(userInfo, req.Secret) {
			message := listofMessagesToString(req.Message)
			storeMesssage(req.User, false, message)
			forawrdMessageToClient(req.Message)

			status = &pb_host.Status{Status: 0}
			recivedWithLocalInfo = true
		}
	}

	if !recivedWithLocalInfo {
		userInfo, err = getUserInfoFromSever(req.User)
		if err != nil {
			return &pb_host.Status{Status: 1}, err
		}

		if verifySecret(userInfo, req.Secret) {
			message := listofMessagesToString(req.Message)
			storeMesssage(req.User, false, message)
			forawrdMessageToClient(req.Message)

			status = &pb_host.Status{Status: 0}
		} else {
			return &pb_host.Status{Status: 1}, err
		}
	}

	_ = updateOrAddUserInfo(userInfo)

	return status, err
}

func (h *host) GetConversation(ctx context.Context, req *pb_host.Username) (*pb_host.Conversation, error) {
	var response *pb_host.Conversation = nil
	var err error = nil
	var messages []Message
	if authenticateClientToken(req.Token) {
		messages, err = getConversationFromDB(req.Username)

		if messages != nil {
			type Convo struct {
				Messages []Message `json:"messages"`
			}

			var convoJsonBytes []byte
			convoJsonBytes, err = json.Marshal(&Convo{Messages: messages})
			convoJson := string(convoJsonBytes)
			response = &pb_host.Conversation{Convo: stringToListofMessages(convoJson)}
		}
	}
	return response, err
}
