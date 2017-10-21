from __future__ import unicode_literals

import logging

import boto3
import tensorflow as tf
from tensorflow.python.estimator import model_fn as model_fn_lib
from tensorflow.python.framework import graph_util
from tensorflow.python.framework import ops

from .train import FEATURES_DIMENSION


BUCKET_NAME = 'chessbot'
OUTPUT_KEY_NAME = 'model/chessbot.pb'


logger = logging.getLogger(__name__)


def export_model(classifier, model_dir):
    logger.info('Exporting model')

    checkpoint_path = tf.train.latest_checkpoint(model_dir)

    predict_input_fn = lambda: {
        'x': tf.placeholder(tf.int32, shape=(None, FEATURES_DIMENSION), name='prediction_inputs')
    }

    with ops.Graph().as_default() as graph:
        classifier._create_and_assert_global_step(graph)

        features = classifier._get_features_from_input_fn(predict_input_fn, model_fn_lib.ModeKeys.PREDICT)
        classifier._call_model_fn(features, None, model_fn_lib.ModeKeys.PREDICT)

        saver = tf.train.Saver()
        with tf.Session(graph=graph) as session:
            saver.restore(session, checkpoint_path)

            output_node_names = ['dnn/head/predictions/class_ids']

            output_graph_def = graph_util.convert_variables_to_constants(
                session,
                graph.as_graph_def(),
                output_node_names,
            )

            export_model_to_s3(output_graph_def)

    logger.info('Export completed')


def export_model_to_s3(output_graph_def):
    s3 = boto3.resource('s3', region_name='eu-west-1')
    bucket = s3.Bucket(BUCKET_NAME)
    bucket.put_object(Key=OUTPUT_KEY_NAME, Body=output_graph_def.SerializeToString())
    logger.info('Model exported to s3://%s/%s', BUCKET_NAME, OUTPUT_KEY_NAME)
