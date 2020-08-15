// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"io"
	"log"
	pb "server/proto"
	"time"

	//"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type echoMsgDetails struct {

	//msg *offerMeetingReq
	//msg []byte
	msg echotestReq
	cli *Client
}

type offerMsgDetails struct {

	//msg *offerMeetingReq
	//msg []byte
	msg offermeetingreq
	cli *Client
}

type ansMsgDetails struct {

	//msg *offerMeetingReq
	//msg []byte
	msg ansmeetingreq
	cli *Client
}

type iceMsgDetails struct {

	//msg *offerMeetingReq
	//msg []byte
	msg icecandidatereq
	cli *Client
}

type MsgDetails struct {
	jsep string
	c    *Client
}

var echodata chan *echoMsgDetails
var offerdata chan *offerMsgDetails
var ansdata chan *ansMsgDetails
var icedata chan *iceMsgDetails

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *pb.Sdp

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *pb.Sdp),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Registered!!!")
			//Makecall(message)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			//message.msg = Makecall(message.msg)

			for client := range h.clients {
				if client.peerid == message.Peerid {
					continue
				}
				select {
				case client.send <- []byte(message.Jsep):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

var client pb.MakecallClient
var grpcconn *grpc.ClientConn

func grpcConnection() {
	var err error
	grpcconn, err = grpc.Dial("203.153.53.181:57778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v", err)
	}
	client = pb.NewMakecallClient(grpcconn)
}

func echoMakecall(msg echoMsgDetails) {

	log.Println("In echotest Makecall()")
	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("Makecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	// tmp := MsgDetails{
	// 	jsep: resp.Jsep,
	// 	c:    msg.cli,
	// }
	msg.cli.hub.broadcast <- resp

}

func offerMakecall(msg offerMsgDetails) {
	log.Println("In offer Makecall()")
	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("offerMakecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	// tmp := MsgDetails{
	// 	jsep: resp.Jsep,
	// 	c:    msg.cli,
	// }
	msg.cli.hub.broadcast <- resp

}

func ansMakecall(msg ansMsgDetails) {

	log.Println("In answer Makecall()")
	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("ansMakecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	// tmp := MsgDetails{
	// 	jsep: resp.Jsep,
	// 	c:    msg.cli,
	// }
	// msg.cli.hub.broadcast <- &tmp
	msg.cli.hub.broadcast <- resp

}

func iceMakecall(msg iceMsgDetails) {

	log.Println("In Ice Makecall()")
	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("iceMakecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	// tmp := MsgDetails{
	// 	jsep: resp.Jsep,
	// 	c:    msg.cli,
	// }
	msg.cli.hub.broadcast <- resp

}

//

//
//
//
//
//
//
//
//
//
//
//

//used in case bidir grpc needed
func BidirCall() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	stream, err := client.Bidirsdpexchange(ctx)
	defer cancel()
	if err != nil {
		log.Fatalf("Failed to receive a note : %v", err)
	}
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				// read done.
				//close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v%v", resp, err)
			}

			log.Println("Bidir Resp:", resp)

			// tmp := MsgDetails{
			// 	jsep: resp.Jsep,
			// 	c:    msg.cli,
			// }
			hub.broadcast <- resp
			// log.Printf("Got message %s at point(%d, %d)", )
		}
	}()
	for {
		select {
		case msg := <-echodata:
			test := &pb.Sdp{Type: msg.msg.Type, Hostid: msg.msg.PeerId, Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
			if err := stream.Send(test); err != nil {
				log.Fatalf("Failed to send a sdp: %v", err)
			}
		case msg := <-offerdata:
			test := &pb.Sdp{Type: msg.msg.Type, Hostid: msg.msg.PeerId, Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
			if err := stream.Send(test); err != nil {
				log.Fatalf("Failed to send a sdp: %v", err)
			}

		case msg := <-ansdata:
			test := &pb.Sdp{Type: msg.msg.Type, Hostid: msg.msg.PeerId, Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
			if err := stream.Send(test); err != nil {
				log.Fatalf("Failed to send a sdp: %v", err)
			}

		case msg := <-icedata:
			test := &pb.Sdp{Type: msg.msg.Type, Hostid: msg.msg.PeerId, Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
			if err := stream.Send(test); err != nil {
				log.Fatalf("Failed to send a sdp: %v", err)
			}
		}
	}
}
