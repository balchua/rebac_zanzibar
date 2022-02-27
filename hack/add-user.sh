#!/usr/bin/env bash

USER=$1
ORG=$2
THIRDPARTY_ROLE=$3
OPERATIONAL_ROLE=$4

if [ "$#" != "4" ]
then
  echo "Invalid number of arguments. Usage <userid to add> <organization> <thirdparty role> <user oper role>"
  exit 1
fi


# first add user to an organization
zed --insecure relationship create organization:$ORG member user:$USER

# now add user to the thirdparty role with his/her operational role
zed --insecure relationship create thirdparty_role:$THIRDPARTY_ROLE $OPERATIONAL_ROLE user:$USER