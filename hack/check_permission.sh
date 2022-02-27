#!/usr/bin/env bash

USER=$1
DEAL=$2
STATE=$3
ACTION=$4

if [ "$#" != "4" ]
then
  echo "Invalid number of arguments. Usage <user id> <deal id> <state> <action>"
  exit 1
fi

# check if user can do
zed --insecure permission check deal:"$DEAL"_"$STATE" $ACTION user:$USER