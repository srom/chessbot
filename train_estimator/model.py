from __future__ import unicode_literals

import tensorflow as tf


INPUT_DIMENSION = 8*8*12  # 768 = 8 x 8 squares x 12 piece types
HIDDEN_UNITS = 2048
KAPPA = 10.0  # Gives emphasize to f(p) = -f(q)
LEARNING_RATE = 0.03
MOMENTUM = 0.9


class ChessDNNEstimator(object):

    def __init__(self):
        self.X = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X')
        self.X_parent = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X_parent')
        self.X_observed = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X_observed')
        self.X_random = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X_random')

        with tf.variable_scope("evaluation_function_scope"):
            self.f = self._get_evaluation_function(self.X)
            tf.get_variable_scope().reuse_variables()
            self.f_parent = self._get_evaluation_function(self.X_parent)
            self.f_observed = self._get_evaluation_function(self.X_observed)
            self.f_random = self._get_evaluation_function(self.X_random)

        self.loss = self._get_loss()
        self.training_op = self._get_training_op()

    def train(self, session, X_parent, X_observed, X_random):
        # TODO: add mini batch
        session.run(self.training_op, feed_dict={
            self.X_parent: X_parent,
            self.X_observed: X_observed,
            self.X_random: X_random,
        })

    def evaluate(self, session, X_parent, X_observed, X_random):
        return session.run(self.loss, feed_dict={
            self.X_parent: X_parent,
            self.X_observed: X_observed,
            self.X_random: X_random,
        })

    def _get_evaluation_function(self, X):
        with tf.name_scope('dnn'):
            hidden_1 = tf.layers.dense(X, HIDDEN_UNITS, activation=tf.nn.relu, name='hidden_1')
            hidden_2 = tf.layers.dense(hidden_1, HIDDEN_UNITS, activation=tf.nn.relu, name='hidden_2')
            output = tf.layers.dense(hidden_2, 1, activation=None, name='output')
            return output


    def _get_loss(self):
        with tf.name_scope('loss'):
            x_observed_random = self.f_observed - self.f_random
            x_parent_observed = KAPPA * (self.f_parent + self.f_observed)

            loss_a = -tf.log(tf.sigmoid(x_observed_random))
            loss_b = -tf.log(tf.sigmoid(x_parent_observed))
            loss_c = -tf.log(tf.sigmoid(-x_parent_observed))

            return tf.reduce_mean(loss_a + loss_b + loss_c, name='loss')

    def _get_training_op(self):
        with tf.name_scope('train'):
            optimizer = tf.train.MomentumOptimizer(learning_rate=LEARNING_RATE, momentum=MOMENTUM)
            return optimizer.minimize(self.loss)
