FROM yummybian/docker-darknet:latest


# Working directory
WORKDIR /kai-service

# Download And Install Kai
RUN go get -u github.com/ZanLabs/kai && \
    go install github.com/ZanLabs/kai && \
    cp /go/bin/kai ./kai

# Copy configurations
RUN cp -Ra /go/src/github.com/ZanLabs/kai/cfg ./ && \
    cp -Ra /go/src/github.com/ZanLabs/kai/data ./ && \
    cp /go/src/github.com/ZanLabs/kai/config.yaml ./config.yaml

# Download weights
RUN curl -O http://pjreddie.com/media/files/yolo.weights >/dev/null 2>&1 

EXPOSE 8000

ENTRYPOINT [ "./kai" ]
