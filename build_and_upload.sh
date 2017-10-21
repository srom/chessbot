#!/bin/bash

TARFILE=chessbot.tar.gz
tar --exclude '.git*' --exclude '*.sqlite3' --exclude 'env' --exclude 'packages/*' --exclude '.idea/*' \
    --exclude '*.pyc' --exclude 'checkpoints/*' --exclude 'model/*' -pczvf $TARFILE .

aws s3 cp $TARFILE s3://romainstrock-packages/
