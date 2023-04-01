package room

import(
	"online_chess/model"
)

func AddNewRoom(inputs create_room_input, admin_id int) (room, error) {
	var new_room room
	res, err := model.GetDBInstance().Exec(
		"INSERT INTO room ( name, admin_id, password ) VALUES ( ?, ?, ? )",
		inputs.Name,
		admin_id,
		inputs.Password,
	)
	if err == nil {
		var new_room_id int64
		new_room_id, err = res.LastInsertId()
		new_room.Id, new_room.Name, new_room.Password = int(new_room_id), inputs.Name, inputs.Password
	}
	return new_room, err
}