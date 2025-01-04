package users

type UserStruct struct {
	ID         uint64 `json:"id"`
	Token      string `json:"token"`
	Password   string `json:"password"`
	Expired    string `json:"expired"`
	ServerCode string `json:"server_code,omitempty"`
	Quota      int    `json:"quota"`
	Relay      string `json:"relay,omitempty"`
	Adblock    bool   `json:"adblock"`
	VPN        string `json:"vpn,omitempty"`
}
