from __future__ import unicode_literals

import tensorflow as tf


MODEL_PREFIX = "chessbot"


def load_graph():
    with tf.gfile.GFile("chessbot.pb", "rb") as f:
        graph_def = tf.GraphDef()
        graph_def.ParseFromString(f.read())

    with tf.Graph().as_default() as graph:
        tf.import_graph_def(
            graph_def,
            input_map=None,
            return_elements=None,
            name=MODEL_PREFIX,
            op_dict=None,
            producer_op_list=None
        )
        return graph


def load_session(graph):
    session = tf.Session(graph=graph)
    session.x = graph.get_tensor_by_name('chessbot/prediction_inputs:0')
    session.y = dict(class_ids=graph.get_tensor_by_name('chessbot/dnn/head/predictions/class_ids:0'))
    return session
