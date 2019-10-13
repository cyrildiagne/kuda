#!/bin/bash

source $KUDA_CMD_DIR/.config.sh

eksctl delete cluster \
--name=$KUDA_CLUSTER_NAME \
--region=$KUDA_CLUSTER_REGION \
--wait