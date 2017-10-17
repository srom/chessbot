from __future__ import unicode_literals

import argparse
import logging
import time

from .export import export_model
from .load import get_inputs_from_s3
from .train import train


logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
logger = logging.getLogger(__name__)


ITERATIONS_BEFORE_EXPORT = 5


def main(model_dir):
    classifier = None
    iteration = 0
    start = time.time()
    for inputs in get_inputs_from_s3():
        iteration += 1
        logger.info('Parse End: iteration %d; Elapsed: %fs', iteration, time.time() - start)
        classifier = train(classifier, inputs, model_dir)
        logger.info('Train End: iteration %d; Elapsed: %fs', iteration, time.time() - start)

        if iteration % ITERATIONS_BEFORE_EXPORT == 0:
            export_model(classifier, model_dir)

    export_model(classifier, model_dir)

    logger.info('DONE')


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('model_dir')
    args = parser.parse_args()
    main(args.model_dir)
