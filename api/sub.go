package api

import (
	"fmt"
	"strings"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
	"github.com/gin-gonic/gin"
)

type apiGetSubStruct struct {
	Pass      string `form:"pass" binding:"omitempty"`
	Free      int8   `form:"free" binding:"omitempty"`
	Premium   int8   `form:"premium" binding:"omitempty"`
	VPN       string `form:"vpn" binding:"omitempty"`
	Format    string `form:"format" binding:"omitempty"`
	Region    string `form:"region" binding:"omitempty"`
	CC        string `form:"cc" binding:"omitempty"`
	Include   string `form:"include" binding:"omitempty"`
	Exclude   string `form:"exclude" binding:"omitempty"`
	TLS       int8   `form:"tls" binding:"omitempty"`
	Transport string `form:"transport" binding:"omitempty"`
	IP        int8   `form:"ip" binding:"omitempty"`
	SNI       string `form:"sni" binding:"omitempty"`
	CDN       string `form:"cdn" binding:"omitempty"`
	Mode      string `form:"mode" binding:"omitempty"`
	Limit     int    `form:"limit" binding:"omitempty"`
	Subdomain string `form:"subdomain" binding:"omitempty"`
}

type whereConditionObject struct {
	conditions []string
	delimiter  string
}

func HandleSubApi(c *gin.Context) {
	var getQuery apiGetSubStruct
	err := c.ShouldBindQuery(&getQuery)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	// Check api password
	if getQuery.Pass == "" {
		c.String(403, "API password not provided!")
		return
	}
	// Check password from database

	condition := buildSqlWhereCondition(getQuery)

	db := database.MakeDatabase()
	proxies, err := db.GetProxiesByCondition(condition)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.JSON(200, proxies)
}

func buildSqlWhereCondition(getQuery apiGetSubStruct) string {
	var (
		limit         = 10
		conditionList = []whereConditionObject{}
	)

	if getQuery.Limit > 0 && getQuery.Limit <= 10 {
		limit = getQuery.Limit
	}

	if getQuery.VPN != "" {
		conditionList = append(conditionList, buildCondition("VPN", getQuery.VPN, "=", " OR "))
	}
	if getQuery.Region != "" {
		conditionList = append(conditionList, buildCondition("REGION", getQuery.Region, "=", " OR "))
	}
	if getQuery.CC != "" {
		conditionList = append(conditionList, buildCondition("COUNTRY_CODE", getQuery.CC, "=", " OR "))
	}
	if getQuery.Transport != "" {
		conditionList = append(conditionList, buildCondition("TRANSPORT", getQuery.Transport, "=", " OR "))
	}
	if getQuery.Mode != "" {
		conditionList = append(conditionList, buildCondition("CONN_MODE", getQuery.Mode, "=", " OR "))
	}
	if getQuery.Include != "" {
		conditionList = append(conditionList, buildCondition("REMARK", strings.ToUpper(getQuery.Include), "LIKE", " OR "))
	}
	if getQuery.Exclude != "" {
		conditionList = append(conditionList, buildCondition("REMARK", strings.ToUpper(getQuery.Exclude), "NOT LIKE", " OR "))
	}
	if getQuery.TLS >= 0 {
		conditionList = append(conditionList, whereConditionObject{
			conditions: []string{fmt.Sprintf("TLS = %d", getQuery.TLS)},
			delimiter:  "",
		})
	}

	whereConditions := []string{}
	for _, cl := range conditionList {
		whereConditions = append(whereConditions, "("+strings.Join(cl.conditions, cl.delimiter)+")")
	}

	return strings.Join(whereConditions, " AND ") + fmt.Sprintf(" ORDER BY RANDOM() LIMIT %d", limit)
}

func buildCondition(key, value, operator, delimiter string) whereConditionObject {
	condition := whereConditionObject{
		delimiter: delimiter,
	}

	for _, v := range strings.Split(value, ",") {
		condition.conditions = append(condition.conditions, fmt.Sprintf("%s %s '%s'", key, operator, v))
	}

	return condition
}
