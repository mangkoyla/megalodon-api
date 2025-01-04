package api

import (
	"strconv"

	"github.com/FoolVPN-ID/megalodon-api/modules/db/users"
	"github.com/gin-gonic/gin"
)

func handleGetUserApi(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(400, err.Error())
		return
	}

	usersTable := users.MakeUsersTableClient()
	user, err := usersTable.GetUserByIdOrToken(userID, nil)

	if err == nil && user != nil {
		c.JSON(200, user)
		return
	} else {
		err := usersTable.NewUser(userID)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		handleGetUserApi(c)
		return
	}
}
