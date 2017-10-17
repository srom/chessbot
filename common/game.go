package common

import "github.com/notnil/chess"

func CopyGame(game *chess.Game) *chess.Game {
	moves := game.Moves()
	return chess.NewGame(moveGameForward(moves))
}

func moveGameForward(moves []*chess.Move) func(*chess.Game) {
	return func(g *chess.Game) {
		for _, move := range moves {
			g.Move(move)
		}
	}
}
