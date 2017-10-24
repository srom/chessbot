from __future__ import unicode_literals

import gzip
import io
import logging
import struct

import boto3
import numpy as np

from common.triplets_pb2 import ChessBotTriplet


logger = logging.getLogger(__name__)


def yield_batch(batch_size, train_test_ratio):
    triplets = []
    logger.info('Loading batch')
    for triplet in yield_triplets():
        triplets.append(triplet)
        if len(triplets) == batch_size:
            triplet_inputs = get_triplet_inputs(triplets)
            yield get_train_and_test_inputs(triplet_inputs, train_test_ratio)
            triplets = []


def get_triplet_inputs(triplets):
    items = []
    for triplet in triplets:
        items.append([
            triplet.parent.pieces,
            triplet.observed.pieces,
            triplet.random.pieces,
        ])
    return np.array(items, dtype=np.float32)


def get_train_and_test_inputs(triplets, train_test_ratio):
    l = len(triplets)
    index = int(train_test_ratio * l) - 1
    np.random.shuffle(triplets)
    return triplets[:index+1], triplets[index+1:]


def yield_triplets():
    for f in yield_triplet_files_from_s3():
        while True:
            sizeBytes = f.read(4)
            if not sizeBytes or len(sizeBytes) < 4:
                break

            next_message_size = struct.unpack("I", sizeBytes)[0]

            message_bytes = f.read(next_message_size)

            if len(message_bytes) < next_message_size:
                logger.error('Truncated message: expected len %d but got %d', next_message_size, len(message_bytes))
                break

            triplet = ChessBotTriplet()
            triplet.ParseFromString(message_bytes)

            yield triplet


def yield_triplet_files_from_s3():
    s3 = boto3.resource('s3', region_name='eu-west-1')
    bucket = s3.Bucket('chessbot')

    object_summary_items = list(bucket.objects.filter(Prefix='triplets_all'))

    logger.info('Triplet files found: %d', len(object_summary_items))

    for index, object_summary in enumerate(object_summary_items):
        logger.info("Processing triplet file %d: %s", index + 1, object_summary.key)
        object = object_summary.get()
        with io.BytesIO(object['Body'].read()) as bytestream:
            with gzip.GzipFile(fileobj=bytestream, mode='rb') as f:
                yield f
