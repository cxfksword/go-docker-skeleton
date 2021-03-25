# go-docker-skeleton [![license](https://img.shields.io/badge/license-Apache%202-blue?style=flat)](https://github.com/cxfksword/go-docker-skeleton/blob/master/LICENSE)[![version](https://img.shields.io/badge/version-0.1.0-blue.svg)](https://github.com/cxfksword/go-docker-skeleton/releases)
Skeleton for run go service in docker


# How to usage

1. copy all file to project directory
2. replace all `github.com/cxfksword/go-docker-skeleton` string to your repo
3. execute shell:
```
go mod vendor
go get -u https://github.com/cosmtrek/air
air -c .air.toml
```


# How to push DockerHub

1. register dockerhub and create a repo
2. goto `Account Settings -> Security` create aceess token
3. on Github, goto repo `Settings -> Secrets` add three github action variable
```
DOCKER_USERNAME
DOCKER_TOKEN
DOCKER_REPOSITORY
```