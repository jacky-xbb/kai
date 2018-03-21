#!/bin/bash


# install CUDA Toolkit v9.1
CUDA_REPO_PKG="cuda-repo-ubuntu1604-9-1-local_9.1.85-1_amd64"
wget https://developer.nvidia.com/compute/cuda/9.1/Prod/local_installers/${CUDA_REPO_PKG}
sudo dpkg -i ${CUDA_REPO_PKG}
sudo apt-key add /var/cuda-repo-9-1-local/7fa2af80.pub
sudo apt-get update
sudo apt-get install --force-yes -y cuda 

# install cuDNN v6.0
CUDNN_TAR_FILE="cudnn-8.0-linux-x64-v6.0.tgz"
wget http://developer.download.nvidia.com/compute/redist/cudnn/v6.0/${CUDNN_TAR_FILE}
tar -xzvf ${CUDNN_TAR_FILE}
sudo cp -P cuda/include/cudnn.h /usr/local/cuda/include
sudo cp -P cuda/lib64/libcudnn* /usr/local/cuda/lib64/
sudo chmod a+r /usr/local/cuda/lib64/libcudnn*

# set environment variables
export PATH=/usr/local/cuda/bin${PATH:+:${PATH}}
export LD_LIBRARY_PATH=/usr/local/cuda/lib64\${LD_LIBRARY_PATH:+:${LD_LIBRARY_PATH}}