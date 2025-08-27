package game

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/jelisavac-l/GBattleships/internal/model"
	"github.com/jelisavac-l/GBattleships/internal/ws"
)

type Game struct {
	ID             string
	Player1        *model.Player
	Player2        *model.Player
	Board1         *model.Board
	Board2         *model.Board
	Turn           bool // true for player1 	false for player2
	State          string
	winnerUsername string
	Wg             sync.WaitGroup
}

func CreateGame(player model.Player, id int) *Game {
	return &Game{
		ID:      strconv.Itoa(id),
		State:   "waiting",
		Player1: &player,
	}
}

func (game *Game) StartGame() bool {
	var wg sync.WaitGroup
	wg.Add(2)
	go game.getBoard(game.Player1, &wg)
	go game.getBoard(game.Player2, &wg)
	wg.Wait()

	game.State = "playing"
	game.tellGameStarted()
	log.Println("Game " + game.ID + " state changed to playing")

	var hit bool
	var err error
	var previousX, previousY int
	for game.State != "finished" {
		x, y := game.getMove(game.Turn, previousX, previousY)
		if game.checkValidMove(x, y, !game.Turn) {
			hit, err = game.PlayMove(x, y)
			if err != nil {
				var player model.Player
				if game.Turn {
					player = *game.Player1
				} else {
					player = *game.Player2
				}
				sendErrorMessage(err.Error(), &player)
				continue
			}
			game.sendHitOrMiss(hit)
		} else {
			var player model.Player
			if game.Turn {
				player = *game.Player1
			} else {
				player = *game.Player2
			}
			sendErrorMessage("x or y out of bounds", &player)
			continue
		}
		game.State = game.checkState()
		previousX = x
		previousY = y
		game.Turn = !game.Turn
	}

	return game.tellResultsAskRematch()
}

func (game *Game) checkState() string {
	if game.Turn { // player1 turn
		if game.Board2.Hits == 17 {
			game.winnerUsername = game.Player1.Username
			return "finished"
		}
	} else if !game.Turn { // player2 turn
		if game.Board1.Hits == 17 {
			game.winnerUsername = game.Player2.Username
			return "finished"
		}
	}
	return "playing"

}

func checkValidBoard(board model.Board) bool {
	dim := len(board.Cells)
	if dim == 0 {
		return false
	}

	visited := make([][]bool, dim)
	for i := range visited {
		visited[i] = make([]bool, len(board.Cells[i]))
	}

	required := map[int]int{
		5: 1,
		4: 1,
		3: 2,
		2: 1,
	}
	found := map[int]int{}

	var dfs func(x, y int) int
	dfs = func(x, y int) int {
		if x < 0 || x >= dim || y < 0 || y >= dim {
			return 0
		}
		if visited[x][y] || board.Cells[x][y] != model.Ship {
			return 0
		}
		visited[x][y] = true

		length := 1

		length += dfs(x+1, y)
		length += dfs(x-1, y)
		length += dfs(x, y+1)
		length += dfs(x, y-1)

		return length
	}

	for i := 0; i < dim; i++ {
		for j := 0; j < dim; j++ {
			if !visited[i][j] && board.Cells[i][j] == model.Ship {
				shipLen := dfs(i, j)
				found[shipLen]++
			}
		}
	}

	for length, reqCount := range required {
		if found[length] != reqCount {
			return false
		}
	}

	for length, count := range found {
		if required[length] == 0 && count > 0 {
			return false
		}
	}

	return true
}

func (game *Game) checkValidMove(x int, y int, boardNo bool) bool {
	if x > game.Board1.Size || y > game.Board1.Size {
		return false
	}
	return true
}

func (game *Game) PlayMove(x int, y int) (bool, error) {
	var ret bool
	var err error
	if game.Turn { // player1 turn
		ret, err = game.Board2.ShootCell(x, y)
	} else { // player2 turn
		ret, err = game.Board1.ShootCell(x, y)
	}
	return ret, err
}

func (game *Game) getBoard(player *model.Player, wg *sync.WaitGroup) {
	defer wg.Done()
	falseIfCompleted := true
	for falseIfCompleted {
		// ask player for board
		req := ws.WSMessage{
			Type:    "GetBoardMessage",
			Payload: ws.GetBoardMessage{},
		}
		if err := player.Conn.WriteJSON(req); err != nil {
			log.Printf("failed to send GetBoardMessage to %s: %v", player.Username, err)
			continue
		}

		// wait for response
		var msg ws.WSMessage
		if err := player.Conn.ReadJSON(&msg); err != nil {
			log.Printf("failed to read board from %s: %v", player.Username, err)
			sendErrorMessage("Failed to read board", player)
			continue
		}
		switch msg.Type {
		case "SendBoardMessage":
			var payload ws.SendBoardMessage

			raw, _ := json.Marshal(msg.Payload)
			if err := json.Unmarshal(raw, &payload); err != nil {
				log.Printf("invalid sendBoard from %s: %v", player.Username, err)
				sendErrorMessage("invalid SendBoardMessage", player)
				continue
			}
			var board = model.NewBoard(10)
			if (checkValidBoard(model.Board{Cells: payload.Cells})) {
				board.Cells = payload.Cells
				if game.Player1.ID == player.ID {
					game.Board1 = &board
				} else if game.Player2.ID == player.ID {
					game.Board2 = &board
				}
				log.Printf("Board received from %s", player.Username)
				falseIfCompleted = false
			} else {
				log.Printf("invalid sendBoard from %s", player.Username)
				sendErrorMessage("Invalid board", player)
				continue
			}

		case "ErrorMessage":
			var payload ws.ErrorMessage
			raw, _ := json.Marshal(msg.Payload)
			json.Unmarshal(raw, &payload)
			log.Printf("error from %s: %s", player.Username, payload.Error)
			continue
		}
	}
}

