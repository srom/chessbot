from __future__ import unicode_literals

import argparse
import logging
import os
import time

import tensorflow as tf

from .export import export_model
from .load import yield_batch
from .model import ChessDNNEstimator


BATCH_SIZE = 1e3
TRAIN_TEST_RATIO = 0.8
ITERATIONS_BETWEEN_EXPORTS = 10


logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
logger = logging.getLogger(__name__)


def main(model_dir, should_export):
    start = time.time()
    save_path = os.path.join(model_dir, 'chessbot')

    estimator = ChessDNNEstimator()

    saver = tf.train.Saver()
    with tf.Session() as session:
        session.run(tf.global_variables_initializer())

        iteration = 0
        best_loss = float('inf')
        best_iteration = 0
        last_exported_iteration = 0
        for X_train, X_test in yield_batch(BATCH_SIZE, TRAIN_TEST_RATIO):
            iteration += 1
            X_p_train, X_o_train, X_r_train = X_train[:, 0, :], X_train[:, 1, :], X_train[:, 2, :]
            X_p_test, X_o_test, X_r_test = X_test[:, 0, :], X_test[:, 1, :], X_test[:, 2, :]

            estimator.train(session, X_p_train, X_o_train, X_r_train)
            loss_test = estimator.compute_loss(session, X_p_test, X_o_test, X_r_test)
            loss_train = estimator.compute_loss(session, X_p_train, X_o_train, X_r_train)

            if loss_test < best_loss:
                best_loss = loss_test
                best_iteration = iteration
                saver.save(session, save_path, global_step=iteration)

            elapsed = int(time.time() - start)
            logger.info('Training batch %d; Elapsed %ds; loss: %f (train: %f); best: %f (%d)',
                        iteration, elapsed, loss_test, loss_train, best_loss, best_iteration)

            if ready_to_export(should_export, iteration, last_exported_iteration, best_iteration):
                export_model(saver, model_dir)
                last_exported_iteration = best_iteration

    logger.info('DONE')


def ready_to_export(should_export, iteration, last_exported_iteration, best_iteration):
    return should_export and iteration % ITERATIONS_BETWEEN_EXPORTS == 0 and best_iteration > last_exported_iteration


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--model_dir', default='checkpoints')
    parser.add_argument('--no_export', action='store_true')
    args = parser.parse_args()
    should_export = not args.no_export
    main(args.model_dir, should_export)
