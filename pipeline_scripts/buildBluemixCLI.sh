#!/bin/bash -e

# Install Bluemix CLI

mkdir bluemix-cli
pushd bluemix-cli
  wget http://public.dhe.ibm.com/cloud/bluemix/cli/bluemix-cli/Bluemix_CLI_0.6.0_amd64.tar.gz

  tar xzf Bluemix_CLI_0.6.0_amd64.tar.gz

popd