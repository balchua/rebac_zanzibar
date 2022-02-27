#!/usr/bin/env bash

USER=$1
THIRDPARTY_ROLE=$2

# check if user can do
zed --insecure permission check thirdparty_role:$THIRDPARTY_ROLE create_deal user:$USER