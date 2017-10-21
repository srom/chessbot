from __future__ import unicode_literals

import chess

from common.parse_board import parse_board
from .load import load_model


def main():
    move_scores = []
    with load_model() as estimator:
        board = chess.Board()
        board.push_san('e4')
        for move in board.legal_moves:
            b = board.copy()
            b.push(move)
            x = parse_board(b)
            move_scores.append((b, move, estimator.evaluate(x)[0][0]))

        for board, move, score in sorted(move_scores, key=lambda t: t[2]):
            print move, score
            print board
            print


if __name__ == '__main__':
    main()
