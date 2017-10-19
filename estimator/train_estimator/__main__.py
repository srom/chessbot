from __future__ import unicode_literals

import argparse
import logging
import time

import tensorflow as tf

from .load import yield_batch
from .model import ChessDNNEstimator


logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
logger = logging.getLogger(__name__)


def main():
    start = time.time()
    estimator = ChessDNNEstimator()

    saver = tf.train.Saver()
    with tf.Session() as session:
        session.run(tf.global_variables_initializer())

        iteration = 0
        best_loss = float('inf')
        best_iteration = 0
        for X_train, X_test in yield_batch():
            iteration += 1
            X_p_train, X_o_train, X_r_train = X_train[:, 0, :], X_train[:, 1, :], X_train[:, 2, :]
            X_p_test, X_o_test, X_r_test = X_test[:, 0, :], X_test[:, 1, :], X_test[:, 2, :]

            estimator.train(session, X_p_train, X_o_train, X_r_train)
            loss_test = estimator.compute_loss(session, X_p_train, X_o_train, X_r_train)
            loss_rain = estimator.compute_loss(session, X_p_test, X_o_test, X_r_test)

            if loss_test < best_loss:
                best_loss = loss_test
                best_iteration = iteration
                saver.save(session, 'checkpoints/chessbot', global_step=iteration)

            elapsed = int(time.time() - start)
            logger.info('Training batch %d; Elapsed %ds; loss: %f (%f); best: %f (%d)',
                        iteration, elapsed, loss_test, loss_rain, best_loss, best_iteration)

    logger.info('DONE')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    # parser.add_argument('model_dir')
    args = parser.parse_args()
    main()
