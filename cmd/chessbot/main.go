package main

import (
	"fmt"
	"log"
	"time"

	"github.com/notnil/chess"
	"github.com/srom/chessbot/play"
	"path/filepath"
)

const DEPTH = 3

func main() {
	modelPath, err := filepath.Abs("../../model/chessbot.pb")
	if err != nil {
		log.Fatalf("Error reading path: %v", err)
	}
	model, err := play.LoadModel(modelPath)
	if err != nil {
		log.Fatalf("Error parsing model: %v", err)
	}
	defer model.Close()

	game := chess.NewGame()

	player := int64(-11)
	moveUnit := 0
	for game.Outcome() == chess.NoOutcome {
		moveUnit += 1
		start := time.Now()
		player = -1 * player
		moveNodes, _ := play.Negamax(model, game, DEPTH, -1 * play.MAX_SCORE, play.MAX_SCORE, player)
		randomMoveNode := play.PickRandomMove(moveNodes)
		game.Move(randomMoveNode.Move)

		moves := []string{}
		for _, moveNode := range moveNodes {
			moves = append(moves, moveNode.Move.String())
		}

		log.Printf("Moves: %d %v %v", len(moveNodes), moves, randomMoveNode.Move.String())
		if moveUnit % 2 == 0 {
			log.Printf("Round %d", moveUnit / 2)
			fmt.Print(game.Position().Board().Draw())
			log.Printf("Elapsed: %v", time.Since(start))
			fmt.Print("\n")
		}
	}

	log.Printf("Outcome: %v", game.Outcome())
	fmt.Print(game.Position().Board().Draw())
	fmt.Print("\n")
}
