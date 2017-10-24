package common

import "github.com/notnil/chess"

func CopyGame(game *chess.Game) (*chess.Game, error) {
	fen, err := chess.FEN(game.FEN())
	if err != nil {
		return nil, err
	}
	return chess.NewGame(fen), nil
}
