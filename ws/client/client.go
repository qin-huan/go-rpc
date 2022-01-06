package main

import (
	"bufio"
	_ "bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mPb "github.com/echoingtech/messager/proto/go.micro.srv.messager"
	pb "github.com/echoingtech/pgw/proto/pgw"
	ucPb "github.com/echoingtech/uc/proto/go.micro.srv.uc"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	_ "os"
	"strconv"
	"strings"
)

//import (
//	"encoding/base64"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/websocket"
//	"time"
//)
//
//type Message struct {
//	Srv  *Srv   `json:"srv"`
//	User *User `json:"user"`
//	Data string `json:"data"`
//}
//type Srv struct {
//	Srv  string `json:"srv"`
//	Pkg  string `json:"pkg"`
//	Md   string `json:"md"`
//}
//
//type LoginRequest struct {
//	Token string `json:"token"`
//	User  *User  `json:"user"`
//	Plat  string `json:"plat"`
//}
//
//type User struct {
//	UserId  string `json:"user_id"`
//	Channel string `json:"channel"`
//	Device  string `json:"device"`
//}
//
//func main() {
//	url := "ws://dev-gapi.echoing.tech/pgw/ws"
//
//	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
//	if err != nil {
//		panic(err)
//	}
//	defer ws.Close()
//
//	// read message
//	go func() {
//		for {
//			_, msg, err := ws.ReadMessage()
//			if err != nil {
//				fmt.Println("err: ", err)
//				continue
//			}
//			fmt.Println("msg: ", string(msg))
//		}
//	}()
//
//	// login
//	loginRequest := &LoginRequest{
//		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjM2NjU2Njg1NTkwMjUwNDEwOSIsInR5cGUiOiJVU0VSIn0.qakbo1JjCoOYTqd3khw1UyTh6vniY9MCNXYG2-OnvkM",
//		User: &User{
//			UserId: "366566855902504109",
//		},
//	}
//	bytes, _ := json.Marshal(loginRequest)
//	request := base64.StdEncoding.EncodeToString(bytes)
//
//	login := &Message{
//		Srv:  &Srv{
//			Srv: "go.micro.srv.pgw",
//			Pkg: "PushGateWay",
//			Md:  "Login",
//		},
//		Data: request,
//	}
//	loginBytes, _ := json.Marshal(login)
//	if err = ws.WriteMessage(websocket.TextMessage, loginBytes); err != nil {
//		panic(err)
//	}
//
//	// send message
//	var str = "hello, 366566855902504109"
//	data := base64.StdEncoding.EncodeToString([]byte(str))
//	msg := &Message{
//		Srv: &Srv{
//			Srv: "foo",
//			Pkg: "pkg",
//			Md: "md",
//		},
//		User: &User{
//			UserId: "366566855902504109",
//		},
//		Data: data,
//	}
//	marshal, _ := json.Marshal(msg)
//
//	if err = ws.WriteMessage(websocket.TextMessage, marshal); err != nil {
//		panic(err)
//	}
//
//	time.Sleep(time.Second*10)
//
//	//reader := bufio.NewReader(os.Stdin)
//	//for {
//	//	str, _ := reader.ReadString('\n')
//	//	fmt.Println(str)
//	//}
//}

func main() {
	// dev
	//url := "ws://dev-gapi.echoing.tech/pgw/ws"

	// local
	url := "ws://127.0.0.1:8080/pgw/ws"

	header := http.Header{}
	header.Add("Content-Type", "application/json")
	ws, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	// read message
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				ws, _, err = websocket.DefaultDialer.Dial(url, nil)
				if err != nil {
					log.Println(err)
				}
				continue
			}
			log.Println("msg: ", string(msg))
		}
	}()

	//log.Println("input senderId...")
	reader := bufio.NewReader(os.Stdin)
	//sender, err := reader.ReadString('\n')
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//senderId, err := strconv.ParseInt(strings.TrimSpace(sender), 10, 64)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println("input receiverId...")
	//recv, err := reader.ReadString('\n')
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//receiverId, err := strconv.ParseInt(strings.TrimSpace(recv), 10, 64)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	var senderId int64 = 366566855902504109
	var receiverId int64 = 86336294443024385

	if err = pushLogin(ws, strconv.FormatInt(senderId, 10)); err != nil {
		log.Fatalln(err)
	}

	for {
		readString, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}

		text := &mPb.TextMessage{
			Text: "hello " + strings.Trim(readString, "\r\n"),
		}
		raw, _ := json.Marshal(text)

		msg := &mPb.SendMessageRequest{
			SenderId:        senderId,
			ChatId:          generateChatId(senderId, receiverId),
			ClientRequestId: generateClientRequestId(),
			Content:         string(raw),
			ContentType:     mPb.Message_MESSAGE_CONTENT_TYPE_TEXT,
		}

		log.Println(msg)
		if err = pushSendMessage(ws, msg); err != nil {
			log.Println(err)
		}

		//if err = pushMessage(ws, strings.TrimSpace(sender), strings.TrimSpace(recv), fmt.Sprintf("hello, %s", readString)); err != nil {
		//	log.Println(err)
		//}

		//if err = pushGetUser(ws, senderId); err != nil {
		//	log.Println(err)
		//}
	}
}

