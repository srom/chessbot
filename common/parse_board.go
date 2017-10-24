package common

import (
	"github.com/notnil/chess"
)

const DIMENSION = 8*8*12

type BoardInput [1][DIMENSION]float32

var PIECES_CHESS = [12]chess.Piece{
	chess.BlackPawn,
	chess.BlackRook,
	chess.BlackKnight,
	chess.BlackBishop,
	chess.BlackQueen,
	chess.BlackKing,
	chess.WhitePawn,
	chess.WhiteRook,
	chess.WhiteKnight,
	chess.WhiteBishop,
	chess.WhiteQueen,
	chess.WhiteKing,
}

var POSITIONS_CHESS = [8 * 8]chess.Square{
	chess.A1,
	chess.B1,
	chess.C1,
	chess.D1,
	chess.E1,
	chess.F1,
	chess.G1,
	chess.H1,
	chess.A2,
	chess.B2,
	chess.C2,
	chess.D2,
	chess.E2,
	chess.F2,
	chess.G2,
	chess.H2,
	chess.A3,
	chess.B3,
	chess.C3,
	chess.D3,
	chess.E3,
	chess.F3,
	chess.G3,
	chess.H3,
	chess.A4,
	chess.B4,
	chess.C4,
	chess.D4,
	chess.E4,
	chess.F4,
	chess.G4,
	chess.H4,
	chess.A5,
	chess.B5,
	chess.C5,
	chess.D5,
	chess.E5,
	chess.F5,
	chess.G5,
	chess.H5,
	chess.A6,
	chess.B6,
	chess.C6,
	chess.D6,
	chess.E6,
	chess.F6,
	chess.G6,
	chess.H6,
	chess.A7,
	chess.B7,
	chess.C7,
	chess.D7,
	chess.E7,
	chess.F7,
	chess.G7,
	chess.H7,
	chess.A8,
	chess.B8,
	chess.C8,
	chess.D8,
	chess.E8,
	chess.F8,
	chess.G8,
	chess.H8,
}

func ParseBoard(board *chess.Board) BoardInput {
	input := [DIMENSION]float32{}
	index := 0
	for _, piece := range PIECES_CHESS {
		for _, position := range POSITIONS_CHESS {
			if board.Piece(position) == piece {
				input[index] = 1.0
			} else {
				input[index] = 0.0
			}
			index += 1
		}
	}
	return BoardInput([1][DIMENSION]float32{input})
}
