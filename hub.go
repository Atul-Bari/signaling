// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"google.golang.org/grpc"
        "golang.org/x/net/context"
	pb "server/proto"
	"flag"

)
// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
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
			//Makecall(message)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func Makecall(msg  []byte) {
	backend := flag.String("b", "localhost:57778", "address of the say backend")
        flag.Parse()

        conn, err := grpc.Dial(*backend, grpc.WithInsecure())
        if err != nil {
                log.Fatalf("could not connect to %s: %v", *backend, err)
        }
        defer conn.Close()

	client := pb.NewMakecallClient(conn)

        text := &pb.Sdp{Fromid: "X", Toid:"Y", Offer:string(msg), Sessionid:"qwerty"}

	resp, err := client.Sdpexchange(context.Background(), text)
	if err != nil {
		log.Fatalf("could not say %s: %v", resp.Offer, err)
        }


}
