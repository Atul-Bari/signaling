// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"log"
	pb "server/proto"

	//"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type MsgDetails struct {
	msg []byte
	cli *Client
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
				if client == message.cli {
					continue
				}
				select {
				case client.send <- message.msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func Makecall(msg *MsgDetails) {
	backend := flag.String("b", "localhost:57778", "address of the say backend")
	flag.Parse()
	log.Println()
	log.Println()
	log.Println()
	log.Println("################# Make call called")
	log.Println()
	log.Println()

	conn, err := grpc.Dial(*backend, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to %s: %v", *backend, err)
	}
	defer conn.Close()

	client := pb.NewMakecallClient(conn)

	text := &pb.Sdp{Fromid: "X", Toid: "Y", Offer: string(msg.msg), Sessionid: "qwerty"}
	log.Println()
	log.Println()
	log.Println()
	log.Println("################# 2Make call called")
	log.Println()
	log.Println()
	resp, err := client.Sdpexchange(context.Background(), text)
	log.Println()
	log.Println()
	log.Println()
	log.Println("################# 3Make call called")
	log.Println()
	log.Println()
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Offer, err)
	}
	//msg.msg = []byte(resp.GetOffer())
	log.Println()
	log.Println()
	log.Println()
	log.Println("#################4 Make call called")
	log.Println()
	log.Println()
	msg.cli.hub.broadcast <- msg

}
