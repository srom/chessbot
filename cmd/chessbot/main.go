package main

import (
	"fmt"
	"log"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/play"
)

func main() {
	model, err := play.LoadModel("/Users/srom/workspace/go/src/github.com/srom/chessbot/model/chessbot.pb")
	if err != nil {
		log.Fatalf("Error parsing model: %v", err)
	}
	defer model.Close()

	for i := 0; i < 20; i++ {
		game := chess.NewGame()
		moves := game.ValidMoves()
		move := moves[i]
		game.Move(move)
		boardInput := play.ParseBoard(game.Position().Board())
		class, err := model.Evaluate(boardInput)
		if err != nil {
			log.Fatalf("Error evaluation model: %v", err)
		}
		fmt.Print(game.Position().Board().Draw())
		log.Printf("%v: %v", move, class)
	}
}
