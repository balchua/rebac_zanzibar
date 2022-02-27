#!/usr/bin/env bash

ORG=$1

if [ "$#" != "1" ]
then
  echo "Invalid number of arguments. Usage <organization>"
  exit 1
fi
# Add these thirdparty roles to an org

zed --insecure relationship create thirdparty_role:agent org organization:$ORG
zed --insecure relationship create thirdparty_role:auditor org organization:$ORG
zed --insecure relationship create thirdparty_role:loan_officer org organization:$ORG