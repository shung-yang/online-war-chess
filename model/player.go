package model

import (
  "fmt"
)

func SetPlayerPassword(new_password []byte) {
	db.Exec("UPDATE player SET password = ? WHERE name = 'wilson'", new_password)
}

func GetPlayerPassword (email string) ([]byte, error) {
  var password []byte
	if err := db.QueryRow("SELECT password FROM player WHERE email = ?", email).Scan(&password); err != nil {
		return nil, fmt.Errorf("QUERY player fail %v", err)
	}
  return password, nil
}