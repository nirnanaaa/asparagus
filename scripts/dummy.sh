#!/bin/bash

ETCD_VER=v3.2.9

# choose either URL
GOOGLE_URL=https://storage.googleapis.com/etcd
GITHUB_URL=https://github.com/coreos/etcd/releases/download
DOWNLOAD_URL=${GOOGLE_URL}

rm -f /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
rm -rf /tmp/etcd-download-test && mkdir -p /tmp/etcd-download-test

curl -L ${DOWNLOAD_URL}/${ETCD_VER}/etcd-${ETCD_VER}-linux-amd64.tar.gz -o /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz
tar xzvf /tmp/etcd-${ETCD_VER}-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1

/tmp/etcd-download-test/etcd --version
<<COMMENT
etcd Version: 3.2.9
Git SHA: f1d7dd8
Go Version: go1.8.4
Go OS/Arch: linux/amd64
COMMENT

ETCDCTL_API=3 /tmp/etcd-download-test/etcdctl version
<<COMMENT
etcdctl version: 3.2.9
API version: 3.2
COMMENT

/tmp/etcd-download-test/etcdctl set /cron/Config '{"SchedulerTickDuration":1000,"Enabled":true,"SecretKey":"","SignRequests":false}'
/tmp/etcd-download-test/etcdctl set /cron/Jobs/Dummy '{"Name":"Dummy","Expression":"* * * * *","URI":"https://httpbin.org/get","URIIsAbsolute":true}'
