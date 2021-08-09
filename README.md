# IOT-PaaS message producer golang microservice

A simple golang microservice to publish data to a message queue. 


## Usage 

```bash
# cd to project directory and build executable
$ make build 

```

## Container build

```bash
# first build the base ubi build image (use golang version as tag)
$ podman build -t quay.io/<your-id>/ubi-base-builder:1.16.6 
$ podman push quay.io/<your-id>/ubu-base-builder:1.16.6

# version of golang as a tag
make container

```
## Container push

```bash
# version of golang as a tag
# use AUTH to set authfile if needed i.e (export AUTH=--authfile=/home/<user>/.docker/config.json)
make push

```

## Curl timing usage
```
curl -w "@curl-timing.txt" -o /dev/null -s "http://site-to-test

```

## Executing tests
```bash
make test 
make cover
# run sonarqube scanner (assuming sonarqube server is running)
# NB the SonarQube host and login will differ - please update it accordingly in the sonar-project.properties file
~/<path-to-sonarqube>/sonar-scanner-3.3.0.1492-linux/bin/sonar-scanner 

```
## Testing container 
```bash

# assume kafka brokers are running
# start the container 
podman run -it <registry>/iotpaas-message-producer
# curl the isalive endpoint
curl -k -H 'Token: xxxxx' -w '@curl-timing.txt'  http://127.0.0.1:9000/api/v1/sys/info/isalive

```

## Update for openshift
- removed all references to GOCD
