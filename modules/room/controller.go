package room

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_chess/modules/player"
	"strings"
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
		c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
	} else {
		room_list, err := GetRoomList(query.Page - 1, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
		} else {
			c.JSON(http.StatusOK, room_list)
		}
	}
}