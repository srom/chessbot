package main

import (
	"fmt"

	"gopkg.in/freeeve/pgn.v1"
)

const DIMENSION = 8*8*12 + 1

var PIECES = [12]pgn.Piece{
	pgn.BlackPawn,
	pgn.BlackRook,
	pgn.BlackKnight,
	pgn.BlackBishop,
	pgn.BlackQueen,
	pgn.BlackKing,
	pgn.WhitePawn,
	pgn.WhiteRook,
	pgn.WhiteKnight,
	pgn.WhiteBishop,
	pgn.WhiteQueen,
	pgn.WhiteKing,
}
var POSITIONS = [8 * 8]pgn.Position{
	pgn.A1,
	pgn.B1,
	pgn.C1,
	pgn.D1,
	pgn.E1,
	pgn.F1,
	pgn.G1,
	pgn.H1,
	pgn.A2,
	pgn.B2,
	pgn.C2,
	pgn.D2,
	pgn.E2,
	pgn.F2,
	pgn.G2,
	pgn.H2,
	pgn.A3,
	pgn.B3,
	pgn.C3,
	pgn.D3,
	pgn.E3,
	pgn.F3,
	pgn.G3,
	pgn.H3,
	pgn.A4,
	pgn.B4,
	pgn.C4,
	pgn.D4,
	pgn.E4,
	pgn.F4,
	pgn.G4,
	pgn.H4,
	pgn.A5,
	pgn.B5,
	pgn.C5,
	pgn.D5,
	pgn.E5,
	pgn.F5,
	pgn.G5,
	pgn.H5,
	pgn.A6,
	pgn.B6,
	pgn.C6,
	pgn.D6,
	pgn.E6,
	pgn.F6,
	pgn.G6,
	pgn.H6,
	pgn.A7,
	pgn.B7,
	pgn.C7,
	pgn.D7,
	pgn.E7,
	pgn.F7,
	pgn.G7,
	pgn.H7,
	pgn.A8,
	pgn.B8,
	pgn.C8,
	pgn.D8,
	pgn.E8,
	pgn.F8,
	pgn.G8,
	pgn.H8,
}

type BoardFeaturesAndResult [DIMENSION]uint8

func parseBoard(board *pgn.Board, result uint8) *BoardFeaturesAndResult {
	input := BoardFeaturesAndResult{}
	index := 0
	for _, piece := range PIECES {
		for _, position := range POSITIONS {
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
