package shared

type ClientHello struct {
	Username string `json:"username"`
}

type ServerHello struct {
	RoomCount  int    `json:"roomCount"`
	WelcomeMsg string `json:"welcomeMsg"`
}

type Msg struct {
	Username string `json:"username"`
	Msg      string `json:"msg"`
	Join     bool   `json:"join"`
	Leave    bool   `json:"leave"`
}
