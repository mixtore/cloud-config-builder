# cloud-config-builder

config-builder is a simple template based centralized config builder
written in Go (Golang). it features a simple command-line to create
config files for cloud applications.

For now it supports the following config files:
* Kubernetes `ConfigMap`
* AppEngine `app.yaml`

## How to use

create a env file with some applications configurations eg:
```
$ cat test.env
>>> PORT=8080
```
```
go run main.go \
    -env-file=./test.env \
    -type=appengine \
    -output-file=./test.yaml \
    \
    -name=test \
    -runtime=ruby \
    -env=flex \
    -command='bundle exec rails server -p $PORT' \
    -scaling-min=1 \
    -scaling-max=4 \
    -scaling-cpu=0.45 \
    -resources-memory=3.7 \
    -resources-cpu-count=4 \
    -disable-healthcheck=true \

go run main.go \
    -env-file=./test.env \
    -type=kubernetes-configmap \
    -output-file=./test.yaml \
    \
    -namespace \
    -name=test \
```


## Licence
> This project is licensed under the terms of the MIT license.
