#!/usr/bin/env bash

ORG=$1

# Add these thirdparty roles to an org

zed --insecure relationship create thirdparty_role:agent org organization:$ORG
zed --insecure relationship create thirdparty_role:auditor org organization:$ORG
zed --insecure relationship create thirdparty_role:loan_officer org organization:$ORG