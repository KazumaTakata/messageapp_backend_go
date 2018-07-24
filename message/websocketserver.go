package message

import (
	"encoding/json"
	"fmt"
	"http_server/db"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var session = db.DBsession()

var connMap = make(map[string]*websocket.Conn)

func Websocketmessage(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	// userid := middleware.RequestIDFromContext(r.Context())
	// connMap[userid] = conn

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		message := new(Messaege)
		if err := json.Unmarshal(msg, message); err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return
		}

		tokenData := Token{}
		_, _ = jwt.ParseWithClaims(message.Myid, &tokenData, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if message.Ping == "hey" {
			connMap[tokenData.ID] = conn
		} else {

			returndata := ReturnMessage{ID: tokenData.ID, Content: message.Content}
			messagejson, err := json.Marshal(returndata)
			// Write message back to browser
			if friendconn, ok := connMap[message.Friendid]; ok {
				if err = friendconn.WriteMessage(msgType, messagejson); err != nil {
					return
				}
			} else {
				session.InsertTalk(message.Friendid, message.Content, tokenData.ID)
			}

		}

	}
}

type Messaege struct {
	Myid     string `json:"myId"`
	Ping     string
	Content  string
	Friendid string `json:"friendId"`
}

type ReturnMessage struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Token struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
