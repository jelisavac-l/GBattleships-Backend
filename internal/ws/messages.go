package ws

import "github.com/jelisavac-l/GBattleships/internal/model"

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type gameStartedMessage struct {
	ID1       string `json:"id1"`
	Username1 string `json:"username1"`
	ID2       string `json:"id2"`
	Username2 string `json:"username2"`
	State     string `json:"state"`
}

type getBoardMessage struct {
}
type sendBoardMessage struct {
	Cells [][]model.CellState `json:"cells"`
}

type getTurnMessage struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type sendTurnMessage struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type turnResultMessage struct {
	Hit bool `json:"hit"`
}

type gameResultMessage struct {
	WinnerUname string `json:"winneruname"`
}
type rematchMessage struct {
	WantsRematch bool `json:"wantsrematch"`
}

type errorMessage struct {
	Error string `json:"error"`
}
