package api

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
	"github.com/FoolVPN-ID/megalodon-api/modules/db/servers"
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

func handleGetSubApi(c *gin.Context) {
	var (
		getQuery apiGetSubStruct
		proxies  []mgdb.ProxyFieldStruct
	)

	err := c.ShouldBindQuery(&getQuery)
	if err != nil {
		c.String(400, err.Error())
		return
	}

	// Menghapus semua batasan Free/Premium, semua pengguna mendapatkan data premium
	db := database.MakeDatabase()
	proxies, err = db.GetProxiesByCondition(buildSqlWhereCondition(getQuery))
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Mendapatkan server premium secara otomatis
	server, err := servers.MakeServersTableClient().GetRandomPremiumServer()
	if err == nil {
		basePremiumProxy := mgdb.ProxyFieldStruct{
			Server:      server.Domain,
			Ip:          server.IP,
			UUID:        "PREMIUM-USER",
			Password:    "PREMIUM-USER",
			Host:        server.Domain,
			Insecure:    true,
			SNI:         server.Domain,
			CountryCode: server.Country,
			VPN:         "Premium",
		}
		proxies = append([]mgdb.ProxyFieldStruct{basePremiumProxy}, proxies...)
	}

	// Menyesuaikan SNI dan CDN jika diperlukan
	cdnDomains := strings.Split(getQuery.CDN, ",")
	sniDomains := strings.Split(getQuery.SNI, ",")
	for i := range proxies {
		switch proxies[i].ConnMode {
		case "cdn":
			if cdnDomains[0] != "" {
				proxies[i].Server = cdnDomains[rand.Intn(len(cdnDomains))]
			}
		case "sni":
			if sniDomains[0] != "" {
				proxies[i].SNI = sniDomains[rand.Intn(len(sniDomains))]
				proxies[i].Host = proxies[i].SNI
			}
		}
	}

	// Mengubah format output sesuai permintaan
	rawProxies := []string{}
	for _, dbProxy := range proxies {
		rawProxies = append(rawProxies, proxy.ConvertDBToURL(&dbProxy).String())
	}
	subProxies, err := subconverter.MakeSubconverterFromConfig(strings.Join(rawProxies, "\n"))
	if err != nil {
		c.String(500, err.Error())
		return
	}

	switch getQuery.Format {
	case "raw":
		c.String(200, strings.Join(rawProxies, "\n"))
	case "sfa":
		if err := subProxies.ToSFA(); err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, subProxies.Result.SFA)
	case "bfr":
		if err := subProxies.ToBFR(); err != nil {
			c.String(500, err.Error())
			return
		}
		c.JSON(200, subProxies.Result.BFR)
	case "sing-box":
		c.JSON(200, subProxies.Outbounds)
	case "clash":
		if err := subProxies.ToClash(); err != nil {
			c.String(500, err.Error())
			return
		}
		c.YAML(200, subProxies.Result.Clash)
	default:
		c.JSON(200, proxies)
	}
}

func buildSqlWhereCondition(getQuery apiGetSubStruct) string {
	var (
		limit         = 10
		conditionList []string
	)

	if getQuery.Limit > 0 && getQuery.Limit <= 10 {
		limit = getQuery.Limit
	}

	if getQuery.VPN != "" {
		conditionList = append(conditionList, fmt.Sprintf("VPN = '%s'", getQuery.VPN))
	}
	if getQuery.Region != "" {
		conditionList = append(conditionList, fmt.Sprintf("REGION = '%s'", getQuery.Region))
	}
	if getQuery.CC != "" {
		conditionList = append(conditionList, fmt.Sprintf("COUNTRY_CODE = '%s'", getQuery.CC))
	}
	if getQuery.Transport != "" {
		conditionList = append(conditionList, fmt.Sprintf("TRANSPORT = '%s'", getQuery.Transport))
	}
	if getQuery.Mode != "" {
		conditionList = append(conditionList, fmt.Sprintf("CONN_MODE = '%s'", getQuery.Mode))
	}
	if getQuery.Include != "" {
		conditionList = append(conditionList, fmt.Sprintf("REMARK LIKE '%%%s%%'", strings.ToUpper(getQuery.Include)))
	}
	if getQuery.Exclude != "" {
		conditionList = append(conditionList, fmt.Sprintf("REMARK NOT LIKE '%%%s%%'", strings.ToUpper(getQuery.Exclude)))
	}
	if getQuery.TLS >= 0 {
		conditionList = append(conditionList, fmt.Sprintf("TLS = %d", getQuery.TLS))
	}

	finalCondition := strings.Join(conditionList, " AND ")
	if finalCondition != "" {
		finalCondition = "WHERE " + finalCondition
	}

	return finalCondition + fmt.Sprintf(" ORDER BY RANDOM() LIMIT %d", limit)
}
