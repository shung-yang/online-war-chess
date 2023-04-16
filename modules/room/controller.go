package room

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"online_chess/modules/player"
	"strings"
	"sync"
)

type create_room_input struct {
	Name string `json:"name" binding:"required,max=50" example:"test room name"`
  Password string  `json:"password" example:"123123123"`
}

type room struct {
	Id int `json:"id" example:"1"`
	Name string `json:"name" example:"room name"`
	Password	string `json:"password" example:"room password"`
	Admin_name string `json:"admin_name" example:"wilson"`
	Admin_level int8 `json:"admin_level" example:"3"`
	Other_player_name string `json:"other_player_name" example:"sam"`
	Other_player_level int8 `json:"other_player_level" example:"2"`
}

type room_list_item struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Admin_id int `json:"admin_id"`
	Other_player_id int `json:"other_player_id"`
}

// @Summary      create room
// @Description  player create room to play with the other player
// @Tags         room
// @Accept       json
// @Produce      json
// @Param create_room_input body create_room_input true "create room"
// @Success      200  {object}  room
// @Failure      401  {object}  object{error=string} "player token not valid"
// @Failure      500  {object}  object{error=string}
// @Router       /room [post]
func Create(c *gin.Context) {
	var inputs create_room_input
	admin_id, err := player.QueryIdByToken(strings.Split(c.GetHeader("Authorization"), " ")[1])
	err = c.ShouldBindJSON(&inputs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
	} else {
		new_room, err := AddNewRoom(inputs, admin_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
		}
		player, err := player.QueryPlayerById(admin_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
		} else {
			new_room.Admin_name, new_room.Admin_level = player.Name, player.Level
			c.JSON(http.StatusOK, new_room)
		}
	}
}

type room_list_query struct {
	Page int `form:"page" binding:"numeric,gt=0"`
}
// @Summary      get room list
// @Description  get room list according to page query
// @Tags         room
// @Accept       json
// @Produce      json
// @Param page query string true "get room list by page number"
// @Success      200  {array}  room_list_item
// @Failure      401  {object}  object{error=string} "player token not valid"
// @Failure      500  {object}  object{error=string}
// @Router       /room [get]
func GetList(c *gin.Context) {
	var query room_list_query
	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
	} else {
		room_list, err := GetRoomList(query.Page - 1, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
		} else {
			c.JSON(http.StatusOK, room_list)
		}
	}
}

func ValidRoomCanJoin(player_id, room_id int) error {
	var valid_room_err error
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		is_exist, err := CheckRoomIsExist(room_id)
		if !is_exist {
			valid_room_err = fmt.Errorf("The room you are trying to join does not exist")
		} else if err != nil {
			valid_room_err = err
		}
	}()

	go func() {
		defer wg.Done()
		is_in_room, _, err := CheckPlayerIsInRoom(player_id)
		if is_in_room {
			valid_room_err = fmt.Errorf("You already join the other room")
		} else if err != nil {
			valid_room_err = err
		}
	}()

	go func() {
		defer wg.Done()
		vancancy_in_room, err := CheckVacancyInRoom(room_id)
		if !vancancy_in_room {
			valid_room_err = fmt.Errorf("The room you want to attend is no longer available")
		} else if err != nil {
			valid_room_err = err
		}
	}()
	wg.Wait()
	return valid_room_err
}

type join_query struct {
	Room int `form:"room" binding:"numeric"`
}
// @Summary      player joins room
// @Description  player joins room
// @Tags         room
// @Accept       json
// @Produce      json
// @Param page query uint true "player join room by room id"
// @Success      200  {array}  string "success to join room"
// @Failure      401  {object}  object{error=string}
// @Failure      500  {object}  object{error=string}
// @Router       /room/join [post]
func Join(c *gin.Context) {
	var query join_query
	player_id, err := player.QueryIdByToken(strings.Split(c.GetHeader("Authorization"), " ")[1])
	err = c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
	}
	valid_room_err := ValidRoomCanJoin(player_id, query.Room)
	if valid_room_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": valid_room_err.Error() })
	} else {
		err := PlayerJoinRoom(player_id, query.Room)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
		} else {
			c.JSON(http.StatusOK, "success to join room")
		}
	}
}