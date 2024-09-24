package entities

type Message struct {
	Message string `json:"message"`
	Player  Player `json:"player"`
}

type Player struct {
	UUID     string `json:"uuid"`
	Username string `json:"Username"`
}
