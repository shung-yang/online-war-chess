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

func GetRoomList(page, rows_per_page int) ([]room_list_item, error) {
	if rows_per_page == 0 {
		rows_per_page = 10
	}
	var rooms []room_list_item
	rows, err := model.GetDBInstance().Query("SELECT id, name, admin_id, COALESCE(other_player_id, 0) from room LIMIT ?, ?", page, rows_per_page)
	defer rows.Close()
	if err != nil {
	  return nil, err
	}
	for rows.Next() {
		var new_room room_list_item
		if err := rows.Scan(&new_room.Id, &new_room.Name, &new_room.Admin_id, &new_room.Other_player_id); err != nil {
			return nil, err
		} else {
			rooms = append(rooms, new_room)
		}
	}
	return rooms, err
}