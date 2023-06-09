package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"online_chess/modules/player"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	send_wait = 10 * time.Second

	// Max time till next pong from peer
	pong_wait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	ping_period = (pong_wait * 9) / 10

	// Maximum message size allowed from peer.
	max_message_size = 10000
)

type Client struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	room     *Room
}

func newClient(conn *websocket.Conn, wsServer *WsServer, id int, name string) *Client {
	return &Client{
		ID:       id,
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
	}
}

func (client *Client) disconnect() {
	close(client.send)
	client.conn.Close()
}

func (client *Client) receivePump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(max_message_size)             //sets the maximum size in bytes for a message read from the peer.
	client.conn.SetReadDeadline(time.Now().Add(pong_wait)) //set the read deadline timestamp on the connection
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pong_wait)); return nil })

	for {
		_, json_message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(json_message)
	}
}

var new_line = []byte{'\n'}

func (client *Client) sendPump() {
	ticker := time.NewTicker(ping_period)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(send_wait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			for i := 0; i < len(client.send); i++ {
				w.Write(new_line)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(send_wait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type token_query struct {
	Token string `form:"token" binding:"required"`
}

func RunClient(c *gin.Context) {
	var query token_query
	err := c.ShouldBindQuery(&query)
	if err != nil {
		log.Println("get token query failed", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := player.QueryPlayerByToken(query.Token)
	if err != nil {
		log.Println("QueryPlayerByToken failed", err)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrader.Upgrade failed", err)
		return
	}
	wsServer := c.MustGet("wsServer").(*WsServer)
	client := newClient(conn, wsServer, player.Id, player.Name)
	client.notify(connect_websocket_success, nil)
	wsServer.addClient(client)

	go client.receivePump()
	go client.sendPump()
}

func (client *Client) handleNewMessage(json_message []byte) {
	var message Message
	if err := json.Unmarshal(json_message, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	switch message.Action {
	case join_room_action:
		client.handleJoinRoomMessage(message)
	case create_room_action:
		client.createRoom(message)
	case get_rooms_action:
		client.handleGetRoomListMessage()
	}
}

func (client *Client) createRoom(message Message) {
	room_name, ok := message.Payload.(string)
	if !ok {
		client.notify(create_room_failed, "Please make sure room name is valid")
		return
	}
	room := client.wsServer.findRoomByName(room_name)
	if room != nil {
		client.notify(create_room_failed, "Room name is already taken")
		return
	}
	client.wsServer.createRoom(room_name, client)
	client.notify(create_room_success, *(client.room))
}

func (client *Client) handleJoinRoomMessage(message Message) {
	room_id_string, ok := message.Payload.(string)
	room_id, err := uuid.Parse(room_id_string)
	if !ok || err != nil {
		client.notify(join_room_failed, "Failed to join room because of the room is not exist")
		return
	}
	client.joinRoom(room_id)
}

func (client *Client) IsInRoom(room *Room) (bool, error) {
	if client.room == nil {
		return false, nil
	} else if client.room != room {
		return false, fmt.Errorf("You have joined other room")
	} else {
		return true, fmt.Errorf("You are already in this room")
	}
}

func (client *Client) ValidRoomCanJoin(room *Room) error {
	var valid_room_err error
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := client.IsInRoom(room)
		if err != nil {
			valid_room_err = err
		}
	}()

	go func() {
		defer wg.Done()
		has_vancancy := room.hasVancancy()
		if !has_vancancy {
			valid_room_err = fmt.Errorf("The room you want to attend is no longer available")
		}
	}()
	wg.Wait()
	return valid_room_err
}

func (client *Client) joinRoom(room_id uuid.UUID) {
	var valid_room_err error
	is_room_exist, room := client.wsServer.findRoomByID(room_id)
	if !is_room_exist {
		valid_room_err = fmt.Errorf("The room you are trying to join does not exist")
	} else {
		valid_room_err = client.ValidRoomCanJoin(room)
	}
	if valid_room_err != nil {
		client.notify(join_room_failed, valid_room_err.Error())
		return
	}

	join_err := room.addClientInRoom(client)
	if join_err != nil {
		client.notify(join_room_failed, join_err.Error())
	} else {
		client.room = room
		success_msg := &Message{
			Action:  join_room_success,
			Payload: client.Name + " success to join to this room",
		}
		room.broadcast <- success_msg
	}
}

func (client *Client) handleGetRoomListMessage() {
	rooms := client.wsServer.GetRoomList()
	client.notify(get_rooms_success, rooms)
}

func (client *Client) notify(action string, payload interface{}) {
	message := &Message{
		Action:  action,
		Payload: payload,
	}
	client.send <- message.encode()
}
