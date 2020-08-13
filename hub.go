// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	pb "server/proto"

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

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *MsgDetails

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan *MsgDetails),
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
				if client == message.c {
					continue
				}
				select {
				case client.send <- []byte(message.jsep):
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func Makecall(msg echoMsgDetails) {

	log.Println("In echotest Makecall()")
	conn, err := grpc.Dial("203.153.53.181:57778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v", err)
	}
	defer conn.Close()

	client := pb.NewMakecallClient(conn)
	// log.Println("Type: ", msg.msg.Type)
	// log.Println("Jsep", msg.msg.Jsep)
	// log.Println("Session", msg.msg.SessionID)
	// log.Println("Peer", msg.msg.PeerId)

	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
	//	text := &pb.Sdp{Type: msg.msg.Type, Hostid: "NO", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("Makecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	//msg.msg = []byte(resp.GetOffer())
	tmp := MsgDetails{
		jsep: resp.Jsep,
		c:    msg.cli,
	}
	msg.cli.hub.broadcast <- &tmp

}

func offerMakecall(msg offerMsgDetails) {
	log.Println("In offer Makecall()")
	conn, err := grpc.Dial("203.153.53.181:57778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v", err)
	}
	defer conn.Close()

	client := pb.NewMakecallClient(conn)
	// log.Println("Type: ", msg.msg.Type)
	// log.Println("Jsep", msg.msg.Jsep)
	// log.Println("Session", msg.msg.SessionId)
	// log.Println("Peer", msg.msg.PeerId)

	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
	//	text := &pb.Sdp{Type: msg.msg.Type, Hostid: "NO", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("offerMakecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	//msg.msg = []byte(resp.GetOffer())
	tmp := MsgDetails{
		jsep: resp.Jsep,
		c:    msg.cli,
	}
	msg.cli.hub.broadcast <- &tmp

}

func ansMakecall(msg ansMsgDetails) {

	log.Println("In answer Makecall()")
	conn, err := grpc.Dial("203.153.53.181:57778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v", err)
	}
	defer conn.Close()

	client := pb.NewMakecallClient(conn)
	// log.Println("Type: ", msg.msg.Type)
	// log.Println("Jsep", msg.msg.Jsep)
	// log.Println("Session", msg.msg.SessionId)
	// log.Println("Peer", msg.msg.PeerId)

	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
	//	text := &pb.Sdp{Type: msg.msg.Type, Hostid: "NO", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("ansMakecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	//msg.msg = []byte(resp.GetOffer())
	tmp := MsgDetails{
		jsep: resp.Jsep,
		c:    msg.cli,
	}
	msg.cli.hub.broadcast <- &tmp

}

func iceMakecall(msg iceMsgDetails) {

	log.Println("In Ice Makecall()")
	conn, err := grpc.Dial("203.153.53.181:57778", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to: %v", err)
	}
	defer conn.Close()

	client := pb.NewMakecallClient(conn)
	// log.Println("Type: ", msg.msg.Type)
	// log.Println("Jsep", msg.msg.Jsep)
	// log.Println("Session", msg.msg.SessionId)
	// log.Println("Peer", msg.msg.PeerId)

	test := &pb.Sdp{Type: msg.msg.Type, Hostid: "No", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionId, Peerid: msg.msg.PeerId}
	//	text := &pb.Sdp{Type: msg.msg.Type, Hostid: "NO", Jsep: msg.msg.Jsep, Sessionid: msg.msg.SessionID, Peerid: msg.msg.PeerId}
	resp, err := client.Sdpexchange(context.Background(), test)
	log.Println("iceMakecall resp Jsep", resp.Jsep)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Jsep, err)
	}
	//msg.msg = []byte(resp.GetOffer())
	tmp := MsgDetails{
		jsep: resp.Jsep,
		c:    msg.cli,
	}
	msg.cli.hub.broadcast <- &tmp

}
