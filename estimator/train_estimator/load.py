from __future__ import unicode_literals

import logging
import struct

import numpy as np

from common.triplets_pb2 import ChessBotTriplet


BATCH_SIZE = 1e3
TRAIN_RATIO = 0.8

logger = logging.getLogger(__name__)


def yield_batch(size=BATCH_SIZE):
    triplets = []
    logger.info('Loading batch')
    for triplet in yield_triplets():
        triplets.append(triplet)
        if len(triplets) == size:
            triplet_inputs = get_triplet_inputs(triplets)
            yield get_train_and_test_inputs(triplet_inputs)
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


def get_train_and_test_inputs(triplets):
    l = len(triplets)
    index = int(TRAIN_RATIO * l) - 1
    np.random.shuffle(triplets)
    return triplets[:index+1], triplets[index+1:]


def yield_triplets():
    # TODO: load from S3
    with open('/Users/srom/Downloads/1508291796.pb', 'rb') as f:
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
