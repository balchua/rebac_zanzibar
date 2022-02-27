#!/usr/bin/env bash

DEAL=$1
NEXT_STATE=$2
ORG=$3

if [ "$#" != "3" ]
then
  echo "Invalid number of arguments. Usage <deal id> <next state> <organization>"
  exit 1
fi
# check state

if [ "$NEXT_STATE" == "reviewed" ]
then
    # delete old relationships
    zed --insecure relationship bulk-delete deal:"$DEAL"_created --force

    # Assign this deal to an organization
    echo "attaching deal "$DEAL"_reviewed to organization $ORG"
    echo "zed --insecure relationship create deal:"$DEAL"_reviewed org organization:$ORG"
    zed --insecure relationship create deal:"$DEAL"_reviewed org organization:$ORG 

    # assign the deal to a thirdparty role, the deal_id must be in this format <id>_<state>
    echo "zed --insecure relationship create deal:"$DEAL"_reviewed thirdparty thirdparty_role:auditor"
    zed --insecure relationship create deal:"$DEAL"_reviewed thirdparty thirdparty_role:auditor
fi

if [ "$NEXT_STATE" == "validated" ]
then
    # delete old relationships
    zed --insecure relationship bulk-delete deal:"$DEAL"_reviewed --force
    # Assign this deal to an organization
    echo "attaching deal "$DEAL"_validated to organization $ORG"
    echo "zed --insecure relationship create deal:"$DEAL"_validated org organization:$ORG"
    zed --insecure relationship create deal:"$DEAL"_validated org organization:$ORG 

    # assign the deal to a thirdparty role, the deal_id must be in this format <id>_<state>
    echo "zed --insecure relationship create deal:"$DEAL"_validated thirdparty thirdparty_role:loan_officer"
    zed --insecure relationship create deal:"$DEAL"_validated thirdparty thirdparty_role:loan_officer 
fi

if [ "$NEXT_STATE" == "processed" ]
then
    # delete old relationships
    zed --insecure relationship bulk-delete deal:"$DEAL"_validated --force
    # Assign this deal to an organization
    echo "attaching deal "$DEAL"_processed to organization $ORG"
    echo "zed --insecure relationship create deal:"$DEAL"_processed org organization:$ORG"
    zed --insecure relationship create deal:"$DEAL"_processed org organization:$ORG 

    # assign the deal to a thirdparty role, the deal_id must be in this format <id>_<state>
    echo "zed --insecure relationship create deal:"$DEAL"_processed thirdparty thirdparty_role:processed"
    zed --insecure relationship create deal:"$DEAL"_processed thirdparty thirdparty_role:loan_officer 
    zed --insecure relationship create deal:"$DEAL"_processed thirdparty thirdparty_role:agent
    zed --insecure relationship create deal:"$DEAL"_processed thirdparty thirdparty_role:auditor
fi