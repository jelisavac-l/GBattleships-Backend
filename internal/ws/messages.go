package ws

import "github.com/jelisavac-l/GBattleships/internal/model"

type WSMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type GameStartedMessage struct {
	ID1       string `json:"id1"`
	Username1 string `json:"username1"`
	ID2       string `json:"id2"`
	Username2 string `json:"username2"`
	State     string `json:"state"`
}

type GetBoardMessage struct {
}
type SendBoardMessage struct {
	Cells [][]model.CellState `json:"cells"`
}

type GetTurnMessage struct {
	X   int  `json:"x"`
	Y   int  `json:"y"`
	Hit bool `json:"hit"`
}
type SendTurnMessage struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type TurnResultMessage struct {
	Hit bool `json:"hit"`
}

type GameResultMessage struct {
	WinnerUname string `json:"winneruname"`
}
type RematchMessage struct {
	WantsRematch bool `json:"wantsrematch"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
