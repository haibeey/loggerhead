package loggerhead

import (
	"log"

	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

var socketServer = CreateChatServer()

type SocketHandler struct{}

func (sh *SocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socketServer.ServeHTTP(w, r)
}

// Event on websocket
var (
	StartChat = "startchat"
	CloseChat = "closechat"
	Message   = "message"
)

func CreateChatServer() *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		return nil
	})

	server.OnEvent("/", "message", func(s socketio.Conn, msg string) string {
		log.Println(server.BroadcastToNamespace("/", "message", msg))
		return ""
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
	})

	go server.Serve()
	return server
}
