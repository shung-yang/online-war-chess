package room

import(
	"errors"
	"fmt"
	"online_chess/model"
	"database/sql"
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

func CheckRoomIsExist(room_id int) (bool, error) {
	var search_result int
	err := model.GetDBInstance().QueryRow("SELECT id FROM room WHERE id = ?", room_id).Scan(&search_result);
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("Fail to check room is full or not, %v", err)
	} else {
		return true, nil
	}
}

func CheckPlayerIsInRoom(player_id int) (bool, int, error) {
	var room_id int
	err := model.GetDBInstance().QueryRow("SELECT id FROM room WHERE admin_id = ? OR other_player_id = ?", player_id, player_id).Scan(&room_id)
	if errors.Is(err, sql.ErrNoRows) {
		return false, 0, nil
	} else if err != nil {
		return false, 0, fmt.Errorf("Fail to check player is in room or not, %v", err)
	} else {
		return true, room_id, nil
	}
}

func CheckVacancyInRoom(room_id int) (bool, error) {
	var other_player_in_room int
	err := model.GetDBInstance().QueryRow("SELECT other_player_id FROM room WHERE room_id = ?", room_id).Scan(&other_player_in_room)
	if other_player_in_room == 0 {
		return true, nil
	} else if err != nil {
		return false, fmt.Errorf("Fail to check vancancy in room or not, %v", err)
	} else {
		return true, nil
	}
	
}

func PlayerJoinRoom(player_id, room_id int) error {
	if _, err := model.GetDBInstance().Exec("UPDATE room SET other_player_id = ? WHERE id = ?", player_id, room_id); err != nil {
		return fmt.Errorf("Player fails to join the room %v", err)
	}
	return nil
}