from __future__ import unicode_literals

import chess
import numpy as np


PIECES = [
    (chess.PAWN, chess.BLACK),
    (chess.ROOK, chess.BLACK),
    (chess.KNIGHT, chess.BLACK),
    (chess.BISHOP, chess.BLACK),
    (chess.QUEEN, chess.BLACK),
    (chess.KING, chess.BLACK),
    (chess.PAWN, chess.WHITE),
    (chess.ROOK, chess.WHITE),
    (chess.KNIGHT, chess.WHITE),
    (chess.BISHOP, chess.WHITE),
    (chess.QUEEN, chess.WHITE),
    (chess.KING, chess.WHITE),
]


def parse_board(board):
    feature = []
    for piece_type, colour in PIECES:
        square_set = board.pieces(piece_type, colour)
        for square in chess.SQUARES:
            mask = chess.BB_SQUARES[square]
            if mask & square_set.mask:
                feature.append(1)
            else:
                feature.append(0)

    return np.array(feature, dtype=np.int32)
