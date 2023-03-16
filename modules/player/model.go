package player

import (
  "fmt"
	"golang.org/x/crypto/bcrypt"
	"online_chess/model"
	"online_chess/util"
)

func SetPlayerPassword(new_password []byte) {
	model.GetDBInstance().Exec("UPDATE player SET password = ? WHERE name = 'wilson'", new_password)
}

func GetPlayerPassword (email string) ([]byte, error) {
  var password []byte
	if err := model.GetDBInstance().QueryRow("SELECT password FROM player WHERE email = ?", email).Scan(&password); err != nil {
		return nil, fmt.Errorf("QUERY player fail %v", err)
	}
  return password, nil
}

func AddNewPlayer(inputs Register_player) (string, error) {
	hash_password, _ := bcrypt.GenerateFromPassword([]byte(inputs.Password), 10)
	res, err := model.GetDBInstance().Exec(
		"INSERT INTO player ( name, email, password ) VALUES ( ?, ?, ? )",
		inputs.Name,
		inputs.Email,
		hash_password,
	)
	fmt.Println("res:", res)
	fmt.Println("err:", err)
	if err != nil {
		return "", err
	} else {
		token := util.GenerateToken()
		return token, nil
	}
}