package database

import (
	"fmt"

	mgdb "github.com/FoolVPN-ID/megalodon/db"
)

func (db *databaseStruct) GetProxiesByCondition(condition string) ([]mgdb.ProxyFieldStruct, error) {
	var results []mgdb.ProxyFieldStruct
	rows, err := db.client.Query(fmt.Sprintf("SELECT * FROM proxies WHERE %s;", condition))
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var (
			result = mgdb.ProxyFieldStruct{}
			id     int
		)

		err := rows.Scan(
			&id,
			&result.Server,
			&result.Ip,
			&result.ServerPort,
			&result.UUID,
			&result.Password,
			&result.Security,
			&result.AlterId,
			&result.Method,
			&result.Plugin,
			&result.PluginOpts,
			&result.Host,
			&result.TLS,
			&result.Transport,
			&result.Path,
			&result.ServiceName,
			&result.Insecure,
			&result.SNI,
			&result.Remark,
			&result.ConnMode,
			&result.CountryCode,
			&result.Region,
			&result.Org,
			&result.VPN,
			&result.Raw,
		)

		if err != nil {
			return results, err
		}

		results = append(results, result)
	}

	return results, err
}
