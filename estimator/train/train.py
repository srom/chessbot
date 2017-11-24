from __future__ import unicode_literals


def train_model(session, estimator, X_train):
    X_p_train, X_o_train, X_r_train = X_train[:, 0, :], X_train[:, 1, :], X_train[:, 2, :]
    estimator.train(session, X_p_train, X_o_train, X_r_train)
    return estimator.compute_loss(session, X_p_train, X_o_train, X_r_train)


def test_model(session, estimator, X_test, detailed=False):
    X_p_test, X_o_test, X_r_test = X_test[:, 0, :], X_test[:, 1, :], X_test[:, 2, :]
    return estimator.compute_loss(session, X_p_test, X_o_test, X_r_test, detailed=detailed)
