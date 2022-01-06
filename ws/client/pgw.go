package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	grpcclient "github.com/echoingtech/go-micro/client/grpc"
	cfgcmd "github.com/echoingtech/go-micro/config/cmd"
	"github.com/golang/protobuf/proto" // nolint
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/echoingtech/pgw/proto/pgw"
)

func main() {
	cmd := cobra.Command{
		Use:   "proxy",
		Short: "Chat Example",
	}

	addr := "ws://dev-gapi.echoing.tech/pgw/ws"
	cmd.PersistentFlags().StringVar(&addr, "addr", addr, "Addr of websocket.")

	notUseJSONBytes := false
	cmd.PersistentFlags().BoolVar(
		&notUseJSONBytes, "not-use-json-bytes", notUseJSONBytes,
		"Don't use bytes in json, avoid golang auto base64.")

	prefix := "go.micro.srv"
	cmd.PersistentFlags().StringVar(&prefix, "prefix", prefix, "Prefix of service.")
	suffix := "local"
	cmd.PersistentFlags().StringVar(&suffix, "suffix", suffix, "Suffix of service.")

	ack := "UNACK"
	cmd.PersistentFlags().StringVar(&ack, "ack", ack, "Level of ack: UNACK,RECV,COMPLETE")

	shouldLogin := true
	cmd.PersistentFlags().BoolVar(&shouldLogin, "login", shouldLogin, "Login before other request.")

	name := "pgw"
	cmd.PersistentFlags().StringVar(&name, "name", name, "Service name.")

	encoding := "json"
	cmd.PersistentFlags().StringVar(&encoding, "encoding", encoding, "Encoding: proto, json.")

	var conn *websocket.Conn
	login := func(*cobra.Command, []string) error {
		msg := pb.Message{
			Srv: &pb.Service{
				Srv: name,
				Pkg: "PushGateWay",
				Md:  "Login",
			},
			Mtype: pb.Message_MICRO_GRPC,
		}
		if ack != "UNACK" {
			id := uuid.New().String()
			msg.Meta = &pb.Message_Meta{
				Id:  id,
				Ack: pb.Message_Meta_Ack(pb.Message_Meta_Ack_value[ack]),
			}
		}
		var data []byte
		var err error
		if encoding == "json" {
			data, err = json.Marshal(&pb.LoginRequest{
				User: &pb.User{
					UserId:  "j2gg0s",
					Device:  "j2gg0s's MacbookPro",
					Channel: "example",
				},
			})
		} else {
			data, err = proto.Marshal(&pb.LoginRequest{
				User: &pb.User{
					UserId:  "j2gg0s",
					Device:  "j2gg0s's MacbookPro",
					Channel: "example",
				},
			})
		}
		if err != nil {
			return fmt.Errorf("proto marshal: %w", err)
		}
		if notUseJSONBytes {
			msg.DataS = string(data)
		} else {
			msg.Data = data
		}

		raw, err := json.Marshal(&msg)
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}
		if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
			return fmt.Errorf("write to ws: %w", err)
		}
		fmt.Println("Login:", string(raw))

		return nil
	}

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		conn, _, err = websocket.DefaultDialer.Dial(
			fmt.Sprintf("%s?not_use_json_bytes=%t", addr, notUseJSONBytes), nil)
		if err != nil {
			return fmt.Errorf("dial %s: %w", addr, err)
		}

		go func() {
			for {
				conn.SetCloseHandler(nil)
				conn.SetCloseHandler(func(code int, text string) error {
					fmt.Println("Websocket recv peer CloseMessage")
					return conn.CloseHandler()(code, text)
				})
				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Websocket read: ", err)
					return
				}
				fmt.Println("Websocket recv msg: ", string(msg))
			}
		}()

		if shouldLogin {
			return login(cmd, args)
		}
		return nil
	}
	cmd.PersistentPostRunE = func(*cobra.Command, []string) error {
		time.Sleep(5 * time.Second)
		if err := conn.Close(); err != nil {
			return fmt.Errorf("close websocket: %w", err)
		}
		return nil
	}

	echo := func(*cobra.Command, []string) error {
		msg := pb.Message{
			Srv: &pb.Service{
				Srv: name,
				Pkg: "PushGateWay",
				Md:  "Echo",
			},
			Mtype: pb.Message_MICRO_GRPC,
		}
		if ack != "UNACK" {
			id := uuid.New().String()
			msg.Meta = &pb.Message_Meta{
				Id:  id,
				Ack: pb.Message_Meta_Ack(pb.Message_Meta_Ack_value[ack]),
			}
		}
		data, err := proto.Marshal(&pb.Message{
			Data: []byte("Hello World"),
		})
		if err != nil {
			return fmt.Errorf("proto marshal: %w", err)
		}
		if notUseJSONBytes {
			msg.DataS = string(data)
		} else {
			msg.Data = data
		}

		raw, err := json.Marshal(&msg)
		if err != nil {
			return fmt.Errorf("json marshal: %w", err)
		}

		if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
			return fmt.Errorf("write to ws: %w", err)
		}
		fmt.Println("Route:", string(raw))
		return nil
	}

	push := func(*cobra.Command, []string) error {
		msg := pb.PushRequest{
			UserId: "j2gg0s",
			Srv: &pb.Service{
				Srv: "foo",
				Pkg: "Foo",
				Md:  "Notify",
			},
			Data: []byte(fmt.Sprintf(
				"Service %s Notify user(j2gg0s) current is %s",
				"foo",
				time.Now().Format(time.RFC3339))),
		}

		cli := grpcclient.NewClient()
		opts := cli.Options()
		err := cfgcmd.Init(
			cfgcmd.Registry(&opts.Registry),
			cfgcmd.Selector(&opts.Selector),
		)
		if err != nil {
			return fmt.Errorf("init config/cmd %w", err)
		}
		err = cli.Init()
		if err != nil {
			return fmt.Errorf("init client %w", err)
		}

		srv := fmt.Sprintf("%s.%s.%s", prefix, name, suffix)
		req := cli.NewRequest(srv, "PushGateWay.Push", &msg)
		resp := emptypb.Empty{}
		if err := cli.Call(context.Background(), req, &resp); err != nil {
			return fmt.Errorf("Push %w", err)
		}

		fmt.Println("Push:", msg.String())

		return nil
	}

	cmd.AddCommand(
		&cobra.Command{Use: "echo", RunE: echo},
		&cobra.Command{Use: "push", RunE: push},
		&cobra.Command{Use: "all", RunE: func(c *cobra.Command, args []string) error {
			if err := echo(c, args); err != nil {
				return err
			}
			if err := push(c, args); err != nil {
				return err
			}
			return nil
		}},
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}