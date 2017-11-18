from __future__ import unicode_literals

import tensorflow as tf


WIDTH = HEIGHT = 8
CHANNELS = 12
FILTERS = 10
KERNEL_SIZE = 3
STRIDES = [1, 1]
PADDING = 'SAME'
DENSE_HIDDEN_UNITS = 1024

KAPPA = 10.0  # Emphasizes f(p) = -f(q)


class ChessConvolutionalNetwork(object):

    def __init__(self, learning_rate, adam_epsilon):
        with tf.variable_scope("input"):
            self.X = tf.placeholder(tf.float32, shape=(None, WIDTH, HEIGHT, CHANNELS), name='X')
            self.X_parent = tf.placeholder(tf.float32, shape=(None, WIDTH, HEIGHT, CHANNELS), name='X_parent')
            self.X_observed = tf.placeholder(tf.float32, shape=(None, WIDTH, HEIGHT, CHANNELS), name='X_observed')
            self.X_random = tf.placeholder(tf.float32, shape=(None, WIDTH, HEIGHT, CHANNELS), name='X_random')

        with tf.variable_scope("f_p"):
            self.training = tf.placeholder_with_default(False, shape=(), name='training')
            self.f = self._get_evaluation_function(self.X)
            tf.get_variable_scope().reuse_variables()
            self.f_parent = self._get_evaluation_function(self.X_parent)
            self.f_observed = self._get_evaluation_function(self.X_observed)
            self.f_random = self._get_evaluation_function(self.X_random)

        with tf.name_scope('loss'):
            self.loss = self._get_loss()

        with tf.name_scope('train'):
            self.training_op = self._get_training_op(learning_rate, adam_epsilon)

    def train(self, session, X_parent, X_observed, X_random):
        session.run(self.training_op, feed_dict={
            self.X_parent: X_parent,
            self.X_observed: X_observed,
            self.X_random: X_random,
            self.training: True,
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
        conv1 = self._get_convolutional_layer(X)
        conv2 = self._get_convolutional_layer(conv1)
        conv2_flat = self._reshape_conv_layer(conv2)
        dense = tf.layers.dense(conv2_flat, DENSE_HIDDEN_UNITS, activation=tf.nn.relu)
        output = tf.layers.dense(dense, 1, activation=None, name='output')
        return output

    def _get_convolutional_layer(self, input):
        return tf.layers.conv2d(
            input,
            filters=FILTERS,
            kernel_size=KERNEL_SIZE,
            strides=STRIDES,
            padding=PADDING,
            activation=tf.nn.relu
        )

    def _reshape_conv_layer(self, conv):
        _, width, height, channels = tf.shape(conv)
        return tf.reshape(conv, [-1, width * height * channels])

    def _get_loss(self):
        x_observed_random = self.f_random - self.f_observed
        x_parent_observed = self.f_parent + self.f_observed

        epsilon_log = 1e-3
        loss_a = -tf.log(epsilon_log + tf.sigmoid(x_observed_random))
        loss_b = -tf.log(epsilon_log + tf.sigmoid(KAPPA * x_parent_observed))
        loss_c = -tf.log(epsilon_log + tf.sigmoid(-KAPPA * x_parent_observed))

        return tf.reduce_mean(loss_a + loss_b + loss_c, name='loss')

    def _get_training_op(self, learning_rate, epsilon):
        optimizer = tf.train.AdamOptimizer(learning_rate=learning_rate, epsilon=epsilon)
        return optimizer.minimize(self.loss)
