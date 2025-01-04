package servers

type ServerStruct struct {
	ID         int    `json:"id"`
	Code       string `json:"code"`
	Domain     string `json:"domain"`
	IP         string `json:"ip"`
	Country    string `json:"country"`
	UsersCount int    `json:"users_count"`
	UsersMax   int    `json:"users_max"`
}
