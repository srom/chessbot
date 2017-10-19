from __future__ import unicode_literals

import tensorflow as tf


INPUT_DIMENSION = 768  # 8 x 8 squares x 12 piece types
HIDDEN_UNITS = 2048
KAPPA = 1.0  # Emphasizes f(p) = -f(q)
INITIAL_LEARNING_RATE = 0.1
MOMENTUM = 0.9
DECAY_STEPS = 1e4
DECAY_RATE = 0.1


class ChessDNNEstimator(object):

    def __init__(self):
        self.X = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X')
        self.X_parent = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X_parent')
        self.X_observed = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X_observed')
        self.X_random = tf.placeholder(tf.float32, shape=(None, INPUT_DIMENSION), name='X_random')

        with tf.variable_scope("dnn"):
            self.f = self._get_evaluation_function(self.X)
            tf.get_variable_scope().reuse_variables()
            self.f_parent = self._get_evaluation_function(self.X_parent)
            self.f_observed = self._get_evaluation_function(self.X_observed)
            self.f_random = self._get_evaluation_function(self.X_random)

        with tf.name_scope('loss'):
            self.loss = self._get_loss()

        with tf.name_scope('train'):
            self.training_op = self._get_training_op()

    def train(self, session, X_parent, X_observed, X_random):
        session.run(self.training_op, feed_dict={
            self.X_parent: X_parent,
            self.X_observed: X_observed,
            self.X_random: X_random,
        })

    def compute_loss(self, session, X_parent, X_observed, X_random):
        return session.run(self.loss, feed_dict={
            self.X_parent: X_parent,
            self.X_observed: X_observed,
            self.X_random: X_random,
        })

    def evaluate(self, session, X):
        return session.run(self.f, feed_dict={
            self.X: X
        })

    def _get_evaluation_function(self, X):
        hidden_1 = tf.layers.dense(X, HIDDEN_UNITS, activation=tf.nn.relu, name='hidden_1')
        hidden_2 = tf.layers.dense(hidden_1, HIDDEN_UNITS, activation=tf.nn.relu, name='hidden_2')
        output = tf.layers.dense(hidden_2, 1, activation=None, name='output')
        return output

    def _get_loss(self):
        x_observed_random = self.f_observed - self.f_random
        x_parent_observed = self.f_parent + self.f_observed

        loss_a = tf.log(1 + tf.sigmoid(x_observed_random))
        loss_b = KAPPA * tf.log(1 + tf.sigmoid(x_parent_observed))
        loss_c = KAPPA * tf.log(1 + tf.sigmoid(-x_parent_observed))

        return tf.reduce_mean(loss_a + loss_b + loss_c, name='loss')

    def _get_training_op(self):
        global_step = tf.Variable(0, trainable=False, name='global_step')
        learning_rate = tf.train.exponential_decay(INITIAL_LEARNING_RATE, global_step, DECAY_STEPS, DECAY_RATE)
        optimizer = tf.train.MomentumOptimizer(learning_rate=learning_rate, momentum=MOMENTUM)
        return optimizer.minimize(self.loss, global_step=global_step)
