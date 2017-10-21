from __future__ import unicode_literals

import numpy as np
from common.parse_board import parse_board


WIN = 1
DRAW = 0
LOSE = -1


def evaluate(session, next_board):
    move_input = parse_board(next_board)

    feature = np.array([move_input], dtype=np.int32)
    y_classes = session.run(session.y, feed_dict={session.x: feature})
    y = y_classes['class_ids'][0]

    if y == 0:
        return LOSE
    elif y == 1:
        return DRAW
    elif y == 2:
        return WIN
    else:
        raise ValueError('Invalid class %d', y)
