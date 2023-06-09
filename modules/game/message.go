package game

import (
	"encoding/json"
	"log"
)

const (
	connect_websocket_success = "connect_websocket_success"
	join_room_action          = "join_room"
	join_room_success         = "join_room_success"
	join_room_failed          = "join_room_failed"
	create_room_action        = "create_room"
	create_room_success       = "create_room_success"
	create_room_failed        = "create_room_failed"
	get_rooms_action          = "get_rooms"
	get_rooms_success         = "get_rooms_success"
)

type Message struct {
	Action  string      `json:"action"`
	Payload interface{} `json:"payload"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
