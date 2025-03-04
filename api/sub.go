package api

import (
	"fmt"
	"math/rand"
	"strings"
	"github.com/FoolVPN-ID/megalodon-api/modules/db"
	"github.com/FoolVPN-ID/megalodon-api/modules/proxy"
	mgdb "github.com/FoolVPN-ID/megalodon/db"
	"github.com/FoolVPN-ID/tool/modules/subconverter"
	"github.com/gin-gonic/gin"
)

type apiGetSubStruct struct {
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

func handleGetSubApi(c *gin.Context) {
	var (
		getQuery apiGetSubStruct
		proxies  = []mgdb.ProxyFieldStruct{}
	)

	err := c.ShouldBindQuery(&getQuery)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	// Re-assign non string query
	if c.Query("tls") == "" {
		getQuery.TLS = -1
	}

	// Removed the password check (pass) and free
	// The logic to check the 'pass' and 'free' has been removed

	// Build the SQL condition query
	condition := buildSqlWhereCondition(getQuery)

	// Access the database with the conditions
	db := database.MakeDatabase()
	proxies, err = db.GetProxiesByCondition(condition)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Assign domain based on CDN or SNI
	var (
		cdnDomains = strings.Split(getQuery.CDN, ",")
		sniDomains = strings.Split(getQuery.SNI, ",")
	)
	for i := range proxies {
		proxy := &proxies[i]
		switch proxy.ConnMode {
		case "cdn":
			if cdnDomains[0] != "" {
				cdnDomain := cdnDomains[rand.Intn(len(cdnDomains))]
				proxy.Server = cdnDomain
			}
		case "sni":
			if sniDomains[0] != "" {
				sniDomain := sniDomains[rand.Intn(len(sniDomains))]
				proxy.SNI = sniDomain
				proxy.Host = sniDomain
			}
		}
	}

	// Convert proxies to raw format
	rawProxies := []string{}
	for _, dbProxy := range proxies {
		rawProxies = append(rawProxies, proxy.ConvertDBToURL(&dbProxy).String())
	}
	subProxies, err := subconverter.MakeSubconverterFromConfig(strings.Join(rawProxies, "\n"))
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Return response based on the requested format
	switch getQuery.Format {
	case "raw":
		c.String(200, strings.Join(rawProxies, "\n"))
		return
	case "sfa":
		if err := subProxies.ToSFA(); err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, subProxies.Result.SFA)
		return
	case "bfr":
		if err := subProxies.ToBFR(); err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, subProxies.Result.BFR)
		return
	case "sing-box":
		c.JSON(200, subProxies.Outbounds)
		return
	case "clash":
		if err := subProxies.ToClash(); err != nil {
			c.String(500, err.Error())
			return
		}
		c.YAML(200, subProxies.Result.Clash)
		return
	default:
		c.JSON(200, proxies)
	}
}

func buildSqlWhereCondition(getQuery apiGetSubStruct) string {
	var (
		limit         = 50
		conditionList = []whereConditionObject{}
	)

	// Set limit based on the query
	if getQuery.Limit > 0 && getQuery.Limit <= 50 {
		limit = getQuery.Limit
	}

	// Add conditions based on query parameters
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
		conditionList = append(conditionList, buildCondition("REMARK", "%%"+strings.ToUpper(getQuery.Include)+"%%", "LIKE", " OR "))
	}
	if getQuery.Exclude != "" {
		conditionList = append(conditionList, buildCondition("REMARK", "%%"+strings.ToUpper(getQuery.Exclude)+"%%", "NOT LIKE", " OR "))
	}
	if getQuery.TLS >= 0 {
		conditionList = append(conditionList, whereConditionObject{
			conditions: []string{fmt.Sprintf("TLS = %d", getQuery.TLS)},
			delimiter:  "",
		})
	}

	// Combine conditions into final SQL condition
	whereConditions := []string{}
	for _, cl := range conditionList {
		whereConditions = append(whereConditions, "("+strings.Join(cl.conditions, cl.delimiter)+")")
	}

	finalCondition := strings.Join(whereConditions, " AND ")
	if finalCondition != "" {
		finalCondition = "WHERE " + finalCondition
	}

	return finalCondition + fmt.Sprintf(" ORDER BY RANDOM() LIMIT %d", limit)
}

func buildCondition(key, value, operator, delimiter string) whereConditionObject {
	condition := whereConditionObject{
		delimiter: delimiter,
	}

	// Split values and build condition
	for _, v := range strings.Split(value, ",") {
		condition.conditions = append(condition.conditions, fmt.Sprintf("%s %s '%s'", key, operator, v))
	}

	return condition
}
