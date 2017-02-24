package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type client struct {
	socket      *websocket.Conn
	sendChannel chan []byte
	room        *room
}

// ソケットからリクエストを読み込み、クライアントへの送信用チャンネルへ送信する。
func (c *client) read() {
	defer c.socket.Close()
	for {
		//if mType, msg, err := c.socket.ReadMessage(); err == nil {
		//	fmt.Println("message type:", mType)
		//	c.room.forwardChannel <- msg
		//} else {
		//	break
		//}

		messageType, message, err := c.socket.ReadMessage()
		if err != nil {
			log.Fatal("client.read", err)
			break
		}
		fmt.Println("messagetype:", messageType)
		fmt.Println("message:", message)
		c.room.forwardChannel <- message
	}
}

// ソケットからリクエストを読み込み、
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.sendChannel {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Fatal("write error:", err.Error())
			break
		}
	}
}