func pushLogin(ws *websocket.Conn, userId string) error {
	pgwSrv := "pgw"

	loginMsg := pb.Message{
		Srv: &pb.Service{
			Srv: pgwSrv,
			Pkg: "PushGateWay",
			Md:  "Login",
		},
		Mtype: pb.Message_MICRO_GRPC,
		Meta: &pb.Message_Meta{
			Id:  userId,
			Ack: pb.Message_Meta_Ack(pb.Message_Meta_Ack_value["UACK"]),
		},
	}
	data, err := json.Marshal(&pb.LoginRequest{
		User: &pb.User{
			UserId:  userId,
			Device:  "lenovo",
			Channel: "test",
		},
	})
	if err != nil {
		return err
	}

	loginMsg.DataS = string(data)
	bytes, err := json.Marshal(&loginMsg)
	if err != nil {
		return err
	}
	if err = ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
		return err
	}

	return nil
}

func pushGetUser(ws *websocket.Conn, userId int64) error {
	ucSrv := "go.micro.srv.uc"

	ucMsg := pb.Message{
		Srv: &pb.Service{
			Srv: ucSrv,
			Pkg: "PushGateWay",
			Md: "GetUser",
		},
		Mtype: pb.Message_MICRO_GRPC,
		Meta: &pb.Message_Meta{
			Id:  strconv.FormatInt(userId, 10),
			Ack: pb.Message_Meta_Ack(pb.Message_Meta_Ack_value["RECV"]),
		},
	}
	bytes, err := json.Marshal(&ucPb.GetUserRequest{
		UserId:       userId,
		DisplayFlags: nil,
	})
	if err != nil {

	}

	ucMsg.Data = bytes
	raw, err := json.Marshal(&ucMsg)
	if err != nil {
		return err
	}

	if err = ws.WriteMessage(websocket.TextMessage, raw); err != nil {
		return err
	}

	return nil
}

func pushMessage(ws *websocket.Conn, userId, receiverId string, msg string) error {
	pgwSrv := "go.micro.srv.pgw"

	pushMsg := pb.Message{
		Srv: &pb.Service{
			Srv: pgwSrv,
			Pkg: "PushGateWay",
			Md:  "Push",
		},
		Mtype: pb.Message_MICRO_GRPC,
		Meta: &pb.Message_Meta{
			Id:  userId,
			Ack: pb.Message_Meta_Ack(pb.Message_Meta_Ack_value["RECV"]),
		},
	}
	data, err := json.Marshal(&pb.PushRequest{
		UserId: receiverId,
		Data:   []byte(msg),
	})
	if err != nil {
		return err
	}

	pushMsg.Data = data
	bytes, err := json.Marshal(&pushMsg)
	if err != nil {
		return err
	}
	if err = ws.WriteMessage(websocket.TextMessage, bytes); err != nil {
		return err
	}

	return nil
}

func pushSendMessage(ws *websocket.Conn, req *mPb.SendMessageRequest) error {
	imSrv := "messager"

	msg := pb.Message{
		Srv: &pb.Service{
			Srv: imSrv,
			Pkg: "Messager",
			Md:  "SendMessage",
		},
		Mtype: pb.Message_MICRO_GRPC,
		Meta: &pb.Message_Meta{
			Id:  strconv.FormatInt(req.SenderId, 10),
			Ack: pb.Message_Meta_Ack(pb.Message_Meta_Ack_value["RECV"]),
		},
	}

	bytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	//marshal, err := proto.Marshal(&pb.Message{
	//	Data: bytes,
	//})

	msg.Data = bytes
	raw, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	if err = ws.WriteMessage(websocket.TextMessage, raw); err != nil {
		return err
	}

	return nil
}

var p2pPrefix = "v1:p2p:"

func generateChatId(foo, bar int64) string {
	if foo > bar {
		foo, bar = bar, foo
	}
	return base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s%d:%d", p2pPrefix, foo, bar)))
}

func generateClientRequestId() string {
	key := make([]byte, 8)
	_, err := rand.Read(key)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(key)
}
