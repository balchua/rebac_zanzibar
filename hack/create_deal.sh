#!/usr/bin/env bash

DEAL=$1
ORG=$2
THIRDPARTY_ROLE=$3

if [ "$#" != "3" ]
then
  echo "Invalid number of arguments. Usage <deal id> <organization> <thirdparty role>"
  exit 1
fi

# Assign this deal to an organization
echo "attaching deal "$DEAL"_created to organization $ORG"
echo "zed --insecure relationship create deal:"$DEAL"_created org organization:$ORG"
zed --insecure relationship create deal:"$DEAL"_created org organization:$ORG 

# assign the deal to a thirdparty role, the deal_id must be in this format <id>_<state>
echo "zed --insecure relationship create deal:"$DEAL"_created thirdparty thirdparty_role:$THIRDPARTY_ROLE"
zed --insecure relationship create deal:"$DEAL"_created thirdparty thirdparty_role:$THIRDPARTY_ROLE 
