package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 2048
)

type echotestReq struct {
	Type      string `json:"type"`
	SessionID string `json:"sessionId"`
	PeerId    string `json:"peerId"`
	Jsep      string `json:"jsep"`
}

type icecandidatereq struct {
	Type      string `json:"type"`
	SessionId string `json:"sessionId"`
	PeerId    string `json:"peerId"`
	Jsep      string `json:"candidate"`
}

type offermeetingreq struct {
	Type      string `json:"type"`
	SessionId string `json:"sessionId"`
	HostId    string `json:"hostid"`
	PeerId    string `json:"peerId"`
	Jsep      string `json:"jsep"`
}

type ansmeetingreq struct {
	Type      string `json:"type"`
	SessionId string `json:"sessionId"`
	HostId    string `json:"hostid"`
	PeerId    string `json:"peerId"`
	Jsep      string `json:"jsep"`
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	peerid string

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

//unmarhsel received from browser
func CreateEchoFromMap(m map[string]interface{}) (echotestReq, error) {
	data, _ := json.Marshal(m)
	var result echotestReq
	err := json.Unmarshal(data, &result)
	return result, err
}

func CreateOfferFromMap(m map[string]interface{}) (offermeetingreq, error) {
	data, _ := json.Marshal(m)
	var result offermeetingreq
	err := json.Unmarshal(data, &result)
	return result, err
}

func CreateAnsFromMap(m map[string]interface{}) (ansmeetingreq, error) {
	data, _ := json.Marshal(m)
	var result ansmeetingreq
	err := json.Unmarshal(data, &result)
	return result, err
}

func CreateIceFromMap(m map[string]interface{}) (icecandidatereq, error) {
	data, _ := json.Marshal(m)
	var result icecandidatereq
	err := json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("Raw Received:", string(message))

		jsonMap := make(map[string]interface{})
		err = json.Unmarshal(message, &jsonMap)

		eventType := fmt.Sprintf("%v", jsonMap["type"])

		//echoreq, _ := Unmarshaloffer(message)
		log.Println()

		//c.hub.broadcast <- &data
		switch eventType {
		case "echoTest":
			log.Println("Received Type Echo Test", eventType)
			requestStruct, err := CreateEchoFromMap(jsonMap)
			if err != nil {
				log.Println("No marshel")
			}
			data := echoMsgDetails{requestStruct, c}
			// echodata <- &data
			go echoMakecall(data)

		case "offerMeeting":
			log.Println("Received Type Offer :", eventType)
			requestStruct, err := CreateOfferFromMap(jsonMap)
			if err != nil {
				log.Println("No marshel")
			}

			data := offerMsgDetails{requestStruct, c}
			// offerdata <- &data
			go offerMakecall(data)
		case "iceCandidate":
			log.Println("Received Type Ice :", eventType)
			requestStruct, err := CreateIceFromMap(jsonMap)
			if err != nil {
				log.Println("No marshel")
			}

			data := iceMsgDetails{requestStruct, c}
			// icedata <- &data
			go iceMakecall(data)
		case "answerMeeting":
			log.Println("Received Type Answer :", eventType)
			requestStruct, err := CreateAnsFromMap(jsonMap)
			if err != nil {
				log.Println("No marshel")
			}

			data := ansMsgDetails{requestStruct, c}
			// ansdata <- &data
			go ansMakecall(data)

		}

	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var count int64

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	n := atomic.AddInt64(&count, 1)
	log.Printf("Total number of connections: %v", n)
	Clientid := fmt.Sprintf("client%d", n)
	log.Printf("Total number of connections: %v", Clientid)
	client := &Client{peerid: Clientid, hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
