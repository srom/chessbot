from __future__ import unicode_literals

import argparse
import json
import logging
import os
import time

import tensorflow as tf

from ..load import yield_batch
from ..model import ChessDNNEstimator
from ..train import train_model, test_model


BATCH_SIZE = 1e3
TRAIN_TEST_RATIO = 0.8

LEARNING_RATES = [1e1, 1e0, 1e-1, 1e-2, 1e-3, 1e-4, 1e-5]
ADAM_EPSILONS = [1, 1e-1, 1e-2, 1e-3, 1e-4, 1e-5, 1e-6, 1e-7, 1e-8]


logger = logging.getLogger(__name__)


class GridSquare(object):
    __slots__ = ('learning_rate', 'epsilon', 'train_losses', 'test_losses')

    def __init__(self, learning_rate, epsilon):
        self.learning_rate = learning_rate
        self.epsilon = epsilon
        self.train_losses = []
        self.test_losses = []


def main(output_path):
    grid = []
    total = len(LEARNING_RATES) * len(ADAM_EPSILONS)
    grid_iteration = 0
    for initial_learning_rate in LEARNING_RATES:
        for epsilon in ADAM_EPSILONS:
            start = time.time()

            grid_iteration += 1
            logger.info(
                '%d / %d: LEARNING RATE %f; EPSILON %f',
                grid_iteration,
                total,
                initial_learning_rate,
                epsilon
            )

            estimator = ChessDNNEstimator(
                learning_rate=initial_learning_rate,
                adam_epsilon=epsilon,
            )
            square = GridSquare(initial_learning_rate, epsilon)

            with tf.Session() as session:
                session.run(tf.global_variables_initializer())

                iteration = 0
                best_loss = float('inf')
                best_iteration = 0
                for X_train, X_test in yield_batch(BATCH_SIZE, TRAIN_TEST_RATIO):
                    iteration += 1
                    loss_train = train_model(session, estimator, X_train)
                    loss_test = test_model(session, estimator, X_test)

                    square.train_losses.append(loss_train)
                    square.test_losses.append(loss_test)

                    if loss_test < best_loss:
                        best_loss = loss_test
                        best_iteration = iteration

                    elapsed = int(time.time() - start)
                    logger.info('Training batch %d; Elapsed %ds; loss: %f (train: %f); best: %f (%d)',
                                iteration, elapsed, loss_test, loss_train, best_loss, best_iteration)

            logger.info('Saving grid...')
            grid.append(square)
            save_grid(output_path, grid)
            logger.info('Saved to %s', output_path)

    logger.info('DONE')


def save_grid(output_path, grid):
    grid = dict(grid=grid)
    with open(output_path, 'w') as f:
        json.dump(grid, f)


if __name__ == '__main__':
    logger.info('Parameters Grid Search')
    logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
    dir_path = os.path.dirname(os.path.realpath(__file__))
    parser = argparse.ArgumentParser()
    parser.add_argument('--output', default=os.path.join(dir_path, 'grid.json'))
    args = parser.parse_args()
    output_path = args.output
    logger.info('Output file path: %s', output_path)
    main(output_path)
