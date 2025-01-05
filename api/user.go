package api

import (
	"strconv"

	"github.com/FoolVPN-ID/megalodon-api/modules/db/kv"
	"github.com/FoolVPN-ID/megalodon-api/modules/db/users"
	"github.com/gin-gonic/gin"
)

func handleGetUserApi(c *gin.Context) {
	// Validate api token
	apiToken := c.Param("apiToken")
	if apiToken == "" {
		c.String(403, "token invalid")
		return
	}

	kvClient := kv.MakeKVTableClient()
	validToken, err := kvClient.GetValueFromKVByKey("apiToken")
	if err != nil {
		c.String(500, err.Error())
		return
	} else if *validToken == "" {
		c.String(500, "token not set")
		return
	} else if *validToken != apiToken {
		c.String(403, "token invalid")
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
