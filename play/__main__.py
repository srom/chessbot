from __future__ import unicode_literals

import logging

import chess
import tensorflow as tf

from .load import load_graph, load_session
from .play import pick_next_move


logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
logger = logging.getLogger(__name__)


def main():
    board = chess.Board()

    graph = load_graph()
    with load_session(graph) as session:
        moves, other_equally_good_moves, score = pick_next_move(session, board)
        board.push(moves[0])

        print board
        print 'Score:', score
        print 'Num other good moves:', len(other_equally_good_moves)


if __name__ == '__main__':
    tf.logging.set_verbosity(tf.logging.INFO)
    main()
