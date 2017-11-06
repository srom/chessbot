from __future__ import unicode_literals

import argparse
import logging
import re


logger = logging.getLogger(__name__)


def main(log_path):
    for log_line in yield_train_log_line(log_path):
        print log_line
        break


class TrainLogLine(object):
    __slots__ = ('iteration', 'elapsed', 'test_loss', 'train_loss', 'best', 'best_iteration')

    def __init__(self, **kwargs):
        for key, value in kwargs.iteritems():
            setattr(self, key, value)

    def __unicode__(self):
        return (
            'Training batch {iteration}; ' +
            'Elapsed {elapsed}; ' +
            'loss: {test_loss} (train: {train_loss}); ' +
            'best: {best} ({best_iteration})'
        ).format(**self.to_dict())

    def __repr__(self):
        return self.__unicode__()

    def to_dict(self):
        return {
            key: getattr(self, key)
            for key in self.__slots__
            if hasattr(self, key)
        }


def yield_train_log_line(log_path):
    with open(log_path, 'r') as f:
        for line in f:
            if is_train_log_line(line):
                yield parse_line(line)


def parse_line(line):
    r = r'^.*Training batch ([0-9]+); Elapsed ([0-9]+)s; loss: ([0-9\.]+) \(train: ([0-9\.]+)\); best: ([0-9\.]+) \(([0-9]+)\).*$'
    m = re.match(r, line)

    if m is None:
        raise ValueError('No match for line {}'.format(line))

    return TrainLogLine(
        iteration=m.group(1),
        elapsed=m.group(2),
        test_loss=m.group(3),
        train_loss=m.group(4),
        best=m.group(5),
        best_iteration=m.group(6),
    )


def is_train_log_line(line):
    return re.search('Training batch', line) is not None


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO, format="%(asctime)s (%(levelname)s) %(message)s")
    parser = argparse.ArgumentParser()
    parser.add_argument('log_path')
    args = parser.parse_args()
    main(args.log_path)
