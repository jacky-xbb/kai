FROM yummybian/docker-darknet:latest

ENV KAI-REPO kai-repo

# Working directory
WORKDIR /kai-service

go get https://github.com/ZanLabs/kai.git

# Download And Install Kai
RUN git clone https://github.com/ZanLabs/kai.git KAI-REPO && \
    cd KAI-REPO && \
    go install && \
    mv kai .. && \
    mv cfg .. && \
    mv data .. && \ 
    mv config.yaml .. && \
    cd .. && \
    rm -rf KAI-REPO 

# Download weights
RUN curl -O http://pjreddie.com/media/files/yolo.weights >/dev/null 2>&1 

CMD [ "./kai" ]



