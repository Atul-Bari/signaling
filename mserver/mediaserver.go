package main

import (
	"flag"
	"fmt"
	"net"
	"io"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	call "samespace.com/signaling/protobuf"

)

const compress = false

func main() {
	port := flag.Int("p", 57778, "port to listen to")
	flag.Parse()

	log.Infof("listening to port %d", *port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("could not listen to port %d: %v", *port, err)
	}

	s := grpc.NewServer()
	call.RegisterMediaprotoServer(s, server{})
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}

type server struct{}

func (server) BiCommonExchange(stream call.Mediaproto_BiCommonExchangeServer) error {
	for {
                log.Println("In server")
                in, err := stream.Recv()
                if err == io.EOF {
                        return nil
                }
                if err != nil {
                        return err
                }
                log.Info("received offer: ",in)
                if err := stream.Send(in); err != nil {
                                return err
                        }
        }
}
/*
func (server) InviteExchange(ctx context.Context, invite *call.Invite) (*call.Invite, error) {
	log.Info("received offer: ",invite)
	/*jsonMap := make(map[string]interface{})
	err := json.Unmarshal(text, &jsonMap)
	if err != nil {
		log.Println("Error while processing json:", err)
		break
	}

	eventType := fmt.Sprintf("%v", jsonMap["type"])
	log.Printf("Recv: %s",text)
	//log.Println("Type Received:", eventType, mt)
	switch eventType {
	case "echoTest":
		log.Println("Received Type Echo Test")
		echoTest(jsonMap, c)
	default:
		log.Printf("Unknown event is received")
	}

	return &call.Invite{Type: invite.Type, Host: "NO", Jsep: invite.Jsep, SessionID: invite.SessionID, Peer: invite.Peer, Record: false}, nil
	//create peerconnection from here

}*/
//bidirsdpexchange
/*func (server) BiInviteExchange(stream call.Makecall_BiInviteExchangeServer) error {
	for {
		log.Println("In server")
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Info("received offer: ",in)
		if err := stream.Send(in); err != nil {
				return err
			}
	}

}

func (server) BiCommonExchange(stream call.Makecall_BiCommonExchangeServer) error {
        for {
                log.Println("In server")
                in, err := stream.Recv()
                if err == io.EOF {
                        return nil
                }
                if err != nil {
                        return err
                }
                log.Info("received offer: ",in)
                if err := stream.Send(in); err != nil {
                                return err
                        }
        }
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	if compress {
		//b = zip(b)
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	if compress {
		//b = unzip(b)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}









/*package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/websocket"
)

// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Home Page")
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")

    upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    // upgrade this connection to a WebSocket
    // connection
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }
    // helpful log statement to show connections
    log.Println("Client Connected")

    err = ws.WriteMessage(1, []byte("Hi Client!"))
    if err != nil {
        log.Println(err)
    }
    reader(ws)

}

func reader(conn *websocket.Conn) {
    for {
    // read in a message
        messageType, p, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
    // print out that message for clarity
        fmt.Println(string(p))

        if err := conn.WriteMessage(messageType, p); err != nil {
            log.Println(err)
            return
        }

    }
}

func setupRoutes() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", wsEndpoint)
}

func main() {
    fmt.Println("Hello")
    setupRoutes()
    fmt.Println("Listening on 57778")
    log.Fatal(http.ListenAndServe("0.0.0.0:57778", nil))
}*/
