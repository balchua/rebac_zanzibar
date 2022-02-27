#!/usr/bin/env bash

for i in {1..100000}
do
   ./create_deal.sh $i singapore agent
done
