from __future__ import unicode_literals

import logging
import struct

import boto3
import gzip
import io
import msgpack
import numpy as np


logger = logging.getLogger(__name__)


def get_inputs_from_s3():
    s3 = boto3.resource('s3', region_name='eu-west-1')
    bucket = s3.Bucket('chessbot')

    logger.info("Fetching input object summary items from S3")

    object_summary_items = list(bucket.objects.filter(Prefix='input'))

    logger.info('Inputs found: %d', len(object_summary_items))

    for index, object_summary in enumerate(object_summary_items):
        logger.info("Processing item %d: %s", index + 1, object_summary.key)
        object = object_summary.get()
        with io.BytesIO(object['Body'].read()) as bytestream:
            with gzip.GzipFile(fileobj=bytestream, mode='rb') as f:
                unpacker = msgpack.Unpacker(f)
                inputs = [struct.unpack("<769B", input) for input in unpacker]
                yield np.array(inputs, dtype=np.int32)
