#!/bin/bash
# WARNING: FOR PIPELINE USE UNLESS YOU WANT TO INSTALL GO
echo "For pipeline use!"

# sudo apt-get -y update
# sudo apt-get -y upgrade

mkdir goinstall

pushd goinstall
  wget https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz

  tar xzf go1.9.linux-amd64.tar.gz

  mkdir workspace

popd
