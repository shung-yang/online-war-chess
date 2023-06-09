package game

import (
	"github.com/google/uuid"
)

type WsServer struct {
	rooms     map[uuid.UUID]*Room
	clients   map[int]*Client
	broadcast chan *Message
}

func NewWebSocketServer() *WsServer {
	return &WsServer{
		rooms:     make(map[uuid.UUID]*Room),
		clients:   make(map[int]*Client),
		broadcast: make(chan *Message, 10),
	}
}

func (server *WsServer) Run() {
	for {
		select {
		case message := <-server.broadcast:
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) addClient(client *Client) {
	server.clients[client.ID] = client
}

func (server *WsServer) createRoom(name string, admin *Client) *Room {
	room := NewRoom(name, admin)
	admin.room = room
	go room.Run()
	server.rooms[room.ID] = room
	return room
}

func (server *WsServer) findRoomByName(name string) *Room {
	var found_room *Room
	for _, room := range server.rooms {
		if room.Name == name {
			found_room = room
			break
		}
	}
	return found_room
}

func (server *WsServer) findRoomByID(room_id uuid.UUID) (bool, *Room) {
	room, is_room_exist := server.rooms[room_id]
	return is_room_exist, room
}

func (server *WsServer) broadcastToClients(message *Message) {
	for _, client := range server.clients {
		client.send <- message.encode()
	}
}

func (server *WsServer) GetRoomList() []*Room {
	var rooms []*Room
	for _, room := range server.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}
