from __future__ import unicode_literals

import chess

from common.parse_board import parse_board
from .load import load_model


def main():
    board = chess.Board()

    with load_model() as estimator:
        for move in board.legal_moves:
            b = board.copy()
            b.push(move)
            x = parse_board(b)

            print 'Move {}: {}'.format(move, estimator.evaluate(x)['class_ids'][0][0])


if __name__ == '__main__':
    main()
