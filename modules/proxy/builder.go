package proxy

import (
	"fmt"
	"strings"

	"github.com/FoolVPN-ID/megalodon-api/modules/db/users"
	"github.com/FoolVPN-ID/megalodon/common/helper"
	database "github.com/FoolVPN-ID/megalodon/db"
	"github.com/Noooste/azuretls-client"
)

var (
	ModeList      = []string{"cdn", "sni"}
	TransportList = []string{"ws", "grpc", "tcp"}
	PortList      = []int{80, 443}
)

func BuildProxyFieldsFromUser(user *users.UserStruct, baseProxyField database.ProxyFieldStruct) []database.ProxyFieldStruct {
	var (
		results  = []database.ProxyFieldStruct{}
		session  = azuretls.NewSession()
		res, err = session.Get("https://" + baseProxyField.Server + "/api/v1/info")

		serverInfo = map[string]string{
			"country":      "XX",
			"country_name": "oceania",
			"ip":           "192.168.1.1",
			"org":          "Megalodon",
			"region":       "Hindia",
		}
	)

	if err == nil && res.StatusCode == 200 {
		res.MustJSON(&serverInfo)
	}

	for _, mode := range ModeList {
		for _, transport := range TransportList {
			for _, port := range PortList {
				// Filters
				if mode == "cdn" {
					if (transport == "grpc" && port == 80) || (transport == "tcp") {
						continue
					}
				} else {
					if port == 80 {
						continue
					}
				}

				tlsStr := "TLS"
				if port == 80 {
					tlsStr = "NTLS"
				}

				proxyField := baseProxyField
				proxyField.Ip = serverInfo["ip"]
				proxyField.ServerPort = port
				proxyField.Transport = transport
				proxyField.TLS = port == 443
				proxyField.Path = "/" + user.VPN
				proxyField.ServiceName = user.VPN
				proxyField.ConnMode = mode
				proxyField.Remark = strings.ToUpper(fmt.Sprintf("%d ✨ %s %s %s %s %s", len(results)+1, helper.CCToEmoji(serverInfo["country"]), serverInfo["org"], transport, mode, tlsStr))

				results = append(results, proxyField)
			}
		}
	}

	return results
}
