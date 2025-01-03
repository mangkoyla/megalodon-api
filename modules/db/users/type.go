package users

type UserStruct struct {
	ID         int    `json:"id"`
	Token      string `json:"token"`
	Password   string `json:"password"`
	Expired    string `json:"expired"`
	DomainCode string `json:"domain_code,omitempty"`
	Quota      int    `json:"quota"`
	Relay      string `json:"relay,omitempty"`
	Adblock    bool   `json:"adblock"`
	VPN        string `json:"vpn,omitempty"`
}
