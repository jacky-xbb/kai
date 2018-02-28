#!/bin/bash

DARKNET="darknet"

# Build Darknet
git clone https://github.com/ZanLabs/darknet.git 
cd ${DARKNET}
# optionally GPU=1
sudo make OPENCV=1 && sudo make install
cd ..
rm -rf ${DARKNET} 