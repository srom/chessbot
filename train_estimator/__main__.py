from __future__ import unicode_literals

import argparse
import logging

import numpy as np
import tensorflow as tf

from .model import ChessDNNEstimator, INPUT_DIMENSION


logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
logger = logging.getLogger(__name__)


def main():
    estimator = ChessDNNEstimator()

    with tf.Session() as session:
        session.run(tf.global_variables_initializer())

        X_p_train = np.random.random_sample((1000, INPUT_DIMENSION))
        X_o_train = np.random.random_sample((1000, INPUT_DIMENSION))
        X_r_train = np.random.random_sample((1000, INPUT_DIMENSION))
        
        X_p_test = np.random.random_sample((10, INPUT_DIMENSION))
        X_o_test = np.random.random_sample((10, INPUT_DIMENSION))
        X_r_test = np.random.random_sample((10, INPUT_DIMENSION))

        estimator.train(session, X_p_train, X_o_train, X_r_train)

        loss = estimator.evaluate(session, X_p_test, X_o_test, X_r_test)
        print "Loss: {}".format(loss)

        instance = session.run(estimator.f, feed_dict={
            estimator.X: np.random.random_sample((1, INPUT_DIMENSION))
        })

        print 'Instance: {}'.format(instance)

    logger.info('DONE')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    # parser.add_argument('model_dir')
    args = parser.parse_args()
    main()