func sendErrorMessage(s string, player *model.Player) {
	errmsg := ws.WSMessage{
		Type: "ErrorMessage",
		Payload: ws.ErrorMessage{
			Error: s,
		},
	}

	if err := player.Conn.WriteJSON(errmsg); err != nil {
		log.Printf("failed to notify %s about error: %v", player.Username, err)
	}
}

func (game *Game) tellGameStarted() {
	msg := ws.WSMessage{
		Type: "gameStartedMessage",
		Payload: ws.GameStartedMessage{
			ID1:       game.Player1.ID,
			Username1: game.Player1.Username,
			ID2:       game.Player2.ID,
			Username2: game.Player2.Username,
			State:     game.State,
		},
	}

	if err := game.Player1.Conn.WriteJSON(msg); err != nil {
		log.Printf("failed to notify %s gameStarted: %v", game.Player1.Username, err)
	}
	if err := game.Player2.Conn.WriteJSON(msg); err != nil {
		log.Printf("failed to notify %s gameStarted: %v", game.Player2.Username, err)
	}
}

func (game *Game) getMove(hit bool, x int, y int) (int, int) {
	var current *model.Player
	if game.Turn {
		current = game.Player1
	} else {
		current = game.Player2
	}

	// tell current player it’s their turn, include result of opponent’s move
	notify := ws.WSMessage{
		Type: "GetTurnMessage",
		Payload: ws.GetTurnMessage{
			X:   x,
			Y:   y,
			Hit: hit,
		},
	}
	if err := current.Conn.WriteJSON(notify); err != nil {
		log.Printf("failed to notify %s about previous move: %v", current.Username, err)
		return -1, -1
	}

	// wait for their reply ( move)
	var msg ws.WSMessage
	if err := current.Conn.ReadJSON(&msg); err != nil {
		log.Printf("failed to read move from %s: %v", current.Username, err)
		return -1, -1
	}
	if msg.Type != "SendTurnMessage" {
		log.Printf("unexpected message type %s from %s", msg.Type, current.Username)
		return -1, -1
	}

	var payload ws.SendTurnMessage
	raw, _ := json.Marshal(msg.Payload)
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("invalid SendTurnMessage from %s: %v", current.Username, err)
		return -1, -1
	}

	return payload.X, payload.Y
}

func (game *Game) sendHitOrMiss(hit bool) {
	var current *model.Player
	if game.Turn {
		current = game.Player1
	} else {
		current = game.Player2
	}

	msg := ws.WSMessage{
		Type: "TurnResultMessage",
		Payload: ws.TurnResultMessage{
			Hit: hit,
		},
	}

	if err := current.Conn.WriteJSON(msg); err != nil {
		log.Printf("failed to send TurnResultMessage to %s: %v", current.Username, err)
	}
}

func (game *Game) tellResultsAskRematch() bool {
	// tell game result
	resultMsg := ws.WSMessage{
		Type: "GameResultMessage",
		Payload: ws.GameResultMessage{
			WinnerUname: game.winnerUsername,
		},
	}
	_ = game.Player1.Conn.WriteJSON(resultMsg)
	_ = game.Player2.Conn.WriteJSON(resultMsg)

	// make channels to collect responses
	respCh := make(chan bool, 2)

	waitForResponse := func(p *model.Player) {
		var msg ws.WSMessage
		if err := p.Conn.ReadJSON(&msg); err != nil {
			respCh <- false
			return
		}
		if msg.Type != "RematchMessage" {
			respCh <- false
			return
		}

		var payload ws.RematchMessage
		raw, _ := json.Marshal(msg.Payload)
		if err := json.Unmarshal(raw, &payload); err != nil {
			respCh <- false
			return
		}
		respCh <- payload.WantsRematch
	}

	go waitForResponse(game.Player1)
	go waitForResponse(game.Player2)

	timeout := time.After(60 * time.Second)
	yesCount := 0

	for i := 0; i < 2; i++ {
		select {
		case resp := <-respCh:
			if resp {
				yesCount++
			}
		case <-timeout:
			return false
		}
	}

	return yesCount == 2
}
