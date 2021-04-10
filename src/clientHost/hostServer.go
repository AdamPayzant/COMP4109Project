package main

import (
	"context"
	"encoding/json"
	"log"
	"math"

	pb_host "pb_host"

	// pb_host "github.com/AdamPayzant/COMP4109Project/src/protos/smvshost"
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

func (h *host) ReKey(ctx context.Context, req *pb_host.Token) (*pb_host.Status, error) {
	/*
		This function is used to update the public key for the RSA encryption.
		The public key can only be changed if the correct auth key is provided to the main server
		The private key should only exist on TRUSTED end user clients.
	*/

	/*
		This is currently just in a state to demo gRPC call
		Plenty of stuff still to do
		TODO:
			Implement key management system
			Implement token management system
	*/
	// 	newkey, err := rsa.GenerateKey(rand.Reader, 2048)
	// 	if err != nil {
	// 		return &pb_host.Status{Status: 1}, errors.New("Key gen error")
	// 	}
	// 	// Get authtoken
	// 	authToken := make([]byte, 64)

	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// 	defer cancel()
	// 	keyBytes, _ := x509.MarshalPKIXPublicKey(&newkey.PublicKey)
	// 	r, err := server.UpdateKey(ctx, &pb_server.KeyUpdate{Username: username,
	// 		AuthKey: authToken,
	// 		NewKey:  keyBytes})
	// 	if err != nil {
	// 		return &pb_host.Status{Status: 2}, err
	// 	}
	// 	return &pb_host.Status{Status: r.Status}, nil
	// }

	// func (h *host) DeleteMessage(ctx context.Context, req *pb_host.DeleteReq) (*pb_host.Status, error) {
	// 	if authenticateClientToken(req.Token) {
	// 		e := deleteMessageFromDB(req.User, req.MessageID)
	// 		if e != nil {
	// 			return &pb_host.Status{Status: 2}, e
	// 		}
	// 	}
	return &pb_host.Status{Status: 0}, nil
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
		clientConnection, e := getConnectionToUser(req.TargetUser)
		if e != nil {
			log.Printf("Failed to connect to user: %v", e)
			return &pb_host.Status{Status: 1}, e
		}

		secret := generateSecret(clientConnection.userInfo)
		msg := listofMessagesToString(req.Message)
		msgEncryptedForSending, er := encryptForSending(msg, clientConnection.userInfo)
		if er != nil {
			log.Println(er)
			return &pb_host.Status{Status: 1}, nil
		}

		startus, err := clientConnection.client.RecieveText(context.Background(), &pb_host.H2HText{Message: stringToListofMessages(msgEncryptedForSending), User: settings.Username, Secret: secret})
		if err != nil {
			return startus, err
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
