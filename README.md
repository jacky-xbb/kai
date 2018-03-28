[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/ZanLabs/kai)
[![Go Report Card](https://goreportcard.com/badge/github.com/labstack/echo?style=flat-square)](https://goreportcard.com/report/github.com/ZanLabs/kai)
[![Build Status](http://img.shields.io/travis/labstack/echo.svg?style=flat-square)](https://travis-ci.org/ZanLabs/kai)
[![Codecov](https://img.shields.io/codecov/c/github/labstack/echo.svg?style=flat-square)](https://codecov.io/gh/ZanLabs/kai) 
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://raw.githubusercontent.com/ZanLabs/kai/master/README.md)

**Kai is a objection detection cloud services based on YOLO/Darknet.**

## Features
- Backend on the YOLO/Darknet of golang binding
- Convenient RESTful to use
- support S3 download and upload
- support Ftp download and upload

## Precondition
- nvidia-docker 2.0
- CUDA
- cudnn 

## Docker
```bash
git pull yummybian/kai
sudo docker run --runtime nvidia -it --rm -p 8000:8000 -v /path/to/config.yaml:/kai-service/config.yaml yummybian/kai bash
or
sudo docker run --runtime nvidia -d --name kai -p 8000:8000 -v /path/to/config.yaml:/kai-service/config.yaml yummybian/kai 
```

## Setting Up

First make sure you have [go-yolo](https://github.com/ZanLabs/go-yolo) installed on your machine. 


You can store jobs on memory or [MongoDB](https://www.mongodb.com/). On your `config.yaml` file:

- For MongoDB, set `dbdriver: "mongo"` and `mongohost: "your.mongo.host"`
- For memory, just set `dbdriver: "memory"` and you're good to go.

Please be aware that in case you use `memory`, Kai will persist the data only while the application is running.

**Finally** download [weight](http://pjreddie.com/media/files/yolo.weights) file from the darknet project into the root directory of kai.

Run!

```
$ make run
```

## Running tests

Make sure you have a local instance of [MongoDB](https://github.com/mongodb/mongo) running.

```
$ make test
```

## Using the API
### Creating a Job
In order to create a job you need to specify a HTTP or S3 address of a source and a S3 address for the destination. 

Given a JSON file called job.json:

```json
{
  "source": "http://www.example.com/example.jpg",
  "destination": "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/DIRECTORY",
  "cate": 0
}
```

or

```json
{
  "source": "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/DIRECTORY/example.jpg",
  "destination": "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/DIRECTORY",
  "cate": 0
}
```

**Note**: key cate 0 for image, 1 for video

Then, make a POST request to the API:

```Bash
$ curl -X POST -H "Content-Type: application/json" -d @job.json  http://api.host.com/jobs
< HTTP/1.1 200 OK
```

### Listing Jobs
```Bash
$ curl -X GET http://api.host.com/jobs
< HTTP/1.1 200 OK
```

```Bash
{
    "id": "7cJ5BtLwcFcQ8Vi8",
    "source": "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/DIRECTORY/example.jpg",
    "destination": "http://AWSKEY:AWSSECRET@BUCKET.s3.amazonaws.com/DIRECTORY",
    "mediaType": {
        "cate": 0,
        "name": "example",
        "container": "jpg"
    },
    "status": "created",
    "details": ""
}
```

### Getting Job Details
With the `Job ID`:

```Bash
$ curl -X GET http://api.host.com/jobs/7cJ5BtLwcFcQ8Vi8
```

### Starting the job
With the `Job ID`:

```Bash
$ curl -X GET http://api.host.com/jobs/7cJ5BtLwcFcQ8Vi8/start
```

Then you should request job details in order to follow the status of each step (downloading, processing, uploading).


## Contributing

1. Fork it
2. Create your feature branch: `git checkout -b my-awesome-new-feature`
3. Commit your changes: `git commit -m 'Add some awesome feature'`
4. Push to the branch: `git push origin my-awesome-new-feature`
5. Submit a pull request

## License

This code is under [Apache 2.0 License](https://github.com/ZanLabs/kai/blob/master/LICENSE).
