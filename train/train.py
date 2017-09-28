from __future__ import unicode_literals

import logging

import numpy as np
import tensorflow as tf
from tensorflow.python.framework import dtypes

FEATURES_DIMENSION = 768  # 8 * 8 squares * 12 piece types
RATIO_TRAIN_VS_TEST = 0.8

logger = logging.getLogger(__name__)


def train(classifier, inputs, model_dir):
    train_inputs, test_inputs = get_train_and_test_features(inputs)

    train_input_fn = tf.estimator.inputs.numpy_input_fn(
        x={"x": train_inputs[:, :-1]},
        y=train_inputs[:, -1:].flatten(),
        num_epochs=1,
        shuffle=False
    )

    test_input_fn = tf.estimator.inputs.numpy_input_fn(
        x={"x": test_inputs[:, :-1]},
        y=test_inputs[:, -1:].flatten(),
        num_epochs=1,
        shuffle=False
    )

    feature_columns = [tf.feature_column.numeric_column("x", shape=[FEATURES_DIMENSION], dtype=dtypes.int32)]

    if classifier is None:
        classifier = tf.estimator.DNNClassifier(
            feature_columns=feature_columns,
            hidden_units=[1024, 1024],
            n_classes=3,
            model_dir=model_dir
        )

    logger.info('Training...')

    classifier = classifier.train(input_fn=train_input_fn)

    logger.info('Training Completed')

    accuracy_score = classifier.evaluate(input_fn=test_input_fn)["accuracy"]

    logger.info('Accuracy score: %f', accuracy_score)

    return classifier


def get_train_and_test_features(inputs):
    l = len(inputs)
    index = int(RATIO_TRAIN_VS_TEST * l)
    np.random.shuffle(inputs)
    return inputs[:index+1], inputs[index+1:]
