from __future__ import unicode_literals

import logging
import time

import numpy as np

from .evaluate import evaluate


DEPTH = 5

logger = logging.getLogger(__name__)


def pick_next_move(session, board):
    start = time.time()

    best_moves, score = negamax(session, board, DEPTH, float('-inf'), float('inf'), 1)

    if not best_moves:
        raise ValueError('No valid move. Is the game over?')

    best_moves = shortest_moves(best_moves)

    # Randomly choose among best candidates
    index = np.random.randint(len(best_moves), size=1)[0]
    moves =  best_moves.pop(index)

    logger.info('pick_next_move took %.2f seconds', time.time() - start)

    return moves, best_moves, score


def negamax(session, board, depth, alpha, beta, color):
    game_over = is_game_over(board)
    if depth <= 0 or game_over:
        score = evaluate(session, board)
        if game_over:
            print 'Game Over:', score
            print board
            result_value = color * (float('inf') if board.is_checkmate() else 0)
            return [], result_value
        else:
            return [], color * score

    best_candidate_moves = []
    best_score = alpha

    for move in heuristic_move_order(board):
        next_board = board.copy()
        next_board.push(move)

        next_best_moves, score = negamax(session, next_board, depth - 1, -beta, -alpha, -color)

        score = -score
        if not best_candidate_moves or score >= best_score:
            if score > best_score:
                best_candidate_moves = []

            best_moves = [[move] + moves for moves in next_best_moves] or [[move]]
            best_candidate_moves = best_candidate_moves + best_moves
            best_score = score

        if score > alpha:
            alpha = score

        if alpha >= beta:
            break

    return best_candidate_moves, best_score


def shortest_moves(best_moves):
    shortest_moves = []
    shortest_path_length = None
    for moves in best_moves:
        l = len(moves)
        if shortest_path_length is None or l <= shortest_path_length:
            if l < shortest_path_length:
                shortest_moves = [moves]
            else:
                shortest_moves.append(moves)

            shortest_path_length = l

    return shortest_moves


def heuristic_move_order(board):
    return sorted(
        board.legal_moves,
        key=lambda move: (
            board.is_capture(move),
            board.is_into_check(move),
        ),
        reverse=True
    )


def is_game_over(board):
    return board.is_game_over(claim_draw=True)
