#!/bin/bash

if [ -z $1 ]; then
  echo "Please supply the parties.json file e.g. '$0 parties.json'"
  exit 1
fi

SPACES=`jq 'keys[]' -r < $1`
echo "*** WARNING ***"
echo "This will delete the spaces (and ALL associated cloud functions and databases) for:"
echo "$SPACES"
echo 
echo "Note: Running the pipeline again will restore the configurations as how they are defined in $1 and the cloudFunctions/ directory in the repository."
echo "*** Press Enter to continue. ***"
read 

IFS=$'\n'
for space in $SPACES; do
  bluemix account space-delete $space -f
done
unset IFS
