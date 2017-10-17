package preprocess

import (
	"fmt"

	"gopkg.in/freeeve/pgn.v1"
	"github.com/srom/chessbot/common"
)

//go:generate msgp -o=input.go -marshal=false

const DIMENSION = 8*8*12 + 1

type BoardFeaturesAndResult [DIMENSION]uint8

func parseBoard(board *pgn.Board, result uint8) *BoardFeaturesAndResult {
	input := BoardFeaturesAndResult{}
	index := 0
	for _, piece := range common.PIECES {
		for _, position := range common.POSITIONS {
			if board.GetPiece(position) == piece {
				input[index] = 1
			} else {
				input[index] = 0
			}
			index += 1
		}
	}
	input[DIMENSION-1] = result
	return &input
}

func parseResult(res string) (uint8, error) {
	if res == "1-0" {
		return 2, nil
	} else if res == "1/2-1/2" {
		return 1, nil
	} else if res == "0-1" {
		return 0, nil
	} else {
		return 255, fmt.Errorf("Unknown result %s", res)
	}
}
