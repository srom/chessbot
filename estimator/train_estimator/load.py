from __future__ import unicode_literals

import logging

import numpy as np

from common.triplets_pb2 import ChessBotTriplets


TRAIN_RATIO = 0.8

logger = logging.getLogger(__name__)


def yield_batch(size=1e2):
    triplets = []
    logger.info('Loading batch')
    for triplet in yield_triplets():
        triplets.append(triplet)
        logger.info('%d', len(triplets))
        if len(triplets) == size:
            triplet_inputs = get_triplet_inputs(triplets)
            logger.info('Yielding batch of size %d', size)
            yield get_train_and_test_inputs(triplet_inputs)
            triplets = []


def get_triplet_inputs(triplets):
    triplet_inputs = []
    for triplet in triplets:
        triplet_inputs.append(np.array([
            triplet.parent.pieces,
            triplet.observed.pieces,
            triplet.random.pieces,
        ], dtype=np.float32))
    return np.array(triplet_inputs, dtype=np.float32)


def get_train_and_test_inputs(triplets):
    l = len(triplets)
    index = int(TRAIN_RATIO * l)
    np.random.shuffle(triplets)
    return triplets[:index+1], triplets[index+1:]


def yield_triplets():
    # TODO: load from S3
    logger.info('Yield triplets')
    with open('/Users/srom/Downloads/1508281722.pb', 'rb') as f:
        triplets = ChessBotTriplets()
        triplets.ParseFromString(f.read())
        logger.info('hello?')

    for triplet in triplets.triplets:
        yield triplet
