from __future__ import unicode_literals

import gzip
import io
import logging

import boto3
from tensorflow.python.framework import graph_util


BUCKET_NAME = 'chessbot'
OUTPUT_KEY_NAME = 'estimator/chessbot.pb.gz'


logger = logging.getLogger(__name__)


def export_model(session, model_dir):
    logger.info('Exporting model')

    output_node_names = ['f_p/output/BiasAdd']

    output_graph_def = graph_util.convert_variables_to_constants(
        session,
        session.graph.as_graph_def(),
        output_node_names,
    )

    export_model_to_s3(output_graph_def)

    logger.info('Export completed')


def export_model_to_s3(output_graph_def):
    s3 = boto3.resource('s3', region_name='eu-west-1')
    bucket = s3.Bucket(BUCKET_NAME)
    bucket.put_object(Key=OUTPUT_KEY_NAME, Body=output_graph_def.SerializeToString())
    logger.info('Model exported to s3://%s/%s', BUCKET_NAME, OUTPUT_KEY_NAME)
