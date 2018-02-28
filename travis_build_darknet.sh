#!/bin/bash

DARKNET="darknet"

# Build Darknet
git clone https://github.com/ZanLabs/darknet.git 
cd ${DARKNET}
# optionally GPU=1
make OPENCV=1 && make install
cd ..
rm -rf ${DARKNET} 