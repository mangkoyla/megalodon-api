package proxy

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	database "github.com/FoolVPN-ID/megalodon/db"
	CS "github.com/sagernet/sing-box/constant"
)

func ConvertDBToURL(field *database.ProxyFieldStruct) *url.URL {
	var rawUrl *url.URL

	switch field.VPN {
	case CS.TypeTrojan:
		rawUrl = trojanToRaw(field)
	case CS.TypeVLESS:
		rawUrl = vlessToRaw(field)
	case CS.TypeVMess:
		rawUrl = vmessToRaw(field)
	case CS.TypeShadowsocks:
		rawUrl = ssToRaw(field)
	}

	return rawUrl
}

func trojanToRaw(field *database.ProxyFieldStruct) *url.URL {
	rawLink, _ := url.Parse(fmt.Sprintf("trojan://%s@%s:%d", field.Password, field.Server, field.ServerPort))
	rawLinkQuery := rawLink.Query()

	// Common value
	rawLink.Fragment = field.Remark
	rawLinkQuery.Add("host", field.Host)

	// Network
	if field.TLS {
		rawLinkQuery.Add("security", "tls")
		rawLinkQuery.Add("sni", field.SNI)
		rawLinkQuery.Add("allowInsecure", "1")
	} else {
		rawLinkQuery.Add("security", "none")
	}

	switch field.Transport {
	case "ws":
		rawLinkQuery.Add("path", field.Path)
		rawLinkQuery.Add("type", field.Transport)
	case "grpc":
		rawLinkQuery.Add("serviceName", field.ServiceName)
		rawLinkQuery.Add("mode", "gun")
		rawLinkQuery.Add("type", field.Transport)
	default:
		rawLinkQuery.Add("type", "tcp")
	}

	rawLink.RawQuery = rawLinkQuery.Encode()
	return rawLink
}

func vlessToRaw(field *database.ProxyFieldStruct) *url.URL {
	rawLink, _ := url.Parse(fmt.Sprintf("vless://%s@%s:%d", field.UUID, field.Server, field.ServerPort))
	rawLinkQuery := rawLink.Query()

	// Common value
	rawLink.Fragment = field.Remark
	rawLinkQuery.Add("host", field.Host)

	// Network
	if field.TLS {
		rawLinkQuery.Add("security", "tls")
		rawLinkQuery.Add("sni", field.SNI)
		rawLinkQuery.Add("allowInsecure", "1")
	} else {
		rawLinkQuery.Add("security", "none")
	}

	switch field.Transport {
	case "ws":
		rawLinkQuery.Add("path", field.Path)
		rawLinkQuery.Add("type", field.Transport)
	case "grpc":
		rawLinkQuery.Add("serviceName", field.ServiceName)
		rawLinkQuery.Add("mode", "gun")
		rawLinkQuery.Add("type", field.Transport)
	default:
		rawLinkQuery.Add("type", "tcp")
	}

	rawLink.RawQuery = rawLinkQuery.Encode()
	return rawLink
}

func vmessToRaw(field *database.ProxyFieldStruct) *url.URL {
	rawLink, _ := url.Parse(fmt.Sprintf("vmess://%s", field.UUID))

	// Some confusing protocol
	params := map[string]any{
		"v":    2,
		"ps":   field.Remark,
		"add":  field.Server,
		"port": field.ServerPort,
		"id":   field.UUID,
		"aid":  field.AlterId,
		"scy":  field.Security,
		"type": "none",
		"host": field.Host,
		"sni":  field.SNI,
	}

	// Network
	if field.TLS {
		params["tls"] = "tls"
	}

	switch field.Transport {
	case "ws":
		params["net"] = field.Transport
		params["path"] = field.Path
	case "grpc":
		params["net"] = field.Transport
		params["path"] = field.ServiceName
	}

	// Encoding
	paramsByte, _ := json.Marshal(params)
	paramsStr := base64.RawStdEncoding.EncodeToString(paramsByte)

	rawLink.Host = paramsStr
	return rawLink
}

func ssToRaw(field *database.ProxyFieldStruct) *url.URL {
	cred := base64.RawStdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", field.Method, field.Password)))
	rawLink, _ := url.Parse(fmt.Sprintf("ss://%s@%s:%d", cred, field.Server, field.ServerPort))
	rawLinkQuery := rawLink.Query()

	// Common value
	rawLink.Fragment = field.Remark

	// Network
	if field.Plugin != "" {
		rawLinkQuery.Add("plugin", fmt.Sprintf("%s;%s", field.Plugin, field.PluginOpts))
	}

	rawLink.RawQuery = rawLinkQuery.Encode()
	return rawLink
}
