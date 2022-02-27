#!/usr/bin/env bash

USER=$1
THIRDPARTY_ROLE=$2

if [ "$#" != "2" ]
then
  echo "Invalid number of arguments. Usage <user id> <thirdparty role>"
  exit 1
fi
# check if user can do
zed --insecure permission check thirdparty_role:$THIRDPARTY_ROLE create_deal user:$USER