package game

import (
	"fmt"

	"github.com/google/uuid"
)

type Room struct {
	ID                 uuid.UUID
	Name               string
	max_clients_number int
	admin_id           int
	clients            map[int]*Client
	applicants         chan *Client
	broadcast          chan *Message
}

func (room *Room) Run() {
	for {
		select {
		case message := <-room.broadcast:
			room.broadcastToClients(message)
		}
	}
}

func NewRoom(name string, admin *Client) *Room {
	new_room := &Room{
		ID:                 uuid.New(),
		Name:               name,
		max_clients_number: 2,
		admin_id:           admin.ID,
		clients:            make(map[int]*Client),
		applicants:         make(chan *Client, 2),
		broadcast:          make(chan *Message, 10),
	}
	new_room.clients[admin.ID] = admin
	new_room.applicants <- admin
	return new_room
}

func (room *Room) addClientInRoom(client *Client) error {
	select {
	case room.applicants <- client:
		room.clients[client.ID] = client
	default:
		return fmt.Errorf("there are no vancancy in the room")
	}
	return nil
}

func (room *Room) broadcastToClients(message *Message) {
	for _, client := range room.clients {
		client.send <- message.encode()
	}
}

func (room *Room) hasVancancy() bool {
	return len(room.clients) < room.max_clients_number
}
