package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

type room struct {
	forwardChannel chan []byte
	joinChannel    chan *client
	leaveChannel   chan *client
	clients        map[*client]bool
}

func newRoom() *room {
	return &room{
		forwardChannel: make(chan []byte),
		joinChannel:    make(chan *client),
		leaveChannel:   make(chan *client),
		clients:        make(map[*client]bool),
	}

}

func (r *room) run() {
	for {
		select {
		case client := <-r.joinChannel:
			r.clients[client] = true
			fmt.Println("client join", client)
		case client := <-r.leaveChannel:
			delete(r.clients, client)
			close(client.sendChannel)
			fmt.Println("client leave", client)
		case msg := <-r.forwardChannel:
			fmt.Println("message", msg)
			for client := range r.clients {
				select {
				case client.sendChannel <- msg:
				default:
					delete(r.clients, client)
					close(client.sendChannel)
				}
			}
		}
	}

}

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// アクション
// WebSocket接続の更新
// クライアントのチャネルからの読み込み、書き込みスレッドの起動
func (room *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	fmt.Println("room.ServeHTTP start", req)

	// Websocket接続の更新
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// クライアントの生成
	client := &client{
		socket:      socket,
		sendChannel: make(chan []byte, messageBufferSize),
		room:        room,
	}

	// 参加チャネルにクライアントを追加
	// クライアントがルームに参加する。
	room.joinChannel <- client

	// クライアントからWebSocket接続が閉じられた時、
	// クライアントがルームから退席する。
	defer func() {
		room.leaveChannel <- client
	}()

	// 別スレッドで書き込み処理を起動
	go client.write()

	// 別スレッドで読み込み処理を起動
	go client.read()

	for {

	}
}
